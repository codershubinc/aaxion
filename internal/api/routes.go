package api

import (
	"aaxion/internal/files"
	"net/http"
)

func RegisterRoutes() {

	http.HandleFunc("/files/view", files.ViewContent)
	http.HandleFunc("/files/create-directory", files.CreateDirectory)
	http.HandleFunc("/files/upload", files.UploadFile)
}
