package music

import (
	"aaxion/internal/db"
	"aaxion/internal/models"
	"aaxion/internal/ws"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type AddTrackResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Count   int    `json:"count"`
}

func AddTrackApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rawURI := r.FormValue("uri")
	if rawURI == "" {
		http.Error(w, "Missing 'uri' parameter", http.StatusBadRequest)
		return
	}

	uri, typeOfUri := CleanYouTubeURL(rawURI)
	var trackURLs []string
	var err error

	switch typeOfUri {
	case "video":
		trackURLs, err = []string{uri}, nil
	case "playlist":
		trackURLs, err = ManageDownloadYoutubePlaylist(uri)
	default:
		trackURLs, err = []string{}, fmt.Errorf("unsupported URI type: %s", typeOfUri)
	}

	if err != nil {
		http.Error(w, "Failed to process URI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	go func(urls []string) {
		for i, u := range urls {
			exactPath, err := DownloadYouTubeAudio(u, DIR)

			if err != nil {
				fmt.Printf("❌ Skip error: %v\n", err)
				ws.Broadcast(map[string]any{
					"type": "TRACK_ERROR",
					"payload": map[string]any{
						"message": "Failed to download track from YouTube",
						"url":     u,
					},
				})
				continue
			}

			trackData, err := ExtractYouTubeMetadata(exactPath)
			fmt.Println("yt meta ", trackData)
			if err != nil {
				fmt.Printf("❌ Metadata error: %v\n", err)
				continue
			}
			if trackData.YtUri == "" {
				trackData.YtUri = u
			}
			t := trackData
			err = db.AddTrack(t)
			if err != nil {
				fmt.Printf("❌ DB error: %v\n", err)
			}
			// 📢 Shout to the WebSocket suckers!
			ws.Broadcast(map[string]any{
				"type": "TRACK_ADDED",
				"state": map[string]any{
					"track":    trackData,
					"progress": fmt.Sprintf("%d/%d", i+1, len(urls)),
				},
			})

		}
	}(trackURLs)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(AddTrackResponse{
		Status:  "success",
		Message: fmt.Sprintf("Queued %d tracks for download", len(trackURLs)),
		Count:   len(trackURLs),
	})
}

func GetTracksApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tracks, err := db.GetAllTracks()
	if err != nil {
		http.Error(w, "Failed to get tracks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

func SearchTracksApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Missing 'title' query parameter", http.StatusBadRequest)
		return
	}

	tracks, err := db.SearchTracksByTitle(title)
	if err != nil {
		http.Error(w, "Failed to search tracks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
}

func GetTrackByIDApi(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(track)
}

func UpdateTrackApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var trackData models.Track
	err := json.NewDecoder(r.Body).Decode(&trackData)
	if err != nil {
		http.Error(w, "Invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	updatedTrack, err := db.UpdateTrack(trackData)
	if err != nil {
		http.Error(w, "Failed to update track: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTrack)
}
