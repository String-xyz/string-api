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

CREATE OR REPLACE TRIGGER update_contact_platform_updated_at
    BEFORE UPDATE
    ON contact_platform
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- DEVICE_INSTRUMENT ----------------------------------------------------
CREATE TABLE device_instrument (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  device_id UUID REFERENCES device (id),
  instrument_id UUID REFERENCES instrument (id)
);

CREATE OR REPLACE TRIGGER update_device_instrument_updated_at
    BEFORE UPDATE
    ON device_instrument
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- TX_LEG ---------------------------------------------------------------
CREATE TABLE tx_leg (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  timestamp TIMESTAMP WITH TIME ZONE,
  amount BIGINT DEFAULT 0, -- will need to define the base for USD and other fiat currencies
  value BIGINT DEFAULT 0, -- the relative quantity in USD at the time of the transaction
  asset UUID REFERENCES asset (id),
  user_id UUID REFERENCES string_user (id), -- this can be null in the case that the recipient is an unknown wallet address
  instrument_id UUID NOT NULL REFERENCES instrument (id)
);

CREATE OR REPLACE TRIGGER update_tx_leg_updated_at
    BEFORE UPDATE
    ON tx_leg
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
  transaction_hash TEXT DEFAULT '',
  network_id UUID NOT NULL REFERENCES network (id),
  network_fee BIGINT DEFAULT 0,
  parameters JSONB DEFAULT '{}'::JSONB, 
  contract_abi JSONB DEFAULT '{}'::JSONB,
  origin_tx_leg_id UUID NOT NULL REFERENCES tx_leg (id),
  receipt_tx_leg_id UUID NOT NULL REFERENCES tx_leg (id),
  response_tx_leg_id UUID NOT NULL REFERENCES tx_leg (id),
  destination_tx_leg_id UUID NOT NULL REFERENCES tx_leg (id),
  processing_fee BIGINT DEFAULT 0,
  processing_fee_asset UUID REFERENCES asset (id),
  string_fee BIGINT DEFAULT 0 -- // always USD?
);

CREATE OR REPLACE TRIGGER update_transaction_updated_at
    BEFORE UPDATE
    ON transaction
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- TRANSACTION ----------------------------------------------------------
DROP TRIGGER IF EXISTS update_transaction_updated_at ON transaction;
DROP TABLE transaction;

-------------------------------------------------------------------------
-- TX_LEG ---------------------------------------------------------------
DROP TRIGGER IF EXISTS update_tx_leg_updated_at ON tx_leg;
DROP TABLE tx_leg;

-------------------------------------------------------------------------
-- DEVICE_INSTRUMENT ----------------------------------------------------
DROP TRIGGER IF EXISTS update_device_instrument_updated_at ON device_instrument;
DROP TABLE device_instrument;

-------------------------------------------------------------------------
-- CONTACT_PLATFORM -----------------------------------------------------
DROP TRIGGER IF EXISTS update_contact_platform_updated_at ON contact_platform;
DROP TABLE contact_platform;