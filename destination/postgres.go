package destination

import (
	"context"
	"errors"
	"github.com/aliarbak/ethereum-connector/configs"
	connectorerrors "github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/aliarbak/ethereum-connector/scripts/postgres"
	"github.com/aliarbak/ethereum-connector/utils"
	"github.com/jackc/pgx/v5"
	"os"
	"time"
)

var (
	transactionColumns       = []string{"hash", "index", "block_hash", "chain_id", "contract_address", "to", "value"}
	transactionLogColumns    = []string{"hash", "index", "transaction_hash", "contract_address", "event_name", "event_alias", "data"}
	rawTransactionLogColumns = []string{"hash", "index", "chain_id", "block_number", "transaction_hash", "contract_address", "log"}
	extractedLogFields       = []string{model.EventNameLogField, model.EventAliasLogField, model.ContractAddressLogField, model.LogHashLogField, model.LogIndexLogField, model.RawDataLogField}
)

type postgresDestinationFactory struct {
	config configs.PostgresDestinationConfig
}

type postgresDestination struct {
	conn                      *pgx.Conn
	persistRawTransactionLogs bool
}

func newPostgresFactory(config configs.PostgresDestinationConfig) Factory {
	return &postgresDestinationFactory{
		config: config,
	}
}

func (f postgresDestinationFactory) CreateDestination(ctx context.Context) (dest Destination, err error) {
	conn, err := pgx.Connect(ctx, f.config.ConnectionString)
	if err != nil {
		return
	}

	return &postgresDestination{
		conn:                      conn,
		persistRawTransactionLogs: f.config.PersistRawTransactionLogs,
	}, err
}

func (r postgresDestination) DeliveryGuarantee() DeliveryGuarantee {
	return ExactlyOnceDeliveryGuarantee
}

func (r postgresDestination) SendBlock(ctx context.Context, block model.Block) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	res, err := tx.Exec(ctx, "insert into block(hash, chain_id, number, time, created_at) VALUES($1, $2, $3, $4, $5)", block.Hash, block.ChainId, block.Number, block.Time, time.Now().UTC().UnixMilli())
	if err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return nil
		}
		return err
	}

	count := res.RowsAffected()
	if count == 0 {
		return errors.New("block insert failed")
	}

	var transactionRows [][]interface{}
	var transactionLogRows [][]interface{}
	var rawTransactionLogRows [][]interface{}
	for _, transaction := range block.Transactions {
		transactionRow := make([]interface{}, 7)
		transactionRow[0] = transaction.Hash
		transactionRow[1] = transaction.Index
		transactionRow[2] = block.Hash
		transactionRow[3] = block.ChainId
		transactionRow[4] = transaction.ContractAddress
		transactionRow[5] = transaction.To
		transactionRow[6] = transaction.Value.String()
		transactionRows = append(transactionRows, transactionRow)

		for _, transactionLog := range transaction.Logs {
			rawTransactionLogRow := make([]interface{}, 7)
			rawTransactionLogRow[0] = transactionLog[model.LogHashLogField]
			rawTransactionLogRow[1] = transactionLog[model.LogIndexLogField]
			rawTransactionLogRow[2] = block.ChainId
			rawTransactionLogRow[3] = block.Number.String()
			rawTransactionLogRow[4] = transaction.Hash
			rawTransactionLogRow[5] = transactionLog[model.ContractAddressLogField]
			rawTransactionLogRow[6] = transactionLog[model.RawDataLogField]
			rawTransactionLogRows = append(rawTransactionLogRows, rawTransactionLogRow)

			if transactionLog[model.EventNameLogField] == nil || len(transactionLog[model.EventNameLogField].(string)) == 0 {
				continue
			}

			transactionLogRow := make([]interface{}, 7)
			transactionLogRow[0] = transactionLog[model.LogHashLogField]
			transactionLogRow[1] = transactionLog[model.LogIndexLogField]
			transactionLogRow[2] = transaction.Hash
			transactionLogRow[3] = transactionLog[model.ContractAddressLogField]
			transactionLogRow[4] = transactionLog[model.EventNameLogField]
			transactionLogRow[5] = transactionLog[model.EventAliasLogField]
			transactionLogRow[6] = getLogDataMap(transactionLog)
			transactionLogRows = append(transactionLogRows, transactionLogRow)
		}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"transaction"},
		transactionColumns,
		pgx.CopyFromRows(transactionRows),
	)

	if err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return nil
		}

		return err
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"transaction_log"},
		transactionLogColumns,
		pgx.CopyFromRows(transactionLogRows),
	)

	if err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return nil
		}

		return err
	}

	if r.persistRawTransactionLogs {
		_, err = tx.CopyFrom(
			ctx,
			pgx.Identifier{"raw_transaction_log"},
			rawTransactionLogColumns,
			pgx.CopyFromRows(rawTransactionLogRows),
		)

		if err != nil {
			if connectorerrors.IsPostgresDuplicateKeyError(err) {
				return nil
			}
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return nil
		}
		return err
	}

	return nil
}

func (r postgresDestination) SendSyncLogs(ctx context.Context, block model.Block) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	var transactionLogRows [][]interface{}
	for _, transaction := range block.Transactions {
		for _, transactionLog := range transaction.Logs {
			if transactionLog[model.EventNameLogField] == nil || len(transactionLog[model.EventNameLogField].(string)) == 0 {
				continue
			}

			transactionLogRow := make([]interface{}, 7)
			transactionLogRow[0] = transactionLog[model.LogHashLogField]
			transactionLogRow[1] = transactionLog[model.LogIndexLogField]
			transactionLogRow[2] = transaction.Hash
			transactionLogRow[3] = transactionLog[model.ContractAddressLogField]
			transactionLogRow[4] = transactionLog[model.EventNameLogField]
			transactionLogRow[5] = transactionLog[model.EventAliasLogField]
			transactionLogRow[6] = getLogDataMap(transactionLog)
			transactionLogRows = append(transactionLogRows, transactionLogRow)
		}
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"transaction_log"},
		transactionLogColumns,
		pgx.CopyFromRows(transactionLogRows),
	)

	if err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return nil
		}

		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return nil
		}
		return err
	}

	return nil
}

func (r *postgresDestination) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func (r postgresDestination) init(ctx context.Context) (err error) {
	if err = r.execSqlFromFile(ctx, postgres.InitDestinationBlockTableScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDestinationBlockTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, postgres.InitDestinationTransactionTableScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDestinationTransactionTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, postgres.InitDestinationTransactionLogTableScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDestinationTransactionTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, postgres.InitDestinationRawTransactionLogTableScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDestinationRawTransactionLogTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, postgres.InitDestinationRawTransactionLogIndexScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDestinationRawTransactionLogIndexScriptFile)
	}

	return
}

func (r postgresDestination) execSqlFromFile(ctx context.Context, file string, params ...any) error {
	sqlFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	sql := string(sqlFile)
	_, err = r.conn.Exec(ctx, sql, params...)
	return err
}

func getLogDataMap(originalMap map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	for k, v := range originalMap {
		if utils.StringArrayContains(extractedLogFields, k) {
			continue
		}

		data[k] = v
	}

	return data
}
