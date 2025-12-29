package main

import (
	"aaxion/internal/api"
	"aaxion/internal/db"
	"fmt"
	"log"
	"net/http"
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
	http.ListenAndServe(":8080", nil)
}
