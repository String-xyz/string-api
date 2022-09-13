-- +goose Up
CREATE TABLE entity_type ( --naming convention?
    id int NOT NULL, -- use bigserial instead of int?
    created date NOT NULL, --timestamp(0) with time zone not null default now(),
    entity text, -- would prefer "name" but reserved word
    description text, -- reserved word
    PRIMARY KEY(id)
);

CREATE TABLE tx_status ( --naming convention?
    id int NOT NULL,
    created date NOT NULL,
    tx_status text, -- would prefer "name" but reserved word
    description text, -- reserved word
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE entity_type;
DROP TABLE tx_status;