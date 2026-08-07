package streamer

import (
	"aaxion/internal/files"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func StreamFileByPathApi(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "Missing 'path' query parameter", http.StatusBadRequest)
		return
	}

	if files.ExpelDotPath(filePath) {
		http.Error(w, "Suspicious path detected", http.StatusBadRequest)
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		http.Error(w, "Path points to a directory, not a file", http.StatusBadRequest)
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	StreamFileRange(w, r, filePath, mimeType)
}
