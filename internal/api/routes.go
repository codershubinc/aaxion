package api

import (
	"aaxion/internal/files"
	"net/http"
)

func RegisterRoutes() {

	http.HandleFunc("/files/view", files.ViewContent)
	http.HandleFunc("/files/create-directory", files.CreateDirectory)
	http.HandleFunc("/files/upload", files.UploadFile)
	http.HandleFunc("/files/upload/chunk/start", files.HandleStartChunkUpload)
	http.HandleFunc("/files/upload/chunk/complete", files.HandleCompleteUpload)
	http.HandleFunc("/files/upload/chunk", files.HandleUploadChunk)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
}
