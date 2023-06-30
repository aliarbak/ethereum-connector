package postgres

const (
	InitDataSourceChainTableScriptFile                = "scripts/postgres/postgres_datasource_01_chain_table.sql"
	InitDataSourceSchemaTableScriptFile               = "scripts/postgres/postgres_datasource_02_schema_table.sql"
	InitDataSourceRawTransactionLogTableScriptFile    = "scripts/postgres/postgres_datasource_03_raw_transaction_log_table.sql"
	InitDataSourceRawTransactionLogIndexScriptFile    = "scripts/postgres/postgres_datasource_04_raw_transaction_log_index.sql"
	InitDataSourceSchemaContractFilterTableScriptFile = "scripts/postgres/postgres_datasource_05_schema_contract_filter_table.sql"
	InitDataSourceSchemaContractFilterIndexScriptFile = "scripts/postgres/postgres_datasource_06_schema_contract_filter_index.sql"
	InitDestinationBlockTableScriptFile               = "scripts/postgres/postgres_destination_01_block_table.sql"
	InitDestinationTransactionTableScriptFile         = "scripts/postgres/postgres_destination_02_transaction_table.sql"
	InitDestinationTransactionLogTableScriptFile      = "scripts/postgres/postgres_destination_03_transaction_log_table.sql"
	InitDestinationRawTransactionLogTableScriptFile   = "scripts/postgres/postgres_destination_04_raw_transaction_log_table.sql"
	InitDestinationRawTransactionLogIndexScriptFile   = "scripts/postgres/postgres_destination_05_raw_transaction_log_index.sql"
)
