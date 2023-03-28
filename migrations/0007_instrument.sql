-- +goose Up

ALTER TABLE instrument
    ADD COLUMN name TEXT DEFAULT NULL,

-- +goose Down

ALTER TABLE instrument
    DROP COLUMN IF EXISTS name;