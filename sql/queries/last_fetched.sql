-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
    ORDER BY last_fetched_at DESC
    LIMIT 10;


-- name: MarkFeedsFetched :exec
UPDATE feeds 
    SET (time_last_fetched, updated_at) = TIMESTAMP;