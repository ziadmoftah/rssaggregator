package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ziadmoftah/rssaggregator/internal/database"
)

func (apiCfg *apiConfig) handlerFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameter{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json: %v", err))
		return
	}
	feed_follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create feed follow: %v", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feed_follow))
}

func (apiCfg *apiConfig) handlerGetFollowedFeedsByUserId(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetFeedFollowsByUserId(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not feeds for selected user: %v", err))
		return
	}
	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feeds))
}

func (apiCfg *apiConfig) handlerDeleteFollowFeedByFeedId(w http.ResponseWriter, r *http.Request, user database.User) {
	feedIdStr := chi.URLParam(r, "feedId")
	feed_id, err := uuid.Parse(feedIdStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get the feed_id from url: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed_id,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not delete feed: %v", err))
		return
	}
	respondWithJSON(w, 204, struct{}{})
}
