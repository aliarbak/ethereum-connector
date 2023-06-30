create unique index if not exists schema_contract_filter_chain_id_schema_id_contract_address_uind
    on public.schema_contract_filter (chain_id, schema_id, contract_address);

create index if not exists schema_contract_filter_schema_id_contract_address_index
    on public.schema_contract_filter (schema_id, contract_address);

create index if not exists schema_contract_filter_sync_completed_index
    on public.schema_contract_filter (((sync ->> 'completed')::bool));