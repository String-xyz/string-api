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
  origin_amount BIGINT DEFAULT 0, -- will need to define the base for USD and other fiat currencies
  origin_asset UUID REFERENCES asset (id),
  origin_value BIGINT DEFAULT 0, -- the relative quantity in USD at the time of the transaction
  origin_user_id UUID NOT NULL REFERENCES string_user (id),
  origin_instrument_id UUID NOT NULL REFERENCES instrument (id),
  receipt_amount BIGINT DEFAULT 0,
  receipt_asset UUID REFERENCES asset (id),
  receipt_value BIGINT DEFAULT 0,
  receipt_user_id UUID NOT NULL REFERENCES string_user (id),
  receipt_instrument_id UUID NOT NULL REFERENCES instrument (id),
  response_amount BIGINT DEFAULT 0,
  response_asset UUID NOT NULL REFERENCES asset (id),
  response_value BIGINT DEFAULT 0,
  response_user_id UUID NOT NULL REFERENCES string_user (id),
  response_instrument_id UUID NOT NULL REFERENCES instrument (id),
  destination_amount BIGINT DEFAULT 0,
  destination_asset UUID NOT NULL REFERENCES asset (id),
  destination_value BIGINT DEFAULT 0,
  destination_user_id UUID NOT NULL REFERENCES string_user (id),
  destination_instrument_id UUID NOT NULL REFERENCES instrument (id),
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
-- DEVICE_INSTRUMENT ----------------------------------------------------
DROP TRIGGER IF EXISTS update_device_instrument_updated_at ON device_instrument;
DROP TABLE device_instrument;

-------------------------------------------------------------------------
-- CONTACT_PLATFORM -----------------------------------------------------
DROP TRIGGER IF EXISTS update_contact_platform_updated_at ON contact_platform;
DROP TABLE contact_platform;