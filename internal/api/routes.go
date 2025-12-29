package api

import (
	"aaxion/internal/files"
	"net/http"
)

func RegisterRoutes() {

	http.HandleFunc("/files/view", files.ViewFiles)
}
