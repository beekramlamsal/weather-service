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

## How to Run It

1. Clone the repo:

```bash
git clone https://github.com/yourusername/weather-service.git
cd weather-service
```

2. Run the app (make sure you’re in the project root):

```bash
go run ./cmd/server
```

3. Hit the endpoint:

```bash
curl http://localhost:5050/api
```

You’ll get something like:

```
The weather in Litchfield is: A chance of showers and thunderstorms after 4pm.
```

You can also check the health endpoint:

```
curl http://localhost:5000/healthz
# → ok
```

---

## Testing

I’ve written unit tests for both upstream clients and config loading. You can run them with:

```bash
go test ./...
```

They include:
- Success cases
- Handling empty/missing data
- Retry logic for 429 responses

---


