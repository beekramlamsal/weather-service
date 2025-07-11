package main

import (
	"log"
	"net/http"

	"github.com/beekramlamsal/weather-service/internal/config"
	"github.com/beekramlamsal/weather-service/internal/handler"
)

func main() {
	// Load configuration from config.yaml file
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatal("Config error:", err)
	}

	// Log the configured timeout for verification
	log.Printf("Using client timeout: %s\n", cfg.Timeout)

	// Start the HTTP server on the configured address and port
	log.Println("Starting server on", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, handler.New(cfg)); err != nil {
		log.Fatal(err)
	}
}
