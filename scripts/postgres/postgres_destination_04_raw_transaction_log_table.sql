create table if not exists public.raw_transaction_log
(
    hash text not null constraint raw_transaction_log_pk primary key,
    index            numeric not null,
    chain_id numeric not null,
    block_number numeric not null,
    transaction_hash text not null,
    contract_address text not null,
    log jsonb not null
);