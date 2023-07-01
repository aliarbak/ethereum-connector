package datasource

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliarbak/ethereum-connector/configs"
	connectorerrors "github.com/aliarbak/ethereum-connector/errors"
	"github.com/aliarbak/ethereum-connector/model"
	"github.com/aliarbak/ethereum-connector/scripts/postgres"
	"github.com/jackc/pgx/v5"
	"math/big"
	"os"
	"time"
)

const (
	chainColumns                   = "id, block_number, paused, rpc_url, block_number_wait_period_in_ms, block_number_read_period_in_ms, concurrent_block_read_limit, updated_at"
	schemaColumns                  = "id, name, abi, events"
	schemaFilterColumns            = "id, chain_id, schema_id, contract_address, sync"
	rawTransactionLogSelectColumns = "hash, index, chain_id, block_number, block_hash, block_time, tx_index, tx_hash, tx_address, tx_to, tx_value, contract_address, log"
)

var (
	rawTransactionLogColumns = []string{"hash", "index", "chain_id", "block_number", "block_hash", "block_time", "tx_index", "tx_hash", "tx_address", "tx_to", "tx_value", "contract_address", "log"}
)

type postgresDataSourceFactory struct {
	config configs.PostgresDataSourceConfig
}

type postgresDataSource struct {
	conn                      *pgx.Conn
	persistRawTransactionLogs bool
}

func newPostgresFactory(config configs.PostgresDataSourceConfig) Factory {
	return &postgresDataSourceFactory{
		config: config,
	}
}

func (f postgresDataSourceFactory) CreateDataSource(ctx context.Context) (source Source, err error) {
	conn, err := pgx.Connect(ctx, f.config.ConnectionString)
	if err != nil {
		return
	}

	return &postgresDataSource{
		conn:                      conn,
		persistRawTransactionLogs: f.config.PersistRawTransactionLogs,
	}, err
}

func (r postgresDataSource) CreateChain(ctx context.Context, chain model.Chain) error {
	query := fmt.Sprintf("INSERT INTO chain(%s) VALUES($1, $2, $3, $4, $5, $6, $7, $8)", chainColumns)
	result, err := r.conn.Exec(ctx, query, chain.Id, chain.BlockNumber.String(), chain.Paused, chain.RpcUrl, chain.BlockNumberWaitPeriodInMS, chain.BlockNumberReadPeriodInMS, chain.ConcurrentBlockReadLimit, chain.UpdatedAt)
	if err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return connectorerrors.AlreadyExists("chain already exists with given id: %d", chain.Id)
		}
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("no rows affected when creating chain")
	}

	return nil
}

func (r *postgresDataSource) UpdateBlock(ctx context.Context, block model.Block) error {
	if !r.persistRawTransactionLogs {
		return r.updateBlockNumber(ctx, block)
	}

	return r.updateBlockNumberAndPersistRawLogs(ctx, block)
}

func (r postgresDataSource) UpdateChain(ctx context.Context, chain model.Chain) error {
	query := `UPDATE chain SET paused = $1, rpc_url = $2, block_number_wait_period_in_ms = $3, block_number_read_period_in_ms = $4, concurrent_block_read_limit = $5, updated_at = $6 WHERE id = $7`
	result, err := r.conn.Exec(ctx, query, chain.Paused, chain.RpcUrl, chain.BlockNumberWaitPeriodInMS, chain.BlockNumberReadPeriodInMS, chain.ConcurrentBlockReadLimit, chain.UpdatedAt, chain.Id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("no rows affected on current chain update")
	}

	return err
}

