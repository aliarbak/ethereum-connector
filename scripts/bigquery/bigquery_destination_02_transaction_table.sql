create table if not exists %s.transaction
(
    chain_id         numeric not null,
    block_number     bignumeric not null,
    block_hash       string not null,
    block_time       numeric not null,
    contract_address string,
    tx_to            string,
    tx_hash          string  not null,
    index            numeric not null,
    value            bignumeric,
    PRIMARY KEY (tx_hash) NOT ENFORCED
) CLUSTER BY chain_id;