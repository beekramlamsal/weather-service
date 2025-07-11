package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/beekramlamsal/weather-service/internal/client"
	"github.com/beekramlamsal/weather-service/internal/config"
)

func New(cfg *config.Config) http.Handler {
	httpClient := &http.Client{Timeout: cfg.Timeout}
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), cfg.Timeout)
		defer cancel()

		loc, err := client.GetRandomLocation(ctx, httpClient, cfg.LocationsURL)
		if err != nil {
			log.Println("🔴 Location fetch failed:", err) // 👈 Add this
			http.Error(w, "location error", http.StatusBadGateway)
			return
		}

		forecast, err := client.GetForecast(ctx, httpClient, cfg.NWSPointURLTmpl, loc.Latitude, loc.Longitude, cfg.RetryCount, cfg.RetryBackoff)
		if err != nil {
			http.Error(w, "forecast error", http.StatusBadGateway)
			return
		}

		fmt.Fprintf(w, "Weather in %s: %s", loc.Name, forecast)
	})

	return mux
}
