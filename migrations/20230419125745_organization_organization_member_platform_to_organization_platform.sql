-------------------------------------------------------------------------
-- +goose Up

-------------------------------------------------------------------------
-- ORGANIZATION ---------------------------------------------------------
-- +goose StatementBegin
CREATE TABLE organization (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  activated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  name TEXT NOT NULL,
  description TEXT DEFAULT ''
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE TRIGGER update_organization_updated_at
  BEFORE UPDATE
  ON organization
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();
-- +goose StatementEnd

-------------------------------------------------------------------------
-- ORGANIZATION-MEMBER --------------------------------------------------
-- +goose StatementBegin
ALTER TABLE platform_member
  RENAME TO organization_member;
-- +goose StatementEnd

-------------------------------------------------------------------------
-- MEMBER_TO_ORGANIZATION -----------------------------------------------
-- +goose StatementBegin
DROP TABLE IF EXISTS member_to_platform;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE member_to_organization (
  member_id UUID REFERENCES organization_member (id),
  organization_id UUID REFERENCES organization (id)
);
-- +goose StatementEnd


-------------------------------------------------------------------------
-- MEMBER_INVITE --------------------------------------------------------
-- +goose StatementBegin
ALTER TABLE member_invite
  DROP COLUMN IF EXISTS platform_id,
  ADD COLUMN organization_id UUID NOT NULL REFERENCES organization (id);
-- +goose StatementEnd


-------------------------------------------------------------------------
-- PLATFORM -------------------------------------------------------------
-- +goose StatementBegin
ALTER TABLE platform
  DROP COLUMN IF EXISTS activated_at,
  ADD COLUMN organization_id UUID NOT NULL REFERENCES organization (id);
-- +goose StatementEnd

-------------------------------------------------------------------------
-- APIKEY ---------------------------------------------------------------
ALTER TABLE apikey
  ADD COLUMN organization_id UUID NOT NULL REFERENCES organization (id);

-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- APIKEY ---------------------------------------------------------------
ALTER TABLE apikey
  DROP COLUMN IF EXISTS organization_id;

-------------------------------------------------------------------------
-- PLATFORM -------------------------------------------------------------
-- +goose StatementBegin
ALTER TABLE platform    
  DROP COLUMN IF EXISTS organization_id,
  ADD COLUMN activated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL;
-- +goose StatementEnd

-------------------------------------------------------------------------
-- ORGANIZATION-MEMBER --------------------------------------------------
-- +goose StatementBegin
ALTER TABLE organization_member
  RENAME TO platform_member;
-- +goose StatementEnd

-------------------------------------------------------------------------
-- MEMBER_INVITE --------------------------------------------------------
-- +goose StatementBegin
ALTER TABLE member_invite
  DROP COLUMN IF EXISTS organization_id,
  DROP COLUMN IF EXISTS organization_member,
  ADD COLUMN platform_member UUID REFERENCES platform_member (id),
  ADD COLUMN platform_id UUID REFERENCES platform (id);
-- +goose StatementEnd

-------------------------------------------------------------------------
-- MEMBER_TO_ORGANIZATION -----------------------------------------------
-- +goose StatementBegin
DROP TABLE IF EXISTS member_to_organization;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE member_to_platform (
  member_id UUID REFERENCES platform_member (id),
  platform_id UUID REFERENCES platform (id)
);
-- +goose StatementEnd

-------------------------------------------------------------------------
-- ORGANIZATION ---------------------------------------------------------
-- +goose StatementBegin
DROP TABLE organization;
-- +goose StatementEnd