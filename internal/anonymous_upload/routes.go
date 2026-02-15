package anonymous_upload

import (
	"aaxion/internal/auth"
	"net/http"
)

// RegisterRoutes registers all anonymous upload routes
func RegisterRoutes() {
	// Token-based upload routes (no auth required, but valid token needed)
	http.HandleFunc("/upload/token/file", TokenUploadFile)
	http.HandleFunc("/upload/token/chunk/start", TokenHandleStartChunkUpload)
	http.HandleFunc("/upload/token/chunk", TokenHandleUploadChunk)
	http.HandleFunc("/upload/token/chunk/complete", TokenHandleCompleteUpload)
	http.HandleFunc("/upload/token/validate", ValidateTokenHandler)

	// Admin routes for token management (auth required)
	http.HandleFunc("/api/upload-tokens/generate", auth.AuthMiddleware(GenerateTokenHandler))
	http.HandleFunc("/api/upload-tokens/revoke", auth.AuthMiddleware(RevokeTokenHandler))
	http.HandleFunc("/api/upload-tokens/list", auth.AuthMiddleware(ListTokensHandler))
	http.HandleFunc("/api/upload-tokens/info", auth.AuthMiddleware(GetTokenInfoHandler))
}
