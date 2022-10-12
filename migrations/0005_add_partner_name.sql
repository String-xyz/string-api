-- +goose Up
ALTER TABLE platform
  ADD COLUMN partner_name TEXT;
  