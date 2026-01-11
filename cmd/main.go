package main

import (
	"aaxion/internal/api"
	"aaxion/internal/db"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	fmt.Println("its aaxion")
	err := db.InitDb()
	if err != nil {
		log.Println("Got err", err)
	}
	startServer()
}

func startServer() {
	fmt.Println("Starting server...")
	api.RegisterRoutes()

	// Wrap the default mux with CORS middleware
	handler := corsMiddleware(http.DefaultServeMux)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	// Get local IP once
	localIP := GetOutboundIP()
	localIPStr := localIP.String()
	fmt.Printf("Allowing CORS for Local IP: %s\n", localIPStr)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow requests from localhost and the dynamic local IP
		if origin != "" {
			if strings.Contains(origin, "localhost") ||
				strings.Contains(origin, "127.0.0.1") ||
				strings.Contains(origin, localIPStr) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IPv4(127, 0, 0, 1)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
