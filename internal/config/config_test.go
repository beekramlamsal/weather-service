package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	f, err := os.CreateTemp("", "config*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.Write([]byte(`
addr: ":5050"
locations_url: "http://example.com"
nws_point_url_tmpl: "http://api.weather.gov/points/%f,%f"
timeout: 5s
retry_count: 2
retry_backoff: 1s
`))
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Addr != ":5050" {
		t.Errorf("unexpected addr: %s", cfg.Addr)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("unexpected timeout: %s", cfg.Timeout)
	}
}
