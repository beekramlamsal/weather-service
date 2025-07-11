package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/beekramlamsal/weather-service/internal/config"
	"github.com/beekramlamsal/weather-service/internal/src"
)

// New sets up the HTTP routes and returns an http.Handler.
// It configures the backend API at /api and serves static files (HTML/JS/images) from the /client directory.
func New(cfg *config.Config) http.Handler {
	// Create two separate clients to allow more generous timeout for forecast
	locationClient := &http.Client{Timeout: cfg.Timeout}
	forecastClient := &http.Client{Timeout: cfg.Timeout * 2}

	mux := http.NewServeMux()

	// Weather API endpoint at /api
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// Use context with timeout for controlling downstream requests
		ctx, cancel := context.WithTimeout(r.Context(), cfg.Timeout*3)
		defer cancel()

		// Call location service
		loc, err := src.GetRandomLocation(ctx, locationClient, cfg.LocationsURL)
		if err != nil {
			// Log the error for troubleshooting
			log.Println("Location fetch failed:", err)
			http.Error(w, "location error", http.StatusBadGateway)
			return
		}

		// Call forecast service
		forecast, err := src.GetForecast(ctx, forecastClient, cfg.NWSPointURLTmpl, loc.Latitude, loc.Longitude, cfg.RetryCount, cfg.RetryBackoff)
		if err != nil {
			// Respond with an error if forecast fetch fails
			log.Println("Forecast fetch failed:", err)
			http.Error(w, "forecast error", http.StatusBadGateway)
			return
		}

		// Return formatted weather response
		fmt.Fprintf(w, "Weather in %s: %s", loc.Name, forecast)
	})

	// Serve static client files at root URL (index.html, script.js, and images)
	// Any request not matching /api will load from client directory
	fs := http.FileServer(http.Dir("client"))
	mux.Handle("/", fs)

	return mux
}
