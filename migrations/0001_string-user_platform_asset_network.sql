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
    NEW.updated_at = now();
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
  type TEXT NOT NULL, -- enum: to be defined at struct level in Go
  status TEXT NOT NULL, -- enum: to be defined at struct level in Go
  tags JSONB DEFAULT '{}'::JSONB, -- platforms should be listed in the tags
  first_name TEXT DEFAULT '', -- name in separate table?
  middle_name TEXT DEFAULT '',
  last_name TEXT DEFAULT ''
);

CREATE OR REPLACE TRIGGER update_string_user_updated_at
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
  type TEXT NOT NULL, -- enum: to be defined at struct level in Go
  status TEXT NOT NULL, -- enum: to be defined at struct level in Go
  name TEXT DEFAULT '',
  api_key TEXT DEFAULT '',
  authentication TEXT DEFAULT '' --enum [email, phone, wallet]
);
CREATE OR REPLACE TRIGGER update_platform_updated_at
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
  network_id TEXT DEFAULT '', -- might actually be big.Int
  chain_id TEXT NOT NULL, -- might actually be big.Int
  gas_token_id TEXT DEFAULT '', -- INDEX CREATED BELOW
  gas_oracle TEXT DEFAULT '', -- the name of the network in oracle (i.e. in owlracle)
  rpc_url TEXT DEFAULT '', -- The RPC used to access the network (ie "https://mainnet.infura.io/v3")
  explorer_url TEXT DEFAULT '' -- The Block Explorer URL used to view transactions and entities in the browser
);
CREATE OR REPLACE TRIGGER update_network_updated_at
    BEFORE UPDATE
    ON network
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- ASSET ----------------------------------------------------------------
CREATE TABLE asset ( -- We will write sql commands to add/update these in bulk.
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  description TEXT DEFAULT '',
  decimals INT DEFAULT 0,
  is_crypto BOOLEAN NOT NULL,
  network_id UUID REFERENCES network (id) DEFAULT NULL,
  value_oracle TEXT DEFAULT '' -- the name of the asset in oracle (i.e. in coingecko).  
);

CREATE OR REPLACE TRIGGER update_asset_updated_at
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
DROP TABLE IF EXISTS asset;

-------------------------------------------------------------------------
-- NETWORK --------------------------------------------------------------
DROP TRIGGER IF EXISTS update_network_updated_at ON network;
DROP TABLE IF EXISTS network;

-------------------------------------------------------------------------
-- PLATFORM -------------------------------------------------------------
DROP TRIGGER IF EXISTS update_platform_updated_at ON platfom;
DROP TABLE IF EXISTS platform;

-------------------------------------------------------------------------
-- STRING_USER ----------------------------------------------------------
DROP TRIGGER IF EXISTS update_string_user_updated_at ON string_user;
DROP TABLE IF EXISTS string_user;

-------------------------------------------------------------------------
-- UPDATE_UPDATED_AT_COLUMN() -------------------------------------------
DROP FUNCTION IF EXISTS update_updated_at_column;

-------------------------------------------------------------------------
-- UUID EXTENSION -------------------------------------------------------
DROP EXTENSION IF EXISTS "uuid-ossp";