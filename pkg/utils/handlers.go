package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssagg/internal/database"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *ApiConfig) authMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		handler(w, r, user)
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
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
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
		Url  string `json:"url"`
	}

	decode := json.NewDecoder(r.Body)
	params := parameters{}
	err := decode.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error parsing json")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not create feed")
		return
	}

	feedFollow, err := cfg.DB.CreateFollowFeed(r.Context(), database.CreateFollowFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not create follow feed")
		return
	}

	RespondWithJSON(w, http.StatusOK, struct {
		feed       Feed
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

	RespondWithJSON(w, http.StatusOK, databaseFeedsToFeeds(feeds))
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
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
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
	fmt.Println("starting parse of feed_id")
	feedid := chi.URLParam(r, "feed_id")
	fmt.Println(feedid)
	feedID, err := uuid.Parse(feedid)
	fmt.Println(feedID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "record not found")
		return
	}
	fmt.Println("Feed_id cast to UUID: ", feedID)

	err = cfg.DB.DeleteFollowFeed(r.Context(), database.DeleteFollowFeedParams{
		ID:     feedID,
		UserID: user.ID,
	})
	fmt.Println("record deleted")
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not delete record")
		return
	}

	type resp struct {
		MSG string `json:"msg"`
	}
	RespondWithJSON(w, http.StatusOK, resp{MSG: "Record Deleted"})
}

// func (cfg *ApiConfig) HandleGetFeedsFromUrl(w http.ResponseWriter, r *http.Request) {
// type parameters struct {
// RootURL string `json:"url"`
// }
// decoder := json.NewDecoder(r.Body)
// params := parameters{}
// err := decoder.Decode(&params)
// if err != nil {
// RespondWithError(w, http.StatusUnauthorized, "unauthorized user")
// return
// }
// rootUrl := params.RootURL
// req, err := http.Get(rootUrl)
// if err != nil {
// RespondWithError(w, http.StatusInternalServerError, "could not get root url")
// }
//
// fmt.Println(rootUrl)
// fmt.Println(req)
//
// resp, err := getFeedsFromUrl(params.RootURL)
// if err != nil {
// RespondWithError(w, http.StatusInternalServerError, "could not retrieve urls")
// return
// }
//
// fmt.Println(resp)
//
// RespondWithJSON(w, http.StatusOK, resp)
// }
//

// AUthenticated endpoint
func (cfg *ApiConfig) handleGetPostByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if specifiedLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = specifiedLimit
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "unable to retrieve posts")
		return
	}
	RespondWithJSON(w, http.StatusOK, databasePostToPosts(posts))
}
