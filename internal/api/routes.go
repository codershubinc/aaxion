package api

import (
	"aaxion/internal/files"
	img "aaxion/internal/image"
	sys "aaxion/internal/system"
	"net/http"
)

func RegisterRoutes() {

	// file management operations
	http.HandleFunc("/api/files/view", files.ViewContent)
	http.HandleFunc("/files/create-directory", files.CreateDirectory)

	// file upload  operations
	http.HandleFunc("/files/upload", files.UploadFile)

	// file upload  operations - chunked
	http.HandleFunc("/files/upload/chunk/start", files.HandleStartChunkUpload)
	http.HandleFunc("/files/upload/chunk/complete", files.HandleCompleteUpload)
	http.HandleFunc("/files/upload/chunk", files.HandleUploadChunk)

	// file download operations
	http.HandleFunc("/files/download", files.DownloadFileApi)
	http.HandleFunc("/files/thumbnail", img.ServeThumbnail)
	http.HandleFunc("/files/view-image", img.ViewImage)

	// temp files sharing
	http.HandleFunc("/files/d/t/{token}", files.FileTempShare)
	http.HandleFunc("/files/d/r", files.RequestFileTempShare)

	// system info
	http.HandleFunc("/api/system/get-root-path", sys.GetSystemRootPath)
	http.HandleFunc("/api/system/storage", sys.GetSystemStorage)

	// this is temp route to serve index.html for testing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
}
