package files

import (
	"net/http"
)

// PublicUploadFile handles file uploads without authentication
func PublicUploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	targetDir := r.URL.Query().Get("dir")
	if targetDir == "" {
		http.Error(w, "Missing 'dir' query parameter", http.StatusBadRequest)
		return
	}
	isSuspicious := ExpelDotPath(targetDir)
	if isSuspicious {
		http.Error(w, "Suspicious path detected", http.StatusBadRequest)
		return
	}

	err := UploadLargeFileToDir(w, r, targetDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// PublicHandleStartChunkUpload initializes chunked upload without authentication
func PublicHandleStartChunkUpload(w http.ResponseWriter, r *http.Request) {
	HandleStartChunkUpload(w, r)
}

// PublicHandleUploadChunk handles chunk upload without authentication
func PublicHandleUploadChunk(w http.ResponseWriter, r *http.Request) {
	HandleUploadChunk(w, r)
}

// PublicHandleCompleteUpload completes chunked upload without authentication
func PublicHandleCompleteUpload(w http.ResponseWriter, r *http.Request) {
	HandleCompleteUpload(w, r)
}
