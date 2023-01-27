-------------------------------------------------------------------------
-- +goose Up

-------------------------------------------------------------------------
-- PLATFORM_MEMBER ------------------------------------------------------
CREATE TABLE platform_member (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL -- how do we maintain this?
);

-------------------------------------------------------------------------
-- PLATFORM_MEMBER ------------------------------------------------------
CREATE TABLE member_to_platform (
  member_id UUID REFERENCES platform_member (id),
  platform_id UUID REFERENCES platform (id)
);

CREATE UNIQUE INDEX member_to_platform_platform_id_member_id_idx ON member_to_platform(platform_id, member_id);

-------------------------------------------------------------------------
-- MEMBER_ROLE ----------------------------------------------------------
CREATE TABLE member_role (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  name TEXT NOT NULL
);

-------------------------------------------------------------------------
-- MEMBER_TO_ROLE -------------------------------------------------------
CREATE TABLE member_to_role (
  member_id UUID REFERENCES platform_member (id),
  role_id UUID REFERENCES role (id)
);

CREATE UNIQUE INDEX member_to_role_member_id_role_id_idx ON member_role(member_id, role_id);

-------------------------------------------------------------------------
-- MEMBER_INVITE --------------------------------------------------------
CREATE TABLE member_invite (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  expired_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  accepted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  email TEXT NOT NULL,
  invited_by UUID REFERENCES platform_member (id),
  platform_id UUID REFERENCES platform (id)
);

-------------------------------------------------------------------------
-- APIKEY ---------------------------------------------------------------
CREATE TABLE apikey (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  type TEXT NOT NULL, -- [public,private] for now all public?
  data TEXT NOT NULL, -- the key itself
  description TEXT NOT NULL,
  created_by UUID REFERENCES platform_member (id),
  platform_id UUID REFERENCES platform (id)
);


-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- PLATFORM_MEMBER ------------------------------------------------------
DROP TABLE IF EXISTS platform_member;

-------------------------------------------------------------------------
-- PLATFORM_TO_MEMBER ---------------------------------------------------
DROP TABLE IF EXISTS member_to_platform;

-------------------------------------------------------------------------
-- MEMBER_ROLE ----------------------------------------------------------
DROP TABLE IF EXISTS member_role;

-------------------------------------------------------------------------
-- MEMBER_TO_ROLE -------------------------------------------------------
DROP TABLE IF EXISTS member_to_role;

-------------------------------------------------------------------------
-- MEMBER_INVITE --------------------------------------------------------
DROP TABLE IF EXISTS member_invite;

-------------------------------------------------------------------------
-- APIKEY ---------------------------------------------------------------
DROP TABLE IF EXISTS apikey;