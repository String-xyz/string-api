-------------------------------------------------------------------------
-- +goose Up

-------------------------------------------------------------------------
-- USER_TO_PLATFORM -----------------------------------------------------
ALTER TABLE user_platform
  RENAME TO user_to_platform;

DROP INDEX IF EXISTS user_platform_user_id_platform_id_idx;

CREATE UNIQUE INDEX user_to_platform_user_id_platform_id_idx ON user_to_platform(user_id, platform_id);


-------------------------------------------------------------------------
-- CONTACT_TO_PLATFORM --------------------------------------------------
ALTER TABLE contact_platform
  RENAME TO contact_to_platform;

DROP INDEX IF EXISTS contact_platform_contact_id_platform_id_idx;

CREATE UNIQUE INDEX contact_to_platform_contact_id_platform_id_idx ON contact_to_platform(contact_id, platform_id);


-------------------------------------------------------------------------
-- DEVICE_TO_INSTRUMENT -------------------------------------------------
ALTER TABLE device_instrument
  RENAME TO device_to_instrument;

DROP INDEX IF EXISTS device_instrument_device_id_instrument_id_idx;

CREATE UNIQUE INDEX device_to_instrument_device_id_instrument_id_idx ON device_to_instrument(device_id, instrument_id);


-------------------------------------------------------------------------
-- PLATFORM -------------------------------------------------------------
ALTER TABLE platform
  DROP COLUMN IF EXISTS type,
  DROP COLUMN IF EXISTS status, 
  DROP COLUMN IF EXISTS name, 
  DROP COLUMN IF EXISTS api_key, 
  DROP COLUMN IF EXISTS authentication,
  ADD COLUMN activated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,  -- for activating prod users
  ADD COLUMN name TEXT NOT NULL,
  ADD COLUMN description TEXT NOT NULL,
  ADD COLUMN domains TEXT[] DEFAULT NULL, -- define which domains can make calls to API (web-to-API)
  ADD COLUMN ip_addresses TEXT[] DEFAULT NULL; -- define which API ips can make calls (API-to-API)


-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- PLATFORM -------------------------------------------------------------
ALTER TABLE platform
  DROP COLUMN IF EXISTS activated_at,
  DROP COLUMN IF EXISTS name,
  DROP COLUMN IF EXISTS description,
  DROP COLUMN IF EXISTS domains,
  DROP COLUMN IF EXISTS ip_addresses,
  ADD COLUMN type TEXT NOT NULL, -- enum: to be defined at struct level in Go
  ADD COLUMN status TEXT NOT NULL, -- enum: to be defined at struct level in Go
  ADD COLUMN name TEXT DEFAULT '', 
  ADD COLUMN api_key TEXT DEFAULT '', 
  ADD COLUMN authentication TEXT DEFAULT ''; --enum [email, phone, wallet]


-------------------------------------------------------------------------
-- USER_PLATFORM -----------------------------------------------------
ALTER TABLE user_to_platform
  RENAME TO user_platform;

DROP INDEX IF EXISTS user_to_platform_user_id_platform_id_idx;

CREATE UNIQUE INDEX user_platform_user_id_platform_id_idx ON user_platform(user_id, platform_id);


-------------------------------------------------------------------------
-- CONTACT_PLATFORM --------------------------------------------------
ALTER TABLE contact_to_platform
  RENAME TO contact_platform;

DROP INDEX IF EXISTS contact_to_platform_contact_id_platform_id_idx;

CREATE UNIQUE INDEX contact_platform_contact_id_platform_id_idx ON contact_platform(contact_id, platform_id);


-------------------------------------------------------------------------
-- DEVICE_INSTRUMENT -------------------------------------------------
ALTER TABLE device_to_instrument
  RENAME TO device_instrument;

DROP INDEX IF EXISTS device_to_instrument_device_id_instrument_id_idx;

CREATE UNIQUE INDEX device_instrument_device_id_instrument_id_idx ON device_instrument(device_id, instrument_id);