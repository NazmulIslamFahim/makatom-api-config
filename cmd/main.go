package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"common/pkg/config"
	"common/pkg/database/mongodb"
	"makatom-api-config/internal/routes"
)

func main() {
	// Initialize configuration
	config.Init()
	cfg := config.GetConfig()

	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mongodb.Manager.Connect(ctx, "config_db", cfg.MongoURI, 30*time.Minute)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongodb.Manager.DisconnectAll(context.Background()); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Register routes and get the mux
	mux := routes.RegisterConfigRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting server on %s", cfg.Port)
	log.Printf("Environment: %s", cfg.Environment)
	log.Printf("Debug mode: %t", cfg.Debug)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
