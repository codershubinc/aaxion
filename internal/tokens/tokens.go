package tokens

import (
	"aaxion/internal/db"
	"encoding/json"
	"net/http"
	"time"
)

func CreateAccessToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token, err := db.CreateAccessToken("access", time.Now().Add(5*time.Hour))
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"token": "` + token + `"}`))

}
func CleanAccessTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type RequestPayload struct {
		Token string `json:"token"`
	}
	var pl RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&pl); err != nil || pl.Token == "" {
		http.Error(w, "Invalid request body or missing token", http.StatusBadRequest)
		return
	}

	err := db.InvalidateAccessToken(pl.Token)
	if err != nil {
		http.Error(w, "Failed to clean access tokens", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Access tokens cleaned"}`))
}
