package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr            string        `yaml:"addr"`
	LocationsURL    string        `yaml:"locations_url"`
	NWSPointURLTmpl string        `yaml:"nws_point_url_tmpl"`
	Timeout         time.Duration `yaml:"timeout"`
	RetryCount      int           `yaml:"retry_count"`
	RetryBackoff    time.Duration `yaml:"retry_backoff"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
