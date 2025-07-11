package main

import (
	"log"
	"net/http"

	"github.com/beekramlamsal/weather-service/internal/config"
	"github.com/beekramlamsal/weather-service/internal/handler"
)

func main() {
	// Load configuration from file
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatal("Config error:", err)
	}

	log.Printf("Using client timeout: %s\n", cfg.Timeout)
	log.Printf("Starting server on %s\n", cfg.Addr)

	h := handler.New(cfg)

	http.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	http.Handle("/", h)

	if err := http.ListenAndServe(cfg.Addr, nil); err != nil {
		log.Fatal(err)
	}
}
