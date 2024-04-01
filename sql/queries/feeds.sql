-- name: CreateFeed :one
INSERT INTO feeds(id, user_id, created_at, updated_at, name, url)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUserID :one
SELECT * FROM feeds 
    WHERE user_id = $1;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
    ORDER BY last_fetched_at ASC NULLS FIRST
    LIMIT $1;


-- name: MarkFeedsFetched :one
UPDATE feeds 
    SET time_last_fetched = NOW(),
    updated_at = NOW()
    WHERE id = $1
    RETURNING *;