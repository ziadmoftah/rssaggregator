package main

import (
	"fmt"
	"net/http"

	"github.com/ziadmoftah/rssaggregator/internal/auth"
	"github.com/ziadmoftah/rssaggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth issues: %v", err))
			return
		}
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("Could not get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}
