create table if not exists public.schema
(
    id     text not null constraint schema_pk primary key,
    name   text not null,
    abi    jsonb not null,
    events jsonb not null
);