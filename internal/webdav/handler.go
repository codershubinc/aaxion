package webdav

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

func NewHandler(basePath string) http.Handler {
	return &webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdav.Dir(basePath),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				fmt.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL.Path, err)
			}
		},
	}
}

func GetRootPath() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}
