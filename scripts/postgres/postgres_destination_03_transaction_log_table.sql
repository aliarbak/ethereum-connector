create table if not exists transaction_log
(
    hash             text constraint transaction_log_pk primary key,
    index            numeric not null,
    transaction_hash text constraint transaction_log_transaction_hash_fk references transaction,
    contract_address text not null,
    event_name       text not null,
    event_alias      text not null,
    data             json not null
);

