package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetForecast(ctx context.Context, client *http.Client, tmpl string, lat, lon float64, retries int, backoff time.Duration) (string, error) {
	pointURL := fmt.Sprintf(tmpl, lat, lon)
	resp, err := client.Get(pointURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var point struct {
		Properties struct {
			Forecast string `json:"forecast"`
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&point); err != nil {
		return "", err
	}

	forecastURL := point.Properties.Forecast

	for i := 0; i <= retries; i++ {
		resp, err := client.Get(forecastURL)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests && i < retries {
			time.Sleep(backoff * time.Duration(i+1))
			continue
		}

		var forecast struct {
			Properties struct {
				Periods []struct {
					DetailedForecast string `json:"detailedForecast"`
				}
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
			return "", err
		}
		if len(forecast.Properties.Periods) == 0 {
			return "", fmt.Errorf("no forecast periods")
		}
		return forecast.Properties.Periods[0].DetailedForecast, nil
	}

	return "", fmt.Errorf("forecast retry limit exceeded")
}
