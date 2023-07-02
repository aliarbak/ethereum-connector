package destination

import (
	"cloud.google.com/go/bigquery"
	storage "cloud.google.com/go/bigquery/storage/apiv1beta2"
	"cloud.google.com/go/bigquery/storage/apiv1beta2/storagepb"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"context"
	"fmt"
	"github.com/aliarbak/ethereum-connector/configs"
	"github.com/aliarbak/ethereum-connector/destination/message"
	connectorerrors "github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	bigqueryscript "github.com/aliarbak/ethereum-connector/scripts/bigquery"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"os"
)

type bigqueryDestinationFactory struct {
	config configs.BigQueryDestinationConfig
}

type bigqueryDestination struct {
	projectId                 string
	datasetName               string
	persistRawTransactionLogs bool
	deliveryGuarantee         DeliveryGuarantee
	client                    *storage.BigQueryWriteClient
}

func newBigQueryFactory(config configs.BigQueryDestinationConfig) Factory {
	return &bigqueryDestinationFactory{
		config: config,
	}
}

func (f bigqueryDestinationFactory) CreateDestination(ctx context.Context) (dest Destination, err error) {
	if f.config.DeliveryGuarantee != string(AtLeastOnceDeliveryGuarantee) && f.config.DeliveryGuarantee != string(AtMostOnceDeliveryGuarantee) {
		return nil, connectorerrors.InvalidInput("invalid bigquery destination delivery guarantee: %s", f.config.DeliveryGuarantee)
	}

	client, err := storage.NewBigQueryWriteClient(ctx)
	if err != nil {
		return dest, err
	}

	return &bigqueryDestination{
		projectId:                 f.config.ProjectId,
		datasetName:               f.config.Dataset,
		persistRawTransactionLogs: f.config.PersistRawTransactionLogs,
		deliveryGuarantee:         DeliveryGuarantee(f.config.DeliveryGuarantee),
		client:                    client,
	}, err
}

func (r bigqueryDestination) DeliveryGuarantee() DeliveryGuarantee {
	return ExactlyOnceDeliveryGuarantee
}

func (r bigqueryDestination) SendBlock(ctx context.Context, block model.Block) (err error) {
	err = r.sendBlockMessage(ctx, block)
	if err != nil {
		return err
	}

	return r.sendTransactions(ctx, block)
}

func (r bigqueryDestination) SendSyncLogs(ctx context.Context, block model.Block) (err error) {
	_, transactionLogData, _, err := r.prepareTransactionMessages(block)
	if err != nil {
		return err
	}

	if len(transactionLogData) > 0 {
		if err = r.sendMessages(ctx, "transaction_log", (&destmessage.TransactionLog{}).ProtoReflect(), transactionLogData); err != nil {
			return err
		}
	}

	return
}

func (r bigqueryDestination) sendBlockMessage(ctx context.Context, block model.Block) (err error) {
	var opts proto.MarshalOptions
	var data [][]byte

	blockMessage := destmessage.NewProtoBlockMessage(block)

	buf, err := opts.Marshal(&blockMessage)
	if err != nil {
		return connectorerrors.From(err, "block protobuf marshall failed")
	}

	data = append(data, buf)

	return r.sendMessages(ctx, "block", blockMessage.ProtoReflect(), data)
}

func (r bigqueryDestination) sendTransactions(ctx context.Context, block model.Block) (err error) {
	transactionData, transactionLogData, rawTransactionLogData, err := r.prepareTransactionMessages(block)
	if err != nil {
		return err
	}

	if err = r.sendMessages(ctx, "transaction", (&destmessage.Transaction{}).ProtoReflect(), transactionData); err != nil {
		return err
	}

	if len(transactionLogData) > 0 {
		if err = r.sendMessages(ctx, "transaction_log", (&destmessage.TransactionLog{}).ProtoReflect(), transactionLogData); err != nil {
			return err
		}
	}

	if r.persistRawTransactionLogs && len(rawTransactionLogData) > 0 {
		if err = r.sendMessages(ctx, "raw_transaction_log", (&destmessage.RawTransactionLog{}).ProtoReflect(), rawTransactionLogData); err != nil {
			return err
		}
	}

	return
}

