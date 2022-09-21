-- +goose Up

-- create extension for UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
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

CREATE TABLE platforms (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  platform_type TEXT DEFAULT '', -- enum:
  api_key TEXT DEFAULT '',
  authentication_type TEXT DEFAULT '', --enum
  tags TEXT[]
);

CREATE TABLE networks (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  network_id INT DEFAULT 0,
  chain_id INT NOT NULL,
  gas_token_id UUID DEFAULT NULL -- CREATE REFERENCE IN SEPARATE MIGRATION
);

CREATE TABLE assets (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  description TEXT DEFAULT '',
  is_crypto BOOLEAN NOT NULL,
  network_id UUID REFERENCES networks (id),
  value_oracle TEXT DEFAULT ''
);

CREATE INDEX network_gas_token_id_fk ON networks (gas_token_id);

-- +goose Down
DROP TABLE users;
DROP TABLE platforms;
DROP TABLE assets;
DROP TABLE networks;