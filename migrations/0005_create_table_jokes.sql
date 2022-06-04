-- +goose Up
CREATE TABLE IF NOT EXISTS jokes (
    id bigint,
    joke_type text,
    joke text,
    setup text,
    delivery text,
    category text,
    created_at timestamp with time zone DEFAULT now(),
    PRIMARY KEY(id)
);

CREATE INDEX IF NOT EXISTS jokes_category_idx ON jokes(category);

-- +goose Down
DROP INDEX IF EXISTS jokes_category_idx;

DROP TABLE IF EXISTS jokes;