-- +goose Up
ALTER TABLE users ADD COLUMN chat_id bigint;

-- +goose Down
ALTER TABLE users DROP COLUMN chat_id;
