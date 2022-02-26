-- +goose Up
CREATE TABLE IF NOT EXISTS vocabulary (
    user_id bigint,
    word text,
    definition text,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE IF EXISTS vocabulary;