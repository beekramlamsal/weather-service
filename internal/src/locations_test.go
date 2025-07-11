package src

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRandomLocation_HTTPErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expectErr  bool
	}{
		{"BadRequest", http.StatusBadRequest, true},
		{"NotFound", http.StatusNotFound, true},
		{"InternalServerError", http.StatusInternalServerError, true},
		{"ServiceUnavailable", http.StatusServiceUnavailable, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(`{"error": "test error"}`))
			}))
			defer srv.Close()

			_, err := GetRandomLocation(context.Background(), srv.Client(), srv.URL)
			if tt.expectErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetRandomLocation_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer srv.Close()

	_, err := GetRandomLocation(context.Background(), srv.Client(), srv.URL)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestGetRandomLocation_NetworkTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"locations": [{"name": "Test", "latitude": 1, "longitude": 2}]}`))
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 50 * time.Millisecond}
	_, err := GetRandomLocation(context.Background(), client, srv.URL)
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestGetRandomLocation_NetworkError(t *testing.T) {
	// Use invalid URL to trigger network error
	_, err := GetRandomLocation(context.Background(), &http.Client{}, "http://invalid-url-12345")
	if err == nil {
		t.Error("expected network error")
	}
}
