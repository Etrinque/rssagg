-- +goose Up
ALTER TABLE feeds
    ADD COLUMN last_fetched_at TIMESTAMP DEFAULT NULL;

-- +goose Down
FROM feeds DROP COLUMN last_fetched_at;

