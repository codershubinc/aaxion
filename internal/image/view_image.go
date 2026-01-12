package image

import (
	"aaxion/internal/files"
	"net/http"
	"os"
)

func ViewImage(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Missing 'path' query parameter", http.StatusBadRequest)
		return
	}

	if files.ExpelDotPath(filePath) {
		http.Error(w, "Suspicious path detected", http.StatusBadRequest)
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Add cache control header (cache for 7 days)
	w.Header().Set("Cache-Control", "public, max-age=604800")

	// Serve the file directly
	http.ServeFile(w, r, filePath)
}
