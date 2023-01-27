-- +goose Up
CREATE TABLE auth_strategy (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  type TEXT NOT NULL DEFAULT '', -- this auth strat is email/password type
  contact_id UUID REFERENCES contact(id) DEFAULT NULL,
  status TEXT NOT NULL DEFAULT 'pending',  -- temp use for API approval
  data TEXT NOT NULL, -- the hashed password
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE OR REPLACE TRIGGER update_auth_strategy_updated_at
  BEFORE UPDATE
  ON auth_strategy
  FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX auth_strategy_status_idx ON auth_strategy(status);
-- +goose Down
DROP INDEX IF EXISTS auth_strategy_status_idx;
DROP TRIGGER IF EXISTS update_auth_strategy_updated_at ON auth_strategy;
DROP TABLE IF EXISTS auth_strategy;