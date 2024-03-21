package utils

import (
	"rssagg/internal/database"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name string `json:"name"`
}

func databaseUserToUser(user database.User) User {
return User{
	ID: user.ID,
	CreatedAt: user.CreatedAt,
	UpdatedAt: user.UpdatedAt,
	Name: user.Name,
}
}

type Feed struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name string `json:"name"`
	Url string `json:"url"`
	UserID uuid.UUID `json:"user_id"`
}

func databaseFeedToFeed(feed database.Feed) Feed {
	return Feed{
		ID: feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name: feed.Name,
		Url:  feed.Url,
		UserID: feed.UserID,
	}
}

func databaseFeedsToFeeds(feeds []database.Feed) []Feed {
	result := make([]Feed, len(feeds))
	for i, feed := range feeds {
		result[i] = databaseFeedToFeed(feed)
	}
	return result
}


type FollowFeed struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
}

func databaseFeedFollowToFeedFollow(followFeed database.FollowFeed) FollowFeed {
		return FollowFeed{
			ID: followFeed.ID,
			FeedID: followFeed.FeedID,
			UserID: followFeed.UserID,
			CreatedAt: followFeed.CreatedAt,
			UpdatedAt: followFeed.UpdatedAt,

		}
}

func databaseFeedFollowsToFeedFollows(feedFollows []database.FollowFeed) []FollowFeed {
	result := make([]FollowFeed, len(feedFollows))
	for i, feedFollow := range feedFollows {
		result[i] = databaseFeedFollowToFeedFollow(feedFollow)
	}
	return result
}