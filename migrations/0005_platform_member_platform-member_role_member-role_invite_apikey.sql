-------------------------------------------------------------------------
-- +goose Up

-------------------------------------------------------------------------
-- PLATFORM -----------------------------------------------------
CREATE TABLE platform (
	id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	activated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL, -- for activating prod users
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	domains TEXT[] DEFAULT NULL, -- define which domains can make calls to API (web-to-API)
	ip_addresses TEXT[] DEFAULT NULL, -- define which API ips can make calls (API-to-API)
);

-------------------------------------------------------------------------
-- MEMBER -----------------------------------------------------
CREATE TABLE member (
	id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	email TEXT NOT NULL,
	password TEXT NOT NULL, -- how do we maintain this?
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
	name TEXT NOT NULL,
);

-------------------------------------------------------------------------
-- MEMBER_ROLE -----------------------------------------------------
CREATE TABLE member_role (
    member_id UUID REFERENCES member (id)
	role_id UUID REFERENCES role (id),
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
	platform_id UUID REFERENCES platform (id),
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
	platform_id UUID REFERENCES platform (id),
);


-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- CONTACT_PLATFORM -----------------------------------------------------
DROP TABLE IF EXISTS platform;

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