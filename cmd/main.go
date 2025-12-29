package main

import (
	"aaxion/internal/api"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("its aaxion")
	startServer()
}

func startServer() {
	fmt.Println("Starting server...")
	api.RegisterRoutes()
	http.ListenAndServe(":8080", nil)
}
