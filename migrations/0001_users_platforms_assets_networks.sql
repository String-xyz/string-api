-------------------------------------------------------------------------
-- +goose Up

-------------------------------------------------------------------------
-- create extension for UUID --------------------------------------------
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-------------------------------------------------------------------------
-- UPDATE_UPDATED_AT_COLUMN() -------------------------------------------
-------------------------------------------------------------------------
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.update_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd
-------------------------------------------------------------------------

-------------------------------------------------------------------------
-- STRING_USER ----------------------------------------------------------
CREATE TABLE string_user (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  user_type TEXT NOT NULL, -- enum: to be defined at struct level in Go
  user_status TEXT NOT NULL, -- enum: to be defined at struct level in Go
  tags JSONB DEFAULT '{}'::JSONB,
  first_name TEXT DEFAULT '', -- name in separate table?
  middle_name TEXT DEFAULT '',
  last_name TEXT DEFAULT ''
);

CREATE TRIGGER update_string_user_updated_at
    BEFORE UPDATE
    ON string_user
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- PLATFORM -------------------------------------------------------------
CREATE TABLE platform (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  platform_type TEXT DEFAULT '', -- enum:
  api_key TEXT DEFAULT '',
  authentication_type TEXT DEFAULT '', --enum
  tags JSONB DEFAULT '[]'::JSONB
);
CREATE TRIGGER update_platform_updated_at
    BEFORE UPDATE
    ON platform
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- NETWORK --------------------------------------------------------------
CREATE TABLE network (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  network_id INT DEFAULT 0,
  chain_id INT NOT NULL,
  gas_token_id UUID DEFAULT NULL -- INDEX CREATED BELOW
);
CREATE TRIGGER update_network_updated_at
    BEFORE UPDATE
    ON network
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- ASSET ----------------------------------------------------------------
CREATE TABLE asset (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  description TEXT DEFAULT '',
  is_crypto BOOLEAN NOT NULL,
  network_id UUID REFERENCES network (id),
  value_oracle TEXT DEFAULT ''
);

CREATE TRIGGER update_asset_updated_at
    BEFORE UPDATE
    ON asset
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX network_gas_token_id_fk ON network (gas_token_id);


-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- ASSET ----------------------------------------------------------------
DROP TRIGGER IF EXISTS update_asset_updated_at ON asset;
DROP INDEX IF EXISTS network_gas_token_id_fk;
DROP TABLE asset;

-------------------------------------------------------------------------
-- NETWORK --------------------------------------------------------------
DROP TRIGGER IF EXISTS update_network_updated_at ON network;
DROP TABLE network;

-------------------------------------------------------------------------
-- PLATFORM -------------------------------------------------------------
DROP TRIGGER IF EXISTS update_platform_updated_at ON platfom;
DROP TABLE platform;

-------------------------------------------------------------------------
-- STRING_USER ----------------------------------------------------------
DROP TRIGGER IF EXISTS update_string_user_updated_at ON string_user;
DROP TABLE string_user;

-------------------------------------------------------------------------
-- UPDATE_UPDATED_AT_COLUMN() -------------------------------------------
DROP FUNCTION update_updated_at_column;

-------------------------------------------------------------------------
-- UUID EXTENSION -------------------------------------------------------
DROP EXTENSION IF EXISTS "uuid-ossp";