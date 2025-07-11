package src

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetForecast_PointAPIErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expectErr  bool
	}{
		{"BadRequest", http.StatusBadRequest, true},
		{"NotFound", http.StatusNotFound, true},
		{"InternalServerError", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pointSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(`{"error": "test error"}`))
			}))
			defer pointSrv.Close()

			_, err := GetForecast(context.Background(), pointSrv.Client(), pointSrv.URL, 1.0, 2.0, 1, 1*time.Millisecond)
			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
		})
	}
}

func TestGetForecast_ForecastAPIErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expectErr  bool
	}{
		{"BadRequest", http.StatusBadRequest, true},
		{"NotFound", http.StatusNotFound, true},
		{"InternalServerError", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First server returns point with forecast URL
			pointSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				forecastSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
					w.Write([]byte(`{"error": "test error"}`))
				}))
				defer forecastSrv.Close()

				response := `{"properties": {"forecast": "` + forecastSrv.URL + `"}}`
				w.Write([]byte(response))
			}))
			defer pointSrv.Close()

			_, err := GetForecast(context.Background(), pointSrv.Client(), pointSrv.URL, 1.0, 2.0, 1, 1*time.Millisecond)
			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
		})
	}
}

func TestGetForecast_RetryLogic429(t *testing.T) {
	attemptCount := 0

	// Forecast server that returns 429 first two times, then success
	forecastSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "too many requests"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"properties": {"periods": [{"detailedForecast": "Sunny"}]}}`))
	}))
	defer forecastSrv.Close()

	// Point server that returns the forecast URL
	pointSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"properties": {"forecast": "` + forecastSrv.URL + `"}}`
		w.Write([]byte(response))
	}))
	defer pointSrv.Close()

	tmpl := pointSrv.URL + "/points/%f,%f"

	forecast, err := GetForecast(context.Background(), pointSrv.Client(), tmpl, 1.0, 2.0, 3, 1*time.Millisecond)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if forecast != "Sunny" {
		t.Errorf("expected Sunny, got %s", forecast)
	}
	if attemptCount != 3 {
		t.Errorf("expected 3 attempts, got %d", attemptCount)
	}
}

func TestGetForecast_InvalidPointJSON(t *testing.T) {
	pointSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer pointSrv.Close()

	_, err := GetForecast(context.Background(), pointSrv.Client(), pointSrv.URL, 1.0, 2.0, 1, 1*time.Millisecond)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestGetForecast_InvalidForecastJSON(t *testing.T) {
	// Forecast server that returns invalid JSON
	forecastSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer forecastSrv.Close()

	// Point server that returns the forecast URL
	pointSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"properties": {"forecast": "` + forecastSrv.URL + `"}}`
		w.Write([]byte(response))
	}))
	defer pointSrv.Close()

	_, err := GetForecast(context.Background(), pointSrv.Client(), pointSrv.URL, 1.0, 2.0, 1, 1*time.Millisecond)
	if err == nil {
		t.Error("expected error for invalid forecast JSON")
	}
}

func TestGetForecast_NoForecastPeriods(t *testing.T) {
	// Forecast server that returns empty periods
	forecastSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"properties": {"periods": []}}`))
	}))
	defer forecastSrv.Close()

	// Point server that returns the forecast URL
	pointSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := `{"properties": {"forecast": "` + forecastSrv.URL + `"}}`
		w.Write([]byte(response))
	}))
	defer pointSrv.Close()

	tmpl := pointSrv.URL + "/points/%f,%f"

	_, err := GetForecast(context.Background(), pointSrv.Client(), tmpl, 1.0, 2.0, 1, 1*time.Millisecond)
	if err == nil {
		t.Error("expected error for no forecast periods")
	}
	if !strings.Contains(err.Error(), "no forecast periods") {
		t.Errorf("expected no forecast periods error, got: %v", err)
	}
}

func TestGetForecast_EmptyForecastURL(t *testing.T) {
	pointSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"properties": {"forecast": ""}}`))
	}))
	defer pointSrv.Close()

	_, err := GetForecast(context.Background(), pointSrv.Client(), pointSrv.URL, 1.0, 2.0, 1, 1*time.Millisecond)
	if err == nil {
		t.Error("expected error for empty forecast URL")
	}
}
