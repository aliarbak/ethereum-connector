create table if not exists public.block
(
    hash     text   not null
    constraint block_pk
    primary key,
    chain_id numeric not null,
    number   numeric not null,
    time     numeric not null,
    created_at numeric
);