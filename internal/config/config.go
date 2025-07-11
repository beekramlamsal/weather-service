package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config holds all the service configuration values loaded from config.yaml
type Config struct {
	Addr            string        `yaml:"addr"`               // Address and port the server will listen on (e.g. ":5000")
	LocationsURL    string        `yaml:"locations_url"`      // URL for fetching random location
	NWSPointURLTmpl string        `yaml:"nws_point_url_tmpl"` // Template URL to get NWS forecast endpoint for given lat/lon
	Timeout         time.Duration `yaml:"timeout"`            // HTTP client timeout duration
	RetryCount      int           `yaml:"retry_count"`        // Number of times to retry on certain HTTP errors like 429
	RetryBackoff    time.Duration `yaml:"retry_backoff"`      // Time delay between retries, increases with each attempt
}

// Load reads and parses the YAML config file from the given path
func Load(path string) (*Config, error) {
	// Open the config file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Decode YAML content into Config struct
	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	// Return the loaded config struct
	return &cfg, nil
}
