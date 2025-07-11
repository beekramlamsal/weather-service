# Weather Service

A production-ready Go web service that:
- Fetches a random U.S. city from a location provider
- Queries the National Weather Service for a forecast
- Returns a friendly weather summary string

---

## Features

- Configurable via `config.yaml`
- Timeout, retries, and backoff for HTTP requests
- Handles HTTP failures (like 429) gracefully
- Logs API issues for debugging
- Easily deployable via Docker or Compose

---

## Run Locally

```bash
go run ./cmd/server
