package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IceWhaleTech/CasaOS-Gateway/route"
	"github.com/IceWhaleTech/CasaOS-Gateway/service"
	"github.com/caddyserver/caddy/v2"
)

func main() {
	// Initialize your management service
	state := service.NewState()
	management := service.NewManagementService(state)

	// Create new Caddy gateway
	gateway := route.NewCaddyGateway(management)

	// Configure Caddy
	config, err := gateway.ConfigureCaddy()
	if err != nil {
		log.Fatalf("Error configuring Caddy: %v", err)
	}

	// Start Caddy
	err = caddy.Run(config)
	if err != nil {
		log.Fatalf("Error starting Caddy: %v", err)
	}

	// Wait for shutdown signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// Graceful shutdown
	err = caddy.Stop()
	if err != nil {
		log.Printf("Error stopping Caddy: %v", err)
	}
}