func (r *postgresDataSource) GetChain(ctx context.Context, chainId int64) (chain model.Chain, err error) {
	query := fmt.Sprintf("SELECT %s FROM chain where id = $1", chainColumns)
	row := r.conn.QueryRow(ctx, query, chainId)

	var blockNumber int64
	err = row.Scan(
		&chain.Id,
		&blockNumber,
		&chain.Paused,
		&chain.RpcUrl,
		&chain.BlockNumberWaitPeriodInMS,
		&chain.BlockNumberReadPeriodInMS,
		&chain.ConcurrentBlockReadLimit,
		&chain.UpdatedAt)

	if err == pgx.ErrNoRows {
		return chain, connectorerrors.NotFound("chain not found by id: %d", chainId)
	}

	chain.BlockNumber = big.NewInt(blockNumber)
	return
}

func (r *postgresDataSource) GetChains(ctx context.Context) (chains []model.Chain, err error) {
	query := fmt.Sprintf("SELECT %s FROM chain", chainColumns)
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return
	}
	for rows.Next() {
		var chain model.Chain
		var blockNumber int64
		err = rows.Scan(
			&chain.Id,
			&blockNumber,
			&chain.Paused,
			&chain.RpcUrl,
			&chain.BlockNumberWaitPeriodInMS,
			&chain.BlockNumberReadPeriodInMS,
			&chain.ConcurrentBlockReadLimit,
			&chain.UpdatedAt)

		if err != nil {
			return
		}
		chain.BlockNumber = big.NewInt(blockNumber)
		chains = append(chains, chain)
	}

	return
}

func (r postgresDataSource) CreateSchema(ctx context.Context, schema model.Schema) (err error) {
	query := fmt.Sprintf("INSERT INTO schema(%s) VALUES($1, $2, $3, $4)", schemaColumns)
	result, err := r.conn.Exec(ctx, query, schema.Id, schema.Name, schema.Abi, schema.Events)
	if err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return connectorerrors.AlreadyExists("schema already exists with given id: %d", schema.Id)
		}
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("no rows affected when creating schema")
	}

	return nil
}

func (r postgresDataSource) GetSchema(ctx context.Context, schemaId string) (schema model.Schema, err error) {
	query := fmt.Sprintf("SELECT %s FROM schema where id = $1", schemaColumns)
	row := r.conn.QueryRow(ctx, query, schemaId)
	err = row.Scan(
		&schema.Id,
		&schema.Name,
		&schema.Abi,
		&schema.Events)

	if err == pgx.ErrNoRows {
		return schema, connectorerrors.NotFound("schema not found by id: %d", schemaId)
	}

	return
}

func (r postgresDataSource) GetSchemas(ctx context.Context) ([]model.Schema, error) {
	query := fmt.Sprintf(`SELECT %s FROM schema`, schemaColumns)
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []model.Schema
	for rows.Next() {
		var schema model.Schema
		err := rows.Scan(&schema.Id, &schema.Name, &schema.Abi, &schema.Events)
		if err != nil {
			return nil, err
		}

		results = append(results, schema)
	}

	return results, nil
}

func (r postgresDataSource) GetSchemasWithFilters(ctx context.Context) (schemas []model.Schema, err error) {
	schemas, err = r.GetSchemas(ctx)
	if err != nil {
		return
	}

	filters, err := r.getAllSchemaFilters(ctx)
	if err != nil {
		return
	}

	model.AddFiltersToSchemas(schemas, filters)
	return
}

func (r postgresDataSource) CreateSchemaFilter(ctx context.Context, filter model.SchemaFilter) (err error) {
	query := fmt.Sprintf("INSERT INTO schema_contract_filter(%s) VALUES($1, $2, $3, $4, $5)", schemaFilterColumns)
	result, err := r.conn.Exec(ctx, query, filter.Id, filter.ChainId, filter.SchemaId, filter.ContractAddress, filter.Sync)
	if err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return connectorerrors.AlreadyExists("schema filter already exists with given contract address and schema")
		}
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("no rows affected when creating schema filter")
	}

	return nil
}

