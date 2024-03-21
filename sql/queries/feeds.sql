-- name: CreateFeed :one
INSERT INTO feeds(id, user_id, created_at, updated_at, name, url)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUserID :one
SELECT * FROM feeds 
    WHERE user_id = $1;
