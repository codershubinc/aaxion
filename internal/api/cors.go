package api

import (
	"aaxion/internal/utils"
	"fmt"
	"net/http"
	"strings"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) headers for incoming requests.
func CORSMiddleware(next http.Handler) http.Handler {
	// Get all local IPs
	localIPs, err := utils.GetAllLocalIPs()
	if err != nil {
		fmt.Printf("Error getting local IPs: %v\n", err)
	}
	// Add common localhost and tunnel domain variants
	localIPs = append(localIPs, "localhost", "127.0.0.1")
	localIPs = append(localIPs, "aaxion-client.codershubinc.com", "aaxioncdn.codershubinc.com", "codershubinc.com", "aaxion-cdn.codershubinc.tech")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow requests from any local IP or registered domain
		if origin != "" {
			for _, ip := range localIPs {
				if strings.Contains(origin, ip) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}
			if w.Header().Get("Access-Control-Allow-Origin") == "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PROPFIND, PROPPATCH, MKCOL, COPY, MOVE, LOCK, UNLOCK")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Depth, Destination, If, Overwrite, Timeout, Range")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges, Content-Type, Content-Disposition")

		if strings.HasPrefix(r.URL.Path, "/webdav") {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
