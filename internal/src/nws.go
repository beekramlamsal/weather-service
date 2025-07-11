package src

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetForecast fetches the weather forecast for a given latitude and longitude.
// It first resolves the forecast URL from the NWS points API,
// then fetches the actual forecast with retry logic for transient errors (e.g. 429).
func GetForecast(ctx context.Context, client *http.Client, tmpl string, lat, lon float64, retries int, backoff time.Duration) (string, error) {
	// Format the initial NWS point URL using lat/lon
	pointURL := fmt.Sprintf(tmpl, lat, lon)

	// Step 1: Fetch the forecast URL from the NWS point API
	pointReq, err := http.NewRequestWithContext(ctx, http.MethodGet, pointURL, nil)
	if err != nil {
		return "", fmt.Errorf("create point request: %w", err)
	}

	resp, err := client.Do(pointReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the forecast URL from the response
	var point struct {
		Properties struct {
			Forecast string `json:"forecast"`
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&point); err != nil {
		return "", err
	}

	forecastURL := point.Properties.Forecast

	// Step 2: Fetch the actual forecast using the URL, with retries
	for i := 0; i <= retries; i++ {
		forecastReq, err := http.NewRequestWithContext(ctx, http.MethodGet, forecastURL, nil)
		if err != nil {
			return "", fmt.Errorf("create forecast request: %w", err)
		}

		resp, err := client.Do(forecastReq)
		if err != nil {
			return "", err
		}

		// Retry if server returns 429 (too many requests), with increasing delay
		if resp.StatusCode == http.StatusTooManyRequests && i < retries {
			resp.Body.Close()
			time.Sleep(backoff * time.Duration(i+1))
			continue
		}

		// Decode the forecast periods array from the JSON response
		var forecast struct {
			Properties struct {
				Periods []struct {
					DetailedForecast string `json:"detailedForecast"`
				}
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
			resp.Body.Close()
			return "", err
		}
		resp.Body.Close()

		// Check for at least one forecast period and return it
		if len(forecast.Properties.Periods) == 0 {
			return "", fmt.Errorf("no forecast periods")
		}

		return forecast.Properties.Periods[0].DetailedForecast, nil
	}

	// If all retries are exhausted, return an error
	return "", fmt.Errorf("forecast retry limit exceeded")
}