func (r bigqueryDestination) sendMessages(ctx context.Context, table string, message protoreflect.Message, data [][]byte) (err error) {
	descriptor, err := adapt.NormalizeDescriptor(message.Descriptor())
	if err != nil {
		return connectorerrors.From(err, "%s descriptor failed", table)
	}

	resp, err := r.client.CreateWriteStream(ctx, &storagepb.CreateWriteStreamRequest{
		Parent: fmt.Sprintf("projects/%s/datasets/%s/tables/%s", r.projectId, r.datasetName, table),
		WriteStream: &storagepb.WriteStream{
			Type: storagepb.WriteStream_COMMITTED,
		},
	})

	if err != nil {
		return connectorerrors.From(err, "%s writeStream failed", table)
	}

	stream, err := r.client.AppendRows(ctx)
	if err != nil {
		return connectorerrors.From(err, "%s appendRows failed", table)
	}

	err = stream.Send(&storagepb.AppendRowsRequest{
		WriteStream: resp.Name,
		TraceId:     "ethereum-connector",
		Rows: &storagepb.AppendRowsRequest_ProtoRows{
			ProtoRows: &storagepb.AppendRowsRequest_ProtoData{
				WriterSchema: &storagepb.ProtoSchema{
					ProtoDescriptor: descriptor,
				},
				Rows: &storagepb.ProtoRows{
					SerializedRows: data,
				},
			},
		},
	})

	if err != nil {
		return connectorerrors.From(err, "%s stream send failed", table)
	}

	rc, err := stream.Recv()
	if err != nil {
		return connectorerrors.From(err, "%s stream recv failed", table)
	}

	if rErr := rc.GetError(); rErr != nil {
		return connectorerrors.New("%s stream send result error: %v", table, rErr)
	} else if rResult := rc.GetAppendResult(); rResult != nil {
		log.Trace().Msgf("%s stream offset is %d", table, rResult.Offset.Value)
	}

	return
}

func (r bigqueryDestination) prepareTransactionMessages(block model.Block) (transactionData [][]byte, transactionLogData [][]byte, rawTransactionLogData [][]byte, err error) {
	var opts proto.MarshalOptions
	for _, transaction := range block.Transactions {
		transactionMessage := destmessage.NewProtoTransactionMessage(block, transaction)
		buf, err := opts.Marshal(&transactionMessage)
		if err != nil {
			return transactionData, transactionLogData, rawTransactionLogData, connectorerrors.From(err, "transaction protobuf marshall failed for txHash: %s", transaction.Hash)
		}
		transactionData = append(transactionData, buf)

		for _, transactionLog := range transaction.Logs {
			rawTransactionLogMessage := destmessage.NewProtoRawTransactionLogMessage(block, transaction, transactionLog)
			buf, err = opts.Marshal(&rawTransactionLogMessage)
			if err != nil {
				return transactionData, transactionLogData, rawTransactionLogData, connectorerrors.From(err, "raw transaction log protobuf marshall failed for txHash: %s, logIndex: %d", transaction.Hash, transactionLog[model.LogIndexLogField])
			}

			rawTransactionLogData = append(rawTransactionLogData, buf)

			if transactionLog[model.EventNameLogField] == nil || len(transactionLog[model.EventNameLogField].(string)) == 0 {
				continue
			}

			transactionLogMessage, err := destmessage.NewProtoTransactionLogMessage(block, transaction, transactionLog)
			if err != nil {
				return transactionData, transactionLogData, rawTransactionLogData, err
			}

			buf, err = opts.Marshal(&transactionLogMessage)
			if err != nil {
				return transactionData, transactionLogData, rawTransactionLogData, connectorerrors.From(err, "transaction log protobuf marshall failed for txHash: %s, logIndex: %d", transaction.Hash, transactionLog[model.LogIndexLogField])
			}

			transactionLogData = append(transactionLogData, buf)
		}
	}

	return
}

func (r *bigqueryDestination) Close(_ context.Context) error {
	return r.client.Close()
}

func (r bigqueryDestination) init(ctx context.Context) (err error) {
	client, err := bigquery.NewClient(ctx, r.projectId)
	if err != nil {
		return
	}

	defer client.Close()

	if err = r.execSqlFromFile(ctx, client, bigqueryscript.InitDestinationBlockTableScriptFile, r.datasetName); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", bigqueryscript.InitDestinationBlockTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, client, bigqueryscript.InitDestinationTransactionTableScriptFile, r.datasetName); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", bigqueryscript.InitDestinationTransactionTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, client, bigqueryscript.InitDestinationTransactionLogTableScriptFile, r.datasetName); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", bigqueryscript.InitDestinationTransactionTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, client, bigqueryscript.InitDestinationRawTransactionLogTableScriptFile, r.datasetName); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", bigqueryscript.InitDestinationRawTransactionLogTableScriptFile)
	}

	return
}

func (r bigqueryDestination) execSqlFromFile(ctx context.Context, client *bigquery.Client, file string, params ...any) error {
	sqlFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	sql := string(sqlFile)
	_, err = client.Query(fmt.Sprintf(sql, params...)).Run(ctx)
	return err
}
