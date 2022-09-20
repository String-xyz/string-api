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
  origin_amount FLOAT, -- // string? should support both crypto and fiat
  origin_asset UUID REFERENCES assets (id),
  origin_value FLOAT,
  origin_user_id UUID NOT NULL REFERENCES users (id),
  origin_instrument_id UUID NOT NULL REFERENCES instruments (id),
  origin_user_location TEXT, 
  destination_amount FLOAT, -- // string? should support both crypto and fiat
  destination_asset UUID NOT NULL REFERENCES assets (id),
  destination_value FLOAT,
  destination_user_id UUID NOT NULL REFERENCES users (id),
  destination_instrument_id UUID NOT NULL REFERENCES instruments (id),
  transaction_hash TEXT DEFAULT '',
  network_id UUID NOT NULL REFERENCES networks (id),
  network_fee INT,
  network_fee_asset TEXT DEFAULT '',
  processing_fee FLOAT,
  processing_asset UUID REFERENCES assets (id),
  string_fee FLOAT -- // always USD?
);

-- +goose Down
DROP TABLE contact_platform;
DROP TABLE transactions;