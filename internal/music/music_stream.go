package music

import (
	"aaxion/internal/db"
	"net/http"
	"strconv"
)

func StreamTrackApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter: "+err.Error(), http.StatusBadRequest)
		return
	}

	track, err := db.GetTrackByID(id)
	if err != nil {
		http.Error(w, "Failed to get track: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, track.FilePath)
}
