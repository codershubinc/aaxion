package music

import (
	"aaxion/internal/db"
	"aaxion/internal/models"
	"aaxion/internal/utils"
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
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	rawURI := r.FormValue("uri")
	if rawURI == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing 'uri' parameter")
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
		utils.WriteError(w, http.StatusInternalServerError, "Failed to process URI: "+err.Error())
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

	err = utils.WriteJSON(w, http.StatusAccepted, AddTrackResponse{
		Status:  "success",
		Message: fmt.Sprintf("Queued %d tracks for download", len(trackURLs)),
		Count:   len(trackURLs),
	})
	if err != nil {
		return
	}
}

func GetTracksApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	tracks, err := db.GetAllTracks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get tracks: "+err.Error())
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, tracks)
	if err != nil {
		return
	}
}

func SearchTracksApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	title := r.URL.Query().Get("title")
	if title == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing 'title' query parameter")
		return
	}

	tracks, err := db.SearchTracksByTitle(title)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to search tracks: "+err.Error())
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, tracks)
	if err != nil {
		return
	}
}

func GetTrackByIDApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "Missing 'id' query parameter")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid 'id' parameter: "+err.Error())
		return
	}

	track, err := db.GetTrackByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get track: "+err.Error())
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, track)
	if err != nil {
		return
	}
}

func UpdateTrackApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var trackData models.Track
	err := json.NewDecoder(r.Body).Decode(&trackData)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON body: "+err.Error())
		return
	}

	updatedTrack, err := db.UpdateTrack(trackData)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update track: "+err.Error())
		return
	}

	err = utils.WriteJSON(w, http.StatusOK, updatedTrack)
	if err != nil {
		return
	}
}
