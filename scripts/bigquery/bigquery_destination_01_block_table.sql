create table if not exists %s.block
(
    number         bignumeric not null,
    block_hash     string   not null,
    chain_id       numeric not null,
    time           numeric not null,
    PRIMARY KEY (block_hash) NOT ENFORCED
    );