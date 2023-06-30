package evm

import (
	"context"
	"github.com/aliarbak/ethereum-connector/datasource"
	"github.com/aliarbak/ethereum-connector/destination"
	"github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/rs/zerolog/log"
	"math/big"
)

func RunContractSyncWorkers(ctx context.Context, dataSourceFactory datasource.Factory, destinationFactory destination.Factory) error {
	dataSource, err := dataSourceFactory.CreateDataSource(ctx)
	if err != nil {
		return errors.From(err, "datasource creation failed for run listeners")
	}

	defer dataSource.Close(ctx)

	filters, err := dataSource.GetOutOfSyncSchemaFilters(ctx)
	if err != nil {
		return err
	}

	for _, filter := range filters {
		err = RunContractSyncWorker(ctx, filter.Id, dataSourceFactory, destinationFactory)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunContractSyncWorker(ctx context.Context, schemaFilterId string, dataSourceFactory datasource.Factory, destinationFactory destination.Factory) error {
	dataSource, err := dataSourceFactory.CreateDataSource(ctx)
	if err != nil {
		return errors.From(err, "datasource creation failed for sync listener for schemaFilterId: %s", schemaFilterId)
	}

	dest, err := destinationFactory.CreateDestination(ctx)
	if err != nil {
		return errors.From(err, "destination creation failed for sync listener for schemaFilterId: %s", schemaFilterId)
	}

	go func() {
		err := runContractSyncWorker(ctx, schemaFilterId, dataSource, dest)
		defer dataSource.Close(ctx)
		defer dest.Close(ctx)
		if err != nil {
			log.Fatal().Err(err).Msgf("sync listener crashed for schemaFilterId: %s", schemaFilterId)
		}

		log.Info().Msgf("sync listener completed for schemaFilterId: %s", schemaFilterId)
	}()

	log.Info().Msgf("sync listener is running for schemaFilterId: %s", schemaFilterId)
	return nil
}

func runContractSyncWorker(ctx context.Context, schemaFilterId string, dataSource datasource.Source, dest destination.Destination) error {
	filter, err := dataSource.GetSchemaFilter(ctx, schemaFilterId)
	if err != nil {
		return errors.From(err, "schema filter read from datasource failed for schemaFilterId: %s", schemaFilterId)
	}

	if filter.Sync == nil {
		return nil
	}

	schema, err := dataSource.GetSchema(ctx, filter.SchemaId)
	if err != nil {
		return errors.From(err, "schema read from datasource failed for schemaId: %s, schemaFilterId: %s", filter.SchemaId, schemaFilterId)
	}

	blockNumbers, err := dataSource.GetRawTransactionBlockNumbersForContract(ctx, filter.ChainId, filter.ContractAddress)
	if err != nil {
		return errors.From(err, "blockNumbers read from datasource failed for schemaFilterId: %s", schemaFilterId)
	}

	for _, blockNumber := range blockNumbers {
		filter, err = dataSource.GetSchemaFilter(ctx, schemaFilterId)
		if err != nil {
			return errors.From(err, "schema filter read from datasource failed for schemaFilterId: %s", schemaFilterId)
		}

		if blockNumber.Int64() < filter.Sync.CurrentBlockNumber {
			continue
		}

		if blockNumber.Int64() > filter.Sync.ToBlockNumber {
			break
		}

		rawTransactionLogs, err := dataSource.GetRawTransactionsForContract(ctx, filter.ChainId, blockNumber, filter.ContractAddress)
		if err != nil {
			return errors.From(err, "raw transaction logs read from datasource failed for schemaFilterId: %s, blockNumber: %s", schemaFilterId, blockNumber.String())
		}

		block, err := mapBlock(filter, schema, blockNumber, rawTransactionLogs)
		if err != nil {
			return errors.From(err, "raw transaction logs mapping block failed for schemaFilterId: %s, blockNumber: %s", schemaFilterId, blockNumber.String())
		}

		err = dest.SendSyncLogs(ctx, block)
		if err != nil {
			return errors.From(err, "send sync logs failed for schemaFilterId: %s, blockNumber: %s", schemaFilterId, blockNumber.String())
		}

		err = dataSource.UpdateSchemaFilterSync(ctx, schemaFilterId, model.SchemaFilterSync{
			FromBlockNumber:    filter.Sync.FromBlockNumber,
			ToBlockNumber:      filter.Sync.ToBlockNumber,
			CurrentBlockNumber: blockNumber.Int64(),
			Completed:          false,
		})

		if err != nil {
			return errors.From(err, "sync update failed for schemaFilterId: %s, blockNumber: %s", schemaFilterId, blockNumber.String())
		}
	}

	filter, err = dataSource.GetSchemaFilter(ctx, schemaFilterId)
	if err != nil {
		return errors.From(err, "schema filter read from datasource failed for schemaFilterId: %s", schemaFilterId)
	}

	err = dataSource.UpdateSchemaFilterSync(ctx, schemaFilterId, model.SchemaFilterSync{
		FromBlockNumber:    filter.Sync.FromBlockNumber,
		ToBlockNumber:      filter.Sync.ToBlockNumber,
		CurrentBlockNumber: filter.Sync.CurrentBlockNumber,
		Completed:          true,
	})

	if err != nil {
		return errors.From(err, "sync update as completed failed for schemaFilterId: %s", schemaFilterId)
	}

	return nil
}

func mapBlock(filter model.SchemaFilter, schema model.Schema, blockNumber *big.Int, rawTransactionLogs []model.RawTransactionLog) (block model.Block, err error) {
	var transactions []model.Transaction
	for _, rawTransactionLog := range rawTransactionLogs {
		log, err := mapRawLog(rawTransactionLog.RawData, schema)
		if err != nil {
			return block, errors.From(err, "raw transaction log mapping failed for schemaFilterId: %s, blockNumber: %s, txHash: %s, logIndex: %d", filter.Id, blockNumber.String(), rawTransactionLog.TxHash, rawTransactionLog.Index)
		}

		if log == nil {
			continue
		}

		transactionFound := false
		for i, transaction := range transactions {
			if transaction.Hash == rawTransactionLog.TxHash {
				transactions[i].Logs = append(transactions[i].Logs, log)
				transactionFound = true
				break
			}
		}

		if !transactionFound {
			to := ""
			if rawTransactionLog.TxTo != nil {
				to = *rawTransactionLog.TxTo
			}

			transactions = append(transactions, model.Transaction{
				ContractAddress: rawTransactionLog.TxAddress,
				To:              to,
				Index:           rawTransactionLog.TxIndex,
				Hash:            rawTransactionLog.TxHash,
				Value:           rawTransactionLog.TxValue,
				Data:            nil,
				Logs:            []map[string]interface{}{log},
			})
		}
	}

	return model.Block{
		ChainId:      filter.ChainId,
		Number:       blockNumber,
		Time:         rawTransactionLogs[0].BlockTime,
		Hash:         rawTransactionLogs[0].BlockHash,
		Transactions: transactions,
	}, nil
}
