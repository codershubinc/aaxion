package api

import (
	"aaxion/internal/auth"
	"aaxion/internal/files"
	img "aaxion/internal/image"
	sys "aaxion/internal/system"
	"net/http"
)

func RegisterRoutes() {

	// Auth routes
	http.HandleFunc("/auth/register", auth.Register)
	http.HandleFunc("/auth/login", auth.Login)

	// file management operations
	http.HandleFunc("/api/files/view", auth.AuthMiddleware(files.ViewContent))
	http.HandleFunc("/files/create-directory", auth.AuthMiddleware(files.CreateDirectory))

	// file upload  operations
	http.HandleFunc("/files/upload", auth.AuthMiddleware(files.UploadFile))

	// file upload  operations - chunked
	http.HandleFunc("/files/upload/chunk/start", auth.AuthMiddleware(files.HandleStartChunkUpload))
	http.HandleFunc("/files/upload/chunk/complete", auth.AuthMiddleware(files.HandleCompleteUpload))
	http.HandleFunc("/files/upload/chunk", auth.AuthMiddleware(files.HandleUploadChunk))

	// file download operations
	http.HandleFunc("/files/download", auth.AuthMiddleware(files.DownloadFileApi))
	http.HandleFunc("/files/thumbnail", auth.AuthMiddleware(img.ServeThumbnail))
	http.HandleFunc("/files/view-image", auth.AuthMiddleware(img.ViewImage))

	// temp files sharing
	http.HandleFunc("/files/d/t/{token}", files.FileTempShare)
	http.HandleFunc("/files/d/r", files.RequestFileTempShare)

	// system info
	http.HandleFunc("/api/system/get-root-path", auth.AuthMiddleware(sys.GetSystemRootPath))
	http.HandleFunc("/api/system/storage", auth.AuthMiddleware(sys.GetSystemStorage))

	// this is temp route to serve index.html for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
}
