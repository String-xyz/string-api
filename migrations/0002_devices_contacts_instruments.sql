-- +goose Up
CREATE TABLE devices (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_used_at TIMESTAMP WITH TIME ZONE NOT NULL,
  validated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  device_type TEXT DEFAULT '', -- enum: to be defined at struct level in Go
  device_description TEXT DEFAULT '',
  user_id UUID NOT NULL REFERENCES users (id)
);

CREATE TABLE contacts (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_authenticated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  contact_type TEXT NOT NULL, -- enum: [phone, email, etc...] to be defined at struct level in Go
  contact_status TEXT DEFAULT '', -- enum: [primary, inactive] to be defined at struct level in Go
  user_id UUID NOT NULL REFERENCES users (id)
);

CREATE TABLE instruments (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  instrument_type TEXT NOT NULL, -- enum:  includes crypto wallet
  instrument_status TEXT NOT NULL, -- enum: 
  instrument_network TEXT NOT NULL, -- enum: 
  public_key TEXT DEFAULT '',
  last_4 TEXT DEFAULT '',
  tags JSONB DEFAULT '[]'::JSONB,
  user_id UUID NOT NULL REFERENCES users (id),
  location_type TEXT DEFAULT '',
  address_number TEXT DEFAULT '',
  unit_number TEXT DEFAULT '',
  street_name TEXT DEFAULT '',
  city TEXT DEFAULT '',
  state TEXT DEFAULT '',
  postal_code TEXT DEFAULT '',
  country TEXT DEFAULT '' -- ISO 3166-1 standard
);

-- +goose Down
DROP TABLE devices;
DROP TABLE contacts;
DROP TABLE instruments;