-- +goose Up
CREATE TABLE transaction (
    id int NOT NULL,
    created date NOT NULL,
    updated date NOT NULL,
    entity_type int, --foreign key: references entity_type (id)
    chain_id int,
    user_address text,
    contract_address text,
    contract_parameters text, -- can't take slice here
    contract_abi text, -- or can you... text[]
    contract_function text,
    tx_value text,
    gas_limit text,
    card_token text, -- all snake case
    tx_status int, -- enum foreign key (TransactionStatus)
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE post;