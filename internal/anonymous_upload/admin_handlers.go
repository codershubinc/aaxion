package anonymous_upload

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// GenerateTokenHandler creates a new upload token (requires auth)
func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse parameters
	targetDir := r.URL.Query().Get("target_dir")
	if targetDir == "" {
		targetDir = "/uploads"
	}

	maxUploadsStr := r.URL.Query().Get("max_uploads")
	maxUploads := 1 // default: one-time use
	if maxUploadsStr != "" {
		if parsed, err := strconv.Atoi(maxUploadsStr); err == nil && parsed > 0 {
			maxUploads = parsed
		}
	}

	expiryHoursStr := r.URL.Query().Get("expiry_hours")
	expiryHours := 24 // default: 24 hours
	if expiryHoursStr != "" {
		if parsed, err := strconv.Atoi(expiryHoursStr); err == nil && parsed > 0 {
			expiryHours = parsed
		}
	}

	maxFileSizeStr := r.URL.Query().Get("max_file_size")
	maxFileSize := int64(11 << 30) // default: 11GB
	if maxFileSizeStr != "" {
		if parsed, err := strconv.ParseInt(maxFileSizeStr, 10, 64); err == nil && parsed > 0 {
			maxFileSize = parsed
		}
	}

	// Generate token
	token, err := GenerateUploadToken(targetDir, maxUploads, expiryHours, maxFileSize)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Build upload URL
	uploadURL := r.Host + "/upload?token=" + token

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token":        token,
		"upload_url":   uploadURL,
		"target_dir":   targetDir,
		"max_uploads":  maxUploads,
		"expiry_hours": expiryHours,
		"max_file_size": maxFileSize,
	})
}

// RevokeTokenHandler revokes a token (requires auth)
func RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	err := RevokeToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Token revoked successfully",
	})
}

// ListTokensHandler lists all active tokens (requires auth)
func ListTokensHandler(w http.ResponseWriter, r *http.Request) {
	tokens := ListAllTokens()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tokens": tokens,
		"count":  len(tokens),
	})
}

// GetTokenInfoHandler returns information about a specific token (requires auth)
func GetTokenInfoHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	tokenInfo, err := GetTokenInfo(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenInfo)
}
