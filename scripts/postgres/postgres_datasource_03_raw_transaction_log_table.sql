create table if not exists public.raw_transaction_log
(
    hash text not null constraint raw_transaction_log_pk primary key,
    index            numeric not null,
    chain_id numeric not null,
    block_number numeric not null,
    block_hash text not null,
    block_time numeric not null,
    tx_index numeric not null,
    tx_hash text not null,
    tx_address text,
    tx_to text,
    tx_value numeric not null,
    contract_address text not null,
    log jsonb not null
);