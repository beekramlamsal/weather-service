package src

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beekramlamsal/weather-service/internal/model"
)

// GetRandomLocation fetches a random location from the external location service API.
// It expects a single location object inside a "locations" array in the JSON response.
func GetRandomLocation(ctx context.Context, client *http.Client, url string) (*model.Location, error) {
	// Build a new GET request with context for timeout control
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Execute the HTTP request using the provided client
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET location: %w", err)
	}
	defer resp.Body.Close()

	// Check if the status code is not 200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	// Decode the JSON response into LocationsResp struct
	var data model.LocationsResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// If no locations returned, report it
	if len(data.Locations) == 0 {
		return nil, fmt.Errorf("empty location list")
	}

	// Return the first location in the list
	return &data.Locations[0], nil
}
