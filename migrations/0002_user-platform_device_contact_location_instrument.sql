-------------------------------------------------------------------------
-- +goose Up

-------------------------------------------------------------------------
-- USER_PLATFORM --------------------------------------------------------
CREATE TABLE user_platform (
  user_id UUID REFERENCES string_user (id),
  platform_id UUID REFERENCES platform (id)
);

CREATE UNIQUE INDEX user_platform_user_id_platform_id_idx ON user_platform(user_id, platform_id);

-------------------------------------------------------------------------
-- DEVICE ---------------------------------------------------------------
CREATE TABLE device (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_used_at TIMESTAMP WITH TIME ZONE NOT NULL,
  validated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  type TEXT DEFAULT '', -- enum: to be defined at struct level in Go
  description TEXT DEFAULT '',
  fingerprint TEXT DEFAULT '',
  ip_addresses JSONB DEFAULT '[]'::JSONB,
  user_id UUID NOT NULL REFERENCES string_user (id)
);

CREATE OR REPLACE TRIGGER update_device_updated_at
    BEFORE UPDATE
    ON device
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- CONTACT ---------------------------------------------------------------
CREATE TABLE contact (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_authenticated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  type TEXT NOT NULL, -- enum: [phone, email, etc...] to be defined at struct level in Go
  status TEXT DEFAULT '', -- enum: [primary, inactive] to be defined at struct level in Go
  data TEXT DEFAULT '', -- the contact information
  user_id UUID NOT NULL REFERENCES string_user (id)
);

CREATE OR REPLACE TRIGGER update_contact_updated_at
    BEFORE UPDATE
    ON contact
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- LOCATION ---------------------------------------------------------------
CREATE TABLE location (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  type TEXT DEFAULT '',
  status TEXT NOT NULL, -- enum: 
  tags JSONB DEFAULT '{}'::JSONB,
  building_number TEXT DEFAULT '',
  unit_number TEXT DEFAULT '',
  street_name TEXT DEFAULT '',
  city TEXT DEFAULT '',
  state TEXT DEFAULT '',
  postal_code TEXT DEFAULT '',
  country TEXT DEFAULT '' -- ISO 3166-1 standard
);

CREATE OR REPLACE TRIGGER update_location_updated_at
    BEFORE UPDATE
    ON location
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- INSTRUMENT ---------------------------------------------------------------
CREATE TABLE instrument (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  type TEXT NOT NULL, -- enum:  includes crypto wallet
  status TEXT NOT NULL, -- enum: 
  tags JSONB DEFAULT '{}'::JSONB,
  network TEXT NOT NULL, -- enum: 
  public_key TEXT DEFAULT '',
  last_4 TEXT DEFAULT '',
  user_id UUID REFERENCES string_user (id), -- instrument can be null in the circumstance that a user sends an asset to an unknown wallet
  location_id UUID REFERENCES location (id) DEFAULT NULL
);

CREATE OR REPLACE TRIGGER update_instrument_updated_at
    BEFORE UPDATE
    ON instrument
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- INSTRUMENT -----------------------------------------------------------
DROP TRIGGER IF EXISTS update_instrument_updated_at ON instrument;
DROP TABLE instrument;

-------------------------------------------------------------------------
-- LOCATION -----------------------------------------------------------
DROP TRIGGER IF EXISTS update_location_updated_at ON location;
DROP TABLE location;

-------------------------------------------------------------------------
-- CONTACT --------------------------------------------------------------
DROP TRIGGER IF EXISTS update_contact_updated_at ON contact;
DROP TABLE contact;

-------------------------------------------------------------------------
-- DEVICE ---------------------------------------------------------------
DROP TRIGGER IF EXISTS update_device_updated_at ON device;
DROP TABLE device;

-------------------------------------------------------------------------
-- CONTACT_PLATFORM -----------------------------------------------------
DROP TABLE contact_platform;