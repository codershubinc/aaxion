package series

import (
	"aaxion/internal/helpers"
	"aaxion/internal/streamer"
	"encoding/json"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// Series Handlers

func ListSeriesApi(w http.ResponseWriter, r *http.Request) {
	seriesList, err := ListSeries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.SetJSONResponce(w, seriesList)
}

func SearchSeriesApi(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing 'q' query parameter", http.StatusBadRequest)
		return
	}
	seriesList, err := SearchSeries(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.SetJSONResponce(w, seriesList)
}

func AddSeriesApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	err := AddSeries(req.Title, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func EditSeriesApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID          int64  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	existing, err := GetSeriesByID(req.ID)
	if err != nil {
		http.Error(w, "Series not found", http.StatusNotFound)
		return
	}
	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	err = UpdateSeries(req.ID, existing.Title, existing.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Episode Handlers

func AddEpisodeApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		SeriesID      int64  `json:"series_id"`
		FileID        int64  `json:"file_id"`
		SeasonNumber  int    `json:"season_number"`
		EpisodeNumber int    `json:"episode_number"`
		Title         string `json:"title"`
		Description   string `json:"description"`
		FilePath      string `json:"file_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.SeriesID == 0 || req.FileID == 0 || req.FilePath == "" {
		http.Error(w, "SeriesID, FileID, and FilePath are required", http.StatusBadRequest)
		return
	}

	// Determine size and mime type
	var size int64
	var mimeType string

	fileInfo, err := os.Stat(req.FilePath)
	if err == nil {
		size = fileInfo.Size()
	}

	mimeType = mime.TypeByExtension(filepath.Ext(req.FilePath))
	if mimeType == "" {
		mimeType = "video/mp4"
	}

	err = AddEpisode(req.SeriesID, req.FileID, req.SeasonNumber, req.EpisodeNumber, req.Title, req.Description, req.FilePath, size, mimeType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func ListEpisodesApi(w http.ResponseWriter, r *http.Request) {
	seriesIDStr := r.URL.Query().Get("series_id")
	if seriesIDStr == "" {
		http.Error(w, "Missing 'series_id' query parameter", http.StatusBadRequest)
		return
	}
	seriesID, err := strconv.ParseInt(seriesIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid series ID format", http.StatusBadRequest)
		return
	}
	episodes, err := ListEpisodes(seriesID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.SetJSONResponce(w, episodes)
}

func StreamEpisodeApi(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	episode, err := GetEpisodeByID(id)
	if err != nil {
		http.Error(w, "Episode not found", http.StatusNotFound)
		return
	}
	mimeType := episode.MimeType
	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(episode.FilePath))
		if mimeType == "" {
			mimeType = "video/mp4"
		}
	}
	streamer.StreamFileRange(w, r, episode.FilePath, mimeType)
}
