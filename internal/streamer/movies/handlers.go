package movies

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

func StreamMovieApi(w http.ResponseWriter, r *http.Request) {
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

	movie, err := GetMovieByID(id)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	mimeType := movie.MimeType
	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(movie.FilePath))
		if mimeType == "" {
			mimeType = "video/mp4"
		}
	}

	streamer.StreamFileRange(w, r, movie.FilePath, mimeType)
}

func SearchMoviesApi(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing 'q' query parameter", http.StatusBadRequest)
		return
	}

	movies, err := SearchMovies(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SetJSONResponce(w, movies)
}

func ListMoviesApi(w http.ResponseWriter, r *http.Request) {
	movies, err := ListMovies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SetJSONResponce(w, movies)
}

func AddMovieApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Title       string `json:"title"`
		FileID      int64  `json:"file_id"`
		FilePath    string `json:"file_path"`
		Description string `json:"description"`
		PosterPath  string `json:"poster_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.FilePath == "" {
		http.Error(w, "FilePath is required", http.StatusBadRequest)
		return
	}

	// Determine size and mime type
	var size int64
	var mimeType string

	fileInfo, err := os.Stat(req.FilePath)
	if err == nil {
		size = fileInfo.Size()
	}

	// mime detection
	mimeType = mime.TypeByExtension(filepath.Ext(req.FilePath))
	if mimeType == "" {
		mimeType = "video/mp4"
	}

	err = AddMovie(req.Title, req.FileID, req.FilePath, req.Description, req.PosterPath, size, mimeType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func EditMovieApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID          int64  `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		PosterPath  string `json:"poster_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	// Optional: Fetch existing movie to ensure it exists or to merge fields if we wanted partial updates
	// For now, full update of provided fields
	// If title is empty, we probably shouldn't erase it? Let's check.
	// User said "edit info". Usually implies updating metadata.
	// Let's first get the movie to fill in blanks if user sends partial data, OR simply assume non-zero values are updates.
	// But `UpdateMovie` takes all 3 strings. If I pass empty string, it overwrites.
	// Safe bet for an "Edit" API is to read existing state, overlay changes, then save.

	existing, err := GetMovieByID(req.ID)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.PosterPath != "" {
		existing.PosterPath = req.PosterPath
	}

	err = UpdateMovie(req.ID, existing.Title, existing.Description, existing.PosterPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func helpGetFileName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
