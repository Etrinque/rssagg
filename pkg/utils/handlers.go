package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssagg/internal/database"
	"time"

	"github.com/google/uuid"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) authMiddleware(handler authedHandler) http.HandlerFunc { return func(w http.ResponseWriter, r *http.Request) {
	apiKey, err := GetApiToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "could not parse APIkey")
		return
	}
	user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "user not found")
		return
	}
	handler(w, r , user)
	}
}


func handleReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, "ok")
	return
}

func handleErr(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	return
}

// Un-Authenticated Endpoint
func (cfg *ApiConfig) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Could not decode json")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not create user")
		return 
	}
	RespondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

// Authenticated Endpoint
func (cfg *ApiConfig) handleUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {

	RespondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}

// Authenticated Endpooint
func (cfg *ApiConfig) handleFeedCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}

	decode := json.NewDecoder(r.Body)
	params := parameters{}
	err := decode.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error parsing json")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url: params.Url,
		UserID: user.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not create feed")
		return
	}

	feedFollow, err := cfg.DB.CreateFollowFeed(r.Context(), database.CreateFollowFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not create follow feed")
		return
	}

	RespondWithJSON(w, http.StatusOK, struct {
		feed Feed
		feedFollow FollowFeed
	}{
		databaseFeedToFeed(feed),
		databaseFeedFollowToFeedFollow(feedFollow),
	})
}

// Un-Authenticated Endpoint
func (cfg *ApiConfig) handleGetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "failure to retrieve feeds")
		return
	}

	RespondWithJSON(w, http.StatusOK, feeds)
}

// Authenticated Endpoint
func (cfg *ApiConfig) handleCreateFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID
	}

	decode := json.NewDecoder(r.Body)
	params := parameters{}
	err := decode.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not parse json")
		return
	}

	followFeed, err := cfg.DB.CreateFollowFeed(r.Context(), database.CreateFollowFeedParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID: user.ID,
			FeedID: params.FeedID,
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not create feed")
		return
	}

	RespondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(followFeed))
}

// Authenticated Endpoint
func (cfg *ApiConfig) handleGetFollowFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	followFeeds, err := cfg.DB.GetFollowFeed(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "Could not retrieve feed")
		return
	}

	RespondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(followFeeds))
}

// UN-Authenticated Endpoint
func (cfg *ApiConfig) handleDeleteFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	// feedid := chi.URLParam(r, "{feed_id}")
	feedid := r.PathValue("feed_id")
	feedID, err := uuid.Parse(feedid)
	fmt.Println(feedID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "record not found")
		return
	}
	fmt.Println("Feed_id cast to UUID: ", feedID)

	err = cfg.DB.DeleteFollowFeed(r.Context(), database.DeleteFollowFeedParams{
		ID: feedID,
		UserID: user.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not delete record")
		return
	}

	type resp struct{ msg string }
	RespondWithJSON(w, http.StatusOK, resp{ msg: "Record Deleted"})
}

// Delete helper function
// func parseDeleteEndpointParam(req *http.Request) (uuid.UUID, error) {
	// feedid := req.PathValue("feed_id")
	// if feedid == "" {
		// return uuid.Nil, errors.New("no provided feed_id")
	// }
	// fmt.Println("Feed_id string from url: ", feedid)
	// 
	// feedID, err := uuid.Parse(feedid)
	// if err != nil {
		// return uuid.Nil, err
	// }
	// fmt.Println("Feed_id cast to UUID: ", feedID)
// 
	// return feedID, nil
// }

