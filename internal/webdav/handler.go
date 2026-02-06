package webdav

import (
	"aaxion/internal/db"
	"fmt"
	"net/http"

	"os"

	"golang.org/x/net/webdav"
)

func NewHandler(basePath string) http.Handler {
	handler := &webdav.Handler{
		Prefix:     "/webdav",
		FileSystem: webdav.Dir(basePath),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				fmt.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL.Path, err)
			}
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow OPTIONS requests without authentication for service discovery
		if r.Method == "OPTIONS" {
			handler.ServeHTTP(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok || !db.VerifyCredentials(user, pass) {
			if user != "" {
				fmt.Printf("WebDAV Auth Failed: User='%s'\n", user)
			}
			w.Header().Set("WWW-Authenticate", `Basic realm="Aaxion"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func GetRootPath() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir
}
