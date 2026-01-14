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
	// Get all local IPs
	localIPs, err := GetAllLocalIPs()
	if err != nil {
		fmt.Printf("Error getting local IPs: %v\n", err)
	}
	// Add common localhost variants
	localIPs = append(localIPs, "localhost", "127.0.0.1")
	localIPs = append(localIPs, "aaxion-cdn.codershubinc.tech")

	fmt.Printf("Allowing CORS for IPs: %v\n", localIPs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow requests from any local IP
		if origin != "" {
			for _, ip := range localIPs {
				if strings.Contains(origin, ip) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
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

func GetAllLocalIPs() ([]string, error) {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && !ip.IsLoopback() {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips, nil
}
