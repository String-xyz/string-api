-------------------------------------------------------------------------
-- +goose Up

-------------------------------------------------------------------------
-- PLATFORM -----------------------------------------------------
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
-- MEMBER -----------------------------------------------------
CREATE TABLE member (
	id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	email TEXT NOT NULL,
	password TEXT NOT NULL -- how do we maintain this?
);

-------------------------------------------------------------------------
-- PLATFORM_MEMBER -----------------------------------------------------
CREATE TABLE platform_member (
  platform_id UUID REFERENCES platform (id),
  member_id UUID REFERENCES member (id)
);

CREATE UNIQUE INDEX platform_member_platform_id_member_id_idx ON platform_member(platform_id, member_id);

-------------------------------------------------------------------------
-- ROLE -----------------------------------------------------
CREATE TABLE role (
	id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	name TEXT NOT NULL
);

-------------------------------------------------------------------------
-- MEMBER_ROLE -----------------------------------------------------
CREATE TABLE member_role (
    member_id UUID REFERENCES member (id),
	role_id UUID REFERENCES role (id)
);

CREATE UNIQUE INDEX member_role_member_id_role_id_idx ON member_role(member_id, role_id);

-------------------------------------------------------------------------
-- INVITE -----------------------------------------------------
CREATE TABLE invite (
	id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	expired_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	email TEXT NOT NULL,
	invited_by UUID REFERENCES member (id),
	platform_id UUID REFERENCES platform (id)
);

-------------------------------------------------------------------------
-- APIKEY -----------------------------------------------------
CREATE TABLE apikey (
	id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	type TEXT NOT NULL, -- [public,private] for now all public?
	data TEXT NOT NULL, -- the key itself
	description TEXT NOT NULL,
	created_by UUID REFERENCES member (id),
	platform_id UUID REFERENCES platform (id)
);


-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- PLATFORM -----------------------------------------------------
ALTER TABLE platform
	DROP COLUMN IF EXISTS activated_at,
	DROP COLUMN IF EXISTS name,
	DROP COLUMN IF EXISTS description,
	DROP COLUMN IF EXISTS domains,
	DROP COLUMN IF EXISTS ip_addresses
	ADD COLUMN type TEXT NOT NULL, -- enum: to be defined at struct level in Go
	ADD COLUMN status TEXT NOT NULL, -- enum: to be defined at struct level in Go
	ADD COLUMN name TEXT DEFAULT '', 
	ADD COLUMN api_key TEXT DEFAULT '', 
	ADD COLUMN authentication TEXT DEFAULT ''; --enum [email, phone, wallet]

-------------------------------------------------------------------------
-- MEMBER -----------------------------------------------------
DROP TABLE IF EXISTS member;

-------------------------------------------------------------------------
-- PLATFORM_MEMBER -----------------------------------------------------
DROP TABLE IF EXISTS platform_member;

-------------------------------------------------------------------------
-- ROLE -----------------------------------------------------
DROP TABLE IF EXISTS role;

-------------------------------------------------------------------------
-- MEMBER_ROLE -----------------------------------------------------
DROP TABLE IF EXISTS member_role;

-------------------------------------------------------------------------
-- INVITE -----------------------------------------------------
DROP TABLE IF EXISTS invite;

-------------------------------------------------------------------------
-- APIKEY -----------------------------------------------------
DROP TABLE IF EXISTS apikey;