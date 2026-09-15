package main

import (
	"aaxion/internal/api"
	"aaxion/internal/cli"
	"aaxion/internal/db"
	"aaxion/internal/discovery"
	"fmt"
	"log"
	"net/http"
)

func init() {
	cli.RegisterCommand(cli.Command{
		Name:        "serve",
		Description: "Start the Aaxion web server (default)",
		Run: func(args []string) {
			err := db.InitDb()
			if err != nil {
				log.Fatal("Got err initializing DB:", err)
			}
			startServer()
		},
	})
}

func main() {
	cli.Execute()
}

func startServer() {
	fmt.Println(`
    _        _    __  __  ___   ___   _   _ 
   / \      / \   \ \/ / |_ _| / _ \ | \ | |
  / _ \    / _ \   >  <   | | | | | ||  \| |
 / ___ \  / ___ \ / /\ \  | | | |_| || |\  |
/_/   \_\/_/   \_\/_/\_\ |___| \___/ |_| \_|
`)

	port := 8080
	fmt.Println("Starting server...")
	api.RegisterRoutes()

	discovery.StartDiscoveryService(port)
	log.Println("mDNS discovery service started at port", port)

	// Wrap the default ServeMux with the CORS middleware
	handler := api.CORSMiddleware(http.DefaultServeMux)

	log.Printf("Listening on :%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
