-- +goose Up
CREATE TABLE IF NOT EXISTS vocabulary (
    user_id bigint,
    word text,
    definition text,
    UNIQUE(user_id, word),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id),
);

CREATE INDEX IF NOT EXISTS vocabulary_word_idx ON vocabulary(word);

-- +goose Down
DROP INDEX IF EXISTS vocabulary_word_idx;

DROP TABLE IF EXISTS vocabulary;