package discovery

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/grandcat/zeroconf"
)

func StartDiscoveryService(port int) {
	hostname, _ := os.Hostname()

	d, err := GetDiscoveryDevices()
	if err != nil || d.ID == "" {
		log.Printf("Failed to get discovery device info: %v", err)
		cd, err := CreateDiscoveryDevice(hostname)
		if err != nil {
			log.Printf("Failed to create discovery device: %v", err)
		}
		log.Printf("Created discovery device with ID: %s", cd)
		d, err = GetDiscoveryDevices()
		if err != nil {
			log.Printf("Failed to get discovery device info after creation: %v", err)
		}
	}
	log.Println("got device", d)
	instanceName := fmt.Sprintf("Aaxion Server + %s", hostname)

	// Metadata to help clients identify the server
	meta := []string{
		"version=unreleased",
		"description=Aaxion File Server",
		fmt.Sprintf("device_id=%s", d.ID),
		fmt.Sprintf("device_name=%s", d.Name),
	}

	// Register the service
	// Instance name: instanceName
	// Service type: "_aaxion._tcp" (clients should search for this)
	// Domain: "local."
	server, err := zeroconf.Register(
		instanceName,
		"_aaxion._tcp",
		"local.",
		port,
		meta,
		nil,
	)
	if err != nil {
		log.Printf("Failed to register discovery service: %v", err)
		return
	}

	log.Printf("Discovery service started on port %d (_aaxion._tcp.local.)", port)

	// Build a shutdown mechanism to deregister the service cleanly
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig

		log.Println("Stopping discovery service...")
		server.Shutdown()
		log.Println("Cleanup complete. Exiting.")
		os.Exit(0)
	}()
}
