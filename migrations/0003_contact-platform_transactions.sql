-- +goose Up
CREATE TABLE contact_platform (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  contact_id UUID REFERENCES contact (id),
  platform_id UUID REFERENCES platform (id)
);

CREATE TABLE transaction (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  txn_type TEXT DEFAULT '', -- enum
  txn_timestamp TIMESTAMP WITH TIME ZONE,
  txn_status TEXT DEFAULT '', --enum
  tags JSONB DEFAULT '[]'::JSONB,
  smart_contract_params JSONB DEFAULT '{}'::JSONB,
  platform_id UUID REFERENCES platform (id),
  origin_amount BIGINT DEFAULT 0, -- will need to define the base for USD and other fiat currencies
  origin_asset UUID REFERENCES asset (id),
  origin_value BIGINT DEFAULT 0,
  origin_user_id UUID NOT NULL REFERENCES string_user (id),
  origin_instrument_id UUID NOT NULL REFERENCES instrument (id),
  origin_user_location TEXT DEFAULT '', 
  destination_amount BIGINT DEFAULT 0, -- // string? should support both crypto and fiat
  destination_asset UUID NOT NULL REFERENCES asset (id),
  destination_value BIGINT DEFAULT 0,
  destination_user_id UUID NOT NULL REFERENCES string_user (id),
  destination_instrument_id UUID NOT NULL REFERENCES instrument (id),
  transaction_hash TEXT DEFAULT '',
  network_id UUID NOT NULL REFERENCES network (id),
  network_fee BIGINT DEFAULT 0,
  network_fee_asset TEXT DEFAULT '',
  processing_fee BIGINT DEFAULT 0,
  processing_asset UUID REFERENCES asset (id),
  string_fee BIGINT DEFAULT 0 -- // always USD?
);

-- +goose Down
DROP TABLE contact_platform;
DROP TABLE transaction;