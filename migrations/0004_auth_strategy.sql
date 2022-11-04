-- +goose Up
CREATE TABLE auth_strategy (
  id UUID PRIMARY KEY NOT NULL DEFAULT UUID_GENERATE_V4(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deactivated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  type TEXT NOT NULL DEFAULT '', -- this auth strat is email/password type
  contact_id UUID REFERENCES contact(id), 
  data TEXT NOT NULL -- the hashed password
);

CREATE OR REPLACE TRIGGER update_auth_strategy_updated_at
    BEFORE UPDATE
    ON auth_strategy
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- +goose Down
DROP TRIGGER IF EXISTS update_auth_strategy_updated_at ON auth_strategy;
DROP TABLE IF EXISTS auth_strategy;