func (r postgresDataSource) GetSchemaFilter(ctx context.Context, schemaFilterId string) (filter model.SchemaFilter, err error) {
	query := fmt.Sprintf("SELECT %s FROM schema_contract_filter where id = $1", schemaFilterColumns)
	row := r.conn.QueryRow(ctx, query, schemaFilterId)
	err = row.Scan(
		&filter.Id,
		&filter.ChainId,
		&filter.SchemaId,
		&filter.ContractAddress,
		&filter.Sync)

	if err == pgx.ErrNoRows {
		return filter, connectorerrors.NotFound("schema filter not found by id: %d", schemaFilterId)
	}

	return
}

func (r postgresDataSource) GetSchemaFiltersBySchemaId(ctx context.Context, schemaId string) ([]model.SchemaFilter, error) {
	query := fmt.Sprintf(`SELECT %s FROM schema_contract_filter WHERE schema_id = $1`, schemaFilterColumns)
	rows, err := r.conn.Query(ctx, query, schemaId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []model.SchemaFilter
	for rows.Next() {
		var filter model.SchemaFilter
		err := rows.Scan(&filter.Id,
			&filter.ChainId,
			&filter.SchemaId,
			&filter.ContractAddress,
			&filter.Sync)

		if err != nil {
			return nil, err
		}

		results = append(results, filter)
	}

	return results, nil
}

func (r postgresDataSource) GetOutOfSyncSchemaFilters(ctx context.Context) ([]model.SchemaFilter, error) {
	query := fmt.Sprintf(`SELECT %s FROM schema_contract_filter WHERE (sync -> 'completed')::bool = false`, schemaFilterColumns)
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []model.SchemaFilter
	for rows.Next() {
		var filter model.SchemaFilter
		err := rows.Scan(&filter.Id,
			&filter.ChainId,
			&filter.SchemaId,
			&filter.ContractAddress,
			&filter.Sync)

		if err != nil {
			return nil, err
		}

		results = append(results, filter)
	}

	return results, nil
}

func (r postgresDataSource) UpdateSchemaFilterSync(ctx context.Context, filterId string, sync model.SchemaFilterSync) error {
	query := `UPDATE schema_contract_filter SET sync = $1 WHERE id = $2`
	result, err := r.conn.Exec(ctx, query, sync, filterId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("no rows affected on schema filter sync update")
	}

	return err
}

func (r postgresDataSource) GetRawTransactionBlockNumbersForContract(ctx context.Context, chainId int64, contractAddress string) (blockNumbers []*big.Int, err error) {
	query := `SELECT DISTINCT(block_number) FROM raw_transaction_log WHERE contract_address = $1 and chain_id = $2 ORDER BY block_number ASC`
	rows, err := r.conn.Query(ctx, query, contractAddress, chainId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []*big.Int
	for rows.Next() {
		var blockNumber int64
		err := rows.Scan(&blockNumber)
		if err != nil {
			return nil, err
		}

		results = append(results, big.NewInt(blockNumber))
	}

	return results, nil
}

func (r postgresDataSource) GetRawTransactionsForContract(ctx context.Context, chainId int64, blockNumber *big.Int, contractAddress string) ([]model.RawTransactionLog, error) {
	query := fmt.Sprintf("SELECT %s FROM raw_transaction_log WHERE chain_id = $1 AND contract_address = $2 AND block_number = $3 ORDER BY index ASC", rawTransactionLogSelectColumns)
	rows, err := r.conn.Query(ctx, query, chainId, contractAddress, blockNumber.String())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []model.RawTransactionLog
	for rows.Next() {
		var log model.RawTransactionLog
		var txValue string
		var bNumber int64
		err := rows.Scan(&log.Hash,
			&log.Index,
			&log.ChainId,
			&bNumber,
			&log.BlockHash,
			&log.BlockTime,
			&log.TxIndex,
			&log.TxHash,
			&log.TxAddress,
			&log.TxTo,
			&txValue,
			&log.ContractAddress,
			&log.RawData)

		log.TxValue, _ = new(big.Int).SetString(txValue, 10)
		log.BlockNumber = blockNumber

		if err != nil {
			return nil, err
		}

		results = append(results, log)
	}

	return results, nil
}

func (r postgresDataSource) RawTransactionLogsPersisted() bool {
	return r.persistRawTransactionLogs
}

func (r *postgresDataSource) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func (r postgresDataSource) updateBlockNumber(ctx context.Context, block model.Block) error {
	query := `UPDATE chain SET block_number = $1, updated_at = $2 WHERE id = $3`
	result, err := r.conn.Exec(ctx, query, block.Number.Int64(), time.Now().UTC().UnixMilli(), block.ChainId)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("no rows affected on current blockNumber update")
	}

	return err
}

func (r postgresDataSource) getAllSchemaFilters(ctx context.Context) ([]model.SchemaFilter, error) {
	query := fmt.Sprintf(`SELECT %s FROM schema_contract_filter`, schemaFilterColumns)
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []model.SchemaFilter
	for rows.Next() {
		var filter model.SchemaFilter
		err := rows.Scan(&filter.Id,
			&filter.ChainId,
			&filter.SchemaId,
			&filter.ContractAddress,
			&filter.Sync)

		if err != nil {
			return nil, err
		}

		results = append(results, filter)
	}

	return results, nil
}

func (r postgresDataSource) updateBlockNumberAndPersistRawLogs(ctx context.Context, block model.Block) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	res, err := tx.Exec(ctx, "UPDATE chain SET block_number = $1, updated_at = $2 WHERE id = $3", block.Number.Int64(), time.Now().UTC().UnixMilli(), block.ChainId)
	if err != nil {
		return err
	}

	count := res.RowsAffected()
	if count == 0 {
		return errors.New("block update failed")
	}

	var rawTransactionLogRows [][]interface{}
	for _, transaction := range block.Transactions {
		for _, transactionLog := range transaction.Logs {
			transactionLogRow := make([]interface{}, 13)
			transactionLogRow[0] = transactionLog[model.LogHashLogField]
			transactionLogRow[1] = transactionLog[model.LogIndexLogField]
			transactionLogRow[2] = block.ChainId
			transactionLogRow[3] = block.Number.String()
			transactionLogRow[4] = block.Hash
			transactionLogRow[5] = block.Time
			transactionLogRow[6] = transaction.Index
			transactionLogRow[7] = transaction.Hash
			transactionLogRow[8] = transaction.ContractAddress
			transactionLogRow[9] = transaction.To
			transactionLogRow[10] = transaction.Value.String()
			transactionLogRow[11] = transactionLog[model.ContractAddressLogField]
			transactionLogRow[12] = transactionLog[model.RawDataLogField]
			rawTransactionLogRows = append(rawTransactionLogRows, transactionLogRow)
		}
	}

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

	if err = tx.Commit(ctx); err != nil {
		if connectorerrors.IsPostgresDuplicateKeyError(err) {
			return nil
		}

		return err
	}

	return err
}

func (r postgresDataSource) init(ctx context.Context) (err error) {
	if err = r.execSqlFromFile(ctx, postgres.InitDataSourceChainTableScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDataSourceChainTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, postgres.InitDataSourceSchemaTableScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDataSourceSchemaTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, postgres.InitDataSourceRawTransactionLogTableScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDataSourceRawTransactionLogTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, postgres.InitDataSourceRawTransactionLogIndexScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDataSourceRawTransactionLogIndexScriptFile)
	}

	if err = r.execSqlFromFile(ctx, postgres.InitDataSourceSchemaContractFilterTableScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDataSourceSchemaContractFilterTableScriptFile)
	}

	if err = r.execSqlFromFile(ctx, postgres.InitDataSourceSchemaContractFilterIndexScriptFile); err != nil {
		return connectorerrors.From(err, "script failed (file: %s)", postgres.InitDataSourceSchemaContractFilterIndexScriptFile)
	}

	return
}

func (r postgresDataSource) execSqlFromFile(ctx context.Context, file string, params ...any) error {
	sqlFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	sql := string(sqlFile)
	_, err = r.conn.Exec(ctx, sql, params...)
	return err
}
