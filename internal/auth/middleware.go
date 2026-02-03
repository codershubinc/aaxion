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
		// Allow token in query param for static resources and streaming where headers are hard to set (e.g. img/video tags)
		if authQuery != "" && (strings.Contains(r.URL.String(), "/files/thumbnail") ||
			strings.Contains(r.URL.String(), "/files/download") ||
			strings.Contains(r.URL.String(), "/api/stream/movie") ||
			strings.Contains(r.URL.String(), "/files/view-image") ||
			strings.Contains(r.URL.String(), "/api/stream/episode")) {
			authHeader = "Bearer " + authQuery
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
	}
}
