package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/beekramlamsal/weather-service/internal/model"
)

func GetRandomLocation(ctx context.Context, client *http.Client, url string) (*model.Location, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET location: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	var data model.LocationsResp
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if len(data.Locations) == 0 {
		return nil, fmt.Errorf("empty location list")
	}
	return &data.Locations[0], nil
}
