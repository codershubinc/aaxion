package stats

import (
	"aaxion/internal/db"
	"aaxion/internal/utils"
	"net/http"
	"strconv"
	"time"
)

// Helper to get user_id from query (defaults to 1 for now if no auth is present)
func getUserId(r *http.Request) int {
	uidStr := r.URL.Query().Get("user_id")
	if uidStr != "" {
		if uid, err := strconv.Atoi(uidStr); err == nil {
			return uid
		}
	}
	return 1
}

// RecordPlayApi handles recording a track play
func RecordPlayApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	trackIdStr := r.FormValue("track_id")
	if trackIdStr == "" {
		trackIdStr = r.URL.Query().Get("track_id")
	}

	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		_ = utils.WriteError(w, http.StatusBadRequest, "Invalid or missing 'track_id'")
		return
	}

	userId := getUserId(r)
	playedAt := time.Now().Format(time.RFC3339)

	db.RecordPlay(trackId, userId, playedAt)

	_ = utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Play recorded successfully",
	})
}

// SetFavoriteApi handles favoriting/unfavoriting a track
func SetFavoriteApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	trackIdStr := r.FormValue("track_id")
	if trackIdStr == "" {
		trackIdStr = r.URL.Query().Get("track_id")
	}

	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		_ = utils.WriteError(w, http.StatusBadRequest, "Invalid or missing 'track_id'")
		return
	}

	isFavoriteStr := r.FormValue("is_favorite")
	if isFavoriteStr == "" {
		isFavoriteStr = r.URL.Query().Get("is_favorite")
	}

	isFavorite := isFavoriteStr == "true" || isFavoriteStr == "1"
	userId := getUserId(r)

	err = db.SetFavorite(trackId, userId, isFavorite)
	if err != nil {
		_ = utils.WriteError(w, http.StatusInternalServerError, "Failed to update favorite status: "+err.Error())
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Favorite status updated",
	})
}

// GetPlayStateApi gets the play state for a single track
func GetPlayStateApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		_ = utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	trackIdStr := r.URL.Query().Get("track_id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		_ = utils.WriteError(w, http.StatusBadRequest, "Invalid or missing 'track_id'")
		return
	}

	userId := getUserId(r)

	state, err := db.GetPlayState(trackId, userId)
	if err != nil {
		_ = utils.WriteError(w, http.StatusInternalServerError, "Failed to get play state: "+err.Error())
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, state)
}

// GetAllPlayStatesApi gets all play states for the current user
func GetAllPlayStatesApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		_ = utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userId := getUserId(r)

	states, err := db.GetAllPlayStates(userId)
	if err != nil {
		_ = utils.WriteError(w, http.StatusInternalServerError, "Failed to get play states: "+err.Error())
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, states)
}

// CheckFavoriteApi checks if a track is favorited by the user
func CheckFavoriteApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		_ = utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	trackIdStr := r.URL.Query().Get("track_id")
	trackId, err := strconv.Atoi(trackIdStr)
	if err != nil {
		_ = utils.WriteError(w, http.StatusBadRequest, "Invalid or missing 'track_id'")
		return
	}

	userId := getUserId(r)
	isFavorite := db.CheckFavorite(trackId, userId)

	_ = utils.WriteJSON(w, http.StatusOK, map[string]bool{
		"is_favorite": isFavorite,
	})
}

// GetLastPlayedApi returns the user's latest played track
func GetLastPlayedApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		_ = utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userId := getUserId(r)

	state, err := db.GetUserLastPlayedTrack(userId)
	if err != nil {
		_ = utils.WriteError(w, http.StatusNotFound, "No recent tracks found")
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, state)
}
