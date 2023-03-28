-- +goose Up
ALTER TABLE instrument
    ADD COLUMN source_id TEXT NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE instrument
    DROP COLUMN IF EXISTS source_id;