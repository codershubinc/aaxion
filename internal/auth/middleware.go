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
		if (authQuery != "") && strings.Contains(r.URL.String(), "/files/thumbnail") || strings.Contains(r.URL.String(), "/files/download") {
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
