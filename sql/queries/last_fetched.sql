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