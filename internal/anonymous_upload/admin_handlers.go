package anonymous_upload

import (
	"aaxion/internal/helpers"
	"net/http"
	"strconv"
)

// GenerateTokenHandler creates a new upload token (requires auth)
func GenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = helpers.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
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
		_ = helpers.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Build upload URL
	uploadURL := r.Host + "/upload?token=" + token

	_ = helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token":         token,
		"upload_url":    uploadURL,
		"target_dir":    targetDir,
		"max_uploads":   maxUploads,
		"expiry_hours":  expiryHours,
		"max_file_size": maxFileSize,
	})
}

// RevokeTokenHandler revokes a token (requires auth)
func RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		_ = helpers.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		_ = helpers.WriteError(w, http.StatusBadRequest, "Missing token")
		return
	}

	err := RevokeToken(token)
	if err != nil {
		_ = helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	_ = helpers.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Token revoked successfully",
	})
}

// ListTokensHandler lists all active tokens (requires auth)
func ListTokensHandler(w http.ResponseWriter, _ *http.Request) {
	tokens := ListAllTokens()

	_ = helpers.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"tokens": tokens,
		"count":  len(tokens),
	})
}

// GetTokenInfoHandler returns information about a specific token (requires auth)
func GetTokenInfoHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		_ = helpers.WriteError(w, http.StatusBadRequest, "Missing token")
		return
	}

	tokenInfo, err := GetTokenInfo(token)
	if err != nil {
		_ = helpers.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	_ = helpers.WriteJSON(w, http.StatusOK, tokenInfo)
}
