package files

import (
	"aaxion/internal/db"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func FileTempShare(w http.ResponseWriter, r *http.Request) {

	token := r.PathValue("token")

	if token == "" {
		http.Error(w, "Missing download token", http.StatusBadRequest)
		return
	}

	filePath, err := db.ValidateFileShareToken(token)
	if err != nil {
		http.Error(w, "Invalid or expired link", http.StatusForbidden)
		return
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	fileName := filepath.Base(filePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	http.ServeFile(w, r, filePath)
}

func RequestFileTempShare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := r.URL.Query().Get("file_path")
	if filePath == "" {
		http.Error(w, "Missing 'file_path' query parameter", http.StatusBadRequest)
		return
	}

	token, err := db.CreateFileShareTempToken(filePath)
	if err != nil {
		http.Error(w, "Failed to create share token", http.StatusInternalServerError)
		return
	}

	shareLink := fmt.Sprintf("/files/d/t/%s", token)
	w.Write([]byte(shareLink))
}
