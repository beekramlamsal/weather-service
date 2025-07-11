package main

import (
	"log"
	"net/http"

	"github.com/beekramlamsal/weather-service/internal/config"
	"github.com/beekramlamsal/weather-service/internal/handler"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatal("Config error:", err)
	}

	log.Println("Starting server on", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, handler.New(cfg)); err != nil {
		log.Fatal(err)
	}
}
