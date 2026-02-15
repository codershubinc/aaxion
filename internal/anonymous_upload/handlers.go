package anonymous_upload

import (
	"aaxion/internal/files"
	"encoding/json"
	"fmt"
	"net/http"
)

// TokenUploadFile handles file uploads with token validation
func TokenUploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing upload token", http.StatusBadRequest)
		return
	}

	// Validate token
	uploadToken, err := ValidateToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Check path security
	isSuspicious := files.ExpelDotPath(uploadToken.TargetDir)
	if isSuspicious {
		http.Error(w, "Invalid target directory", http.StatusBadRequest)
		return
	}

	// Upload file
	err = files.UploadLargeFileToDir(w, r, uploadToken.TargetDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Increment upload count
	IncrementUploadCount(token)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Upload successful",
		"uploads_remaining": uploadToken.MaxUploads - uploadToken.UploadCount - 1,
	})
}

// TokenHandleStartChunkUpload initializes chunked upload with token
func TokenHandleStartChunkUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing upload token", http.StatusBadRequest)
		return
	}

	// Validate token
	_, err := ValidateToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	files.HandleStartChunkUpload(w, r)
}

// TokenHandleUploadChunk handles chunk upload with token
func TokenHandleUploadChunk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing upload token", http.StatusBadRequest)
		return
	}

	// Validate token
	_, err := ValidateToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	files.HandleUploadChunk(w, r)
}

// TokenHandleCompleteUpload completes chunked upload with token
func TokenHandleCompleteUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing upload token", http.StatusBadRequest)
		return
	}

	// Validate token
	uploadToken, err := ValidateToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Use token's target directory
	r.URL.RawQuery = fmt.Sprintf("filename=%s&dir=%s", filename, uploadToken.TargetDir)

	files.HandleCompleteUpload(w, r)

	// Increment upload count after successful completion
	IncrementUploadCount(token)
}

// ValidateTokenHandler checks if a token is valid (for frontend)
func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	uploadToken, err := ValidateToken(token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":             true,
		"target_dir":        uploadToken.TargetDir,
		"max_uploads":       uploadToken.MaxUploads,
		"uploads_remaining": uploadToken.MaxUploads - uploadToken.UploadCount,
		"expires_at":        uploadToken.ExpiresAt,
	})
}
