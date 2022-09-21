-- +goose Up
CREATE TABLE contact_platform (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  contact_id UUID REFERENCES contacts (id),
  platform_id UUID REFERENCES platforms (id)
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  txn_type TEXT DEFAULT '', -- enum
  txn_timestamp TIMESTAMP WITH TIME ZONE,
  txn_status TEXT DEFAULT '', --enum
  tags JSONB DEFAULT '[]'::JSONB,
  smart_contract_params JSONB DEFAULT '{}'::JSONB,
  platform_id UUID REFERENCES platforms (id),
  origin_amount BIGINT DEFAULT 0, -- will need to define the base for USD and other fiat currencies
  origin_asset UUID REFERENCES assets (id),
  origin_value BIGINT DEFAULT 0,
  origin_user_id UUID NOT NULL REFERENCES users (id),
  origin_instrument_id UUID NOT NULL REFERENCES instruments (id),
  origin_user_location TEXT DEFAULT '', 
  destination_amount BIGINT DEFAULT 0, -- // string? should support both crypto and fiat
  destination_asset UUID NOT NULL REFERENCES assets (id),
  destination_value BIGINT DEFAULT 0,
  destination_user_id UUID NOT NULL REFERENCES users (id),
  destination_instrument_id UUID NOT NULL REFERENCES instruments (id),
  transaction_hash TEXT DEFAULT '',
  network_id UUID NOT NULL REFERENCES networks (id),
  network_fee BIGINT DEFAULT 0,
  network_fee_asset TEXT DEFAULT '',
  processing_fee BIGINT DEFAULT 0,
  processing_asset UUID REFERENCES assets (id),
  string_fee BIGINT DEFAULT 0 -- // always USD?
);

-- +goose Down
DROP TABLE contact_platform;
DROP TABLE transactions;