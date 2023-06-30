create table if not exists public.chain
(
    id bigint not null constraint chain_pk primary key,
    block_number numeric,
    paused       boolean default true not null,
    rpc_url      text not null,
    block_number_wait_period_in_ms numeric default 10000 not null,
    block_number_read_period_in_ms numeric default 100 not null,
    concurrent_block_read_limit numeric default 1 not null,
    updated_at numeric not null
);