-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id bigint,
    username text,
    first_name text,
    last_name text,
    PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE IF EXISTS users;