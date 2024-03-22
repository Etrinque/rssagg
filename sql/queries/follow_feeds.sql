-- name: CreateFollowFeed :one
INSERT INTO follow_feeds(id, feed_id, user_id, created_at, updated_at)
VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFollowFeed :many
SELECT * FROM follow_feeds WHERE user_id = $1;

-- name: GetFollowFeedsAll :many
SELECT * FROM follow_feeds;

-- name: DeleteFollowFeed :exec
DELETE FROM follow_feeds WHERE id = $1 and user_id = $2;