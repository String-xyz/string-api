-------------------------------------------------------------------------
-- +goose Up

-------------------------------------------------------------------------
-- CONTACT_PLATFORM -----------------------------------------------------
CREATE TABLE contact_platform (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  contact_id UUID REFERENCES contact (id),
  platform_id UUID REFERENCES platform (id)
);

CREATE TRIGGER update_contact_platform_updated_at
    BEFORE UPDATE
    ON contact_platform
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- BLOCKCHAIN_TX ----------------------------------------------------------
CREATE TABLE blockchain_tx (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  timestamp TIMESTAMP WITH TIME ZONE,
  status TEXT DEFAULT '', --enum
  tags JSONB DEFAULT '[]'::JSONB,
  transaction_hash TEXT DEFAULT '',
  network_id UUID NOT NULL REFERENCES network (id),
  network_fee BIGINT DEFAULT 0,
  parameters JSONB DEFAULT '{}'::JSONB,
  contract_abi JSONB DEFAULT '{}'::JSONB,
  sender_amount BIGINT DEFAULT 0,
  sender_asset UUID REFERENCES asset (id),
  sender_value BIGINT DEFAULT 0,
  sender_user_id UUID NOT NULL REFERENCES string_user (id),
  sender_instrument_id UUID NOT NULL REFERENCES instrument (id),
  reciever_amount BIGINT DEFAULT 0,
  reciever_asset UUID NOT NULL REFERENCES asset (id),
  reciever_value BIGINT DEFAULT 0,
  reciever_user_id UUID NOT NULL REFERENCES string_user (id),
  reciever_instrument_id UUID NOT NULL REFERENCES instrument (id)
);

CREATE TRIGGER update_blockchain_tx_updated_at
    BEFORE UPDATE
    ON blockchain_tx
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- TRANSACTION ----------------------------------------------------------
CREATE TABLE transaction (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  timestamp TIMESTAMP WITH TIME ZONE,
  type TEXT DEFAULT '', -- enum [fiat-to-crypto, crypto-to-fiat]
  status TEXT DEFAULT '', --enum
  tags JSONB DEFAULT '[]'::JSONB,
  device_id UUID REFERENCES device (id),
  ip_address TEXT DEFAULT '',
  platform_id UUID REFERENCES platform (id),
  sender_amount BIGINT DEFAULT 0, -- will need to define the base for USD and other fiat currencies
  sender_asset UUID REFERENCES asset (id),
  sender_value BIGINT DEFAULT 0, -- the relative quantity in USD at the time of the transaction
  sender_user_id UUID NOT NULL REFERENCES string_user (id),
  sender_instrument_id UUID NOT NULL REFERENCES instrument (id),
  reciever_amount BIGINT DEFAULT 0,
  reciever_asset UUID NOT NULL REFERENCES asset (id),
  reciever_value BIGINT DEFAULT 0,
  reciever_user_id UUID NOT NULL REFERENCES string_user (id),
  reciever_instrument_id UUID NOT NULL REFERENCES instrument (id),
  blockchain_tx_id UUID NOT NULL REFERENCES blockchain_tx (id),
  processing_fee BIGINT DEFAULT 0,
  processing_fee_asset UUID REFERENCES asset (id),
  string_fee BIGINT DEFAULT 0 -- // always USD?
);

CREATE TRIGGER update_transaction_updated_at
    BEFORE UPDATE
    ON transaction
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- TRANSACTION ----------------------------------------------------------
DROP TRIGGER IF EXISTS update_transaction_updated_at ON device;
DROP TABLE transaction;

-------------------------------------------------------------------------
-- TRANSACTION ----------------------------------------------------------
DROP TRIGGER IF EXISTS update_blockchain_tx_updated_at ON device;
DROP TABLE blockchain_tx;

-------------------------------------------------------------------------
-- CONTACT_PLATFORM -----------------------------------------------------
DROP TRIGGER IF EXISTS update_contact_platform_updated_at ON device;
DROP TABLE contact_platform;