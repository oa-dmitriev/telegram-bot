-- +goose Up
ALTER TABLE vocabulary ADD COLUMN created_at timestamp with time zone DEFAULT now();

-- +goose Down
ALTER TABLE vocabulary DROP COLUMN created_at;