create table if not exists public.transaction
(
    hash             text    not null constraint transaction_pk primary key,
    index            numeric not null,
    block_hash text constraint transaction_block_hash_fk references block,
    chain_id         numeric not null,
    contract_address text,
    "to"             text,
    value            numeric
);
