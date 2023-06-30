create table if not exists public.schema_contract_filter
(
    id               text not null constraint schema_contract_filter_pk primary key,
    chain_id         numeric not null,
    schema_id        text not null,
    contract_address text not null,
    sync             jsonb
);
