create table if not exists %s.raw_transaction_log
(
    log_hash                string not null,
    chain_id                numeric not null,
    tx_hash                 string,
    index                   numeric not null,
    contract_address        string not null,
    block_number            bignumeric not null,
    block_hash              string not null,
    block_time              numeric not null,
    tx_index                numeric not null,
    tx_contract_address     string,
    tx_to                   string,
    tx_value                bignumeric not null,
    data                    json,
    PRIMARY KEY (log_hash) NOT ENFORCED
) CLUSTER BY chain_id;