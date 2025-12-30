package api

import (
	"aaxion/internal/files"
	"net/http"
)

func RegisterRoutes() {

	http.HandleFunc("/api/files/view", files.ViewContent)
	http.HandleFunc("/files/create-directory", files.CreateDirectory)
	http.HandleFunc("/files/upload", files.UploadFile)
	http.HandleFunc("/files/upload/chunk/start", files.HandleStartChunkUpload)
	http.HandleFunc("/files/upload/chunk/complete", files.HandleCompleteUpload)
	http.HandleFunc("/files/upload/chunk", files.HandleUploadChunk)

	// temp files sharing
	http.HandleFunc("/files/d/t/{token}", files.FileTempShare)
	http.HandleFunc("/files/d/r", files.RequestFileTempShare)

	// system info
	http.HandleFunc("/api/system/get-root-path", files.GetSystemRootPath)

	// this is temp route to serve index.html for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
}
