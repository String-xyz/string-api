-------------------------------------------------------------------------
-- +goose Up

-------------------------------------------------------------------------
-- CONTRACT -------------------------------------------------------------
CREATE TABLE contract (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  name TEXT DEFAULT '',
  address TEXT NOT NULL,
  functions TEXT[] DEFAULT '{}'::TEXT[],
  network_id UUID NOT NULL REFERENCES network (id),
  platform_id UUID NOT NULL REFERENCES platform (id)
);

CREATE OR REPLACE TRIGGER update_contract_updated_at
  BEFORE UPDATE
  ON contract
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-------------------------------------------------------------------------
-- APIKEY ---------------------------------------------------------------
ALTER TABLE apikey 
  ADD COLUMN hint TEXT NOT NULL;


-------------------------------------------------------------------------
-- +goose Down

-------------------------------------------------------------------------
-- CONTRACT -------------------------------------------------------------
DROP TRIGGER IF EXISTS update_contract_updated_at ON contract;
DROP TABLE IF EXISTS contract;

-------------------------------------------------------------------------
-- APIKEY ---------------------------------------------------------------
ALTER TABLE apikey 
  DROP COLUMN IF EXISTS hint;
