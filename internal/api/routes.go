package api

import (
	"aaxion/internal/files"
	"aaxion/internal/system"
	"net/http"
)

func RegisterRoutes() {

	http.HandleFunc("/api/files/view", files.ViewContent)
	http.HandleFunc("/files/create-directory", files.CreateDirectory)
	http.HandleFunc("/files/upload", files.UploadFile)
	http.HandleFunc("/files/upload/chunk/start", files.HandleStartChunkUpload)
	http.HandleFunc("/files/upload/chunk/complete", files.HandleCompleteUpload)
	http.HandleFunc("/files/upload/chunk", files.HandleUploadChunk)
	http.HandleFunc("/files/download", files.DownloadFileApi)

	// temp files sharing
	http.HandleFunc("/files/d/t/{token}", files.FileTempShare)
	http.HandleFunc("/files/d/r", files.RequestFileTempShare)

	// system info
	http.HandleFunc("/api/system/get-root-path", system.GetSystemRootPath)
	http.HandleFunc("/api/system/storage", system.GetSystemStorage)

	// this is temp route to serve index.html for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
}
