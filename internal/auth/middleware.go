package auth

import (
	"aaxion/internal/db"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		authQuery := r.URL.Query().Get("tkn")

		if authQuery != "" {
			allowedPaths := []string{
				"/api/v1/images/thumbnail",
				"/api/v1/files/download",
				"/api/v1/movies/stream",
				"/api/v1/images/view",
				"/api/v1/series/episodes/stream",
				"/api/v1/files/stream",
			}

			for _, path := range allowedPaths {
				if strings.Contains(r.URL.Path, path) {
					authHeader = "Bearer " + authQuery
					break
				}
			}
		}

		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid auth header format", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		// For access token

		if strings.Contains(r.URL.String(), "/api/token/generate") {
			valid, err := db.VerifyToken(token)
			if err != nil {
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}

			if !valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			next(w, r)
			return
		}
		valid, err := db.VerifyAccessToken(token)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		if !valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
