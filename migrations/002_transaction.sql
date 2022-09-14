-- +goose Up
CREATE TABLE transaction (
    id UUID PRIMARY KEY  NOT NULL DEFAULT UUID_V4(),
    created date NOT NULL,
    updated date NOT NULL,
    entity_type int, --foreign key: references entity_type (id)
    chain_id INT,
    user_address TEXT,
    contract_address TEXT,
    contract_parameters TEXT, -- can't take slice here
    contract_abi TEXT, -- or can you... text[]
    contract_function TEXT,
    tx_value TEXT,
    gas_limit TEXT,
    card_token TEXT, -- all snake case
    tx_status INT, -- enum foreign key (TransactionStatus)
);

-- +goose Down
DROP TABLE post;