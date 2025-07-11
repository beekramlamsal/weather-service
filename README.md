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

# TODOs for Production Readiness

This service covers the critical pieces of functionality, but there are a few areas I would prioritize next to make this truly production-ready:

---

## Error Resilience & Retries

- [ ] Add circuit breaker pattern to avoid hammering upstream APIs when they’re down
- [ ] Use exponential backoff with jitter for retries (currently fixed delay)
- [ ] Support fallback (e.g. cached or static responses) when upstream is unavailable

---

## Observability

- [ ] Add structured logging (e.g. zap or zerolog) instead of standard log.Println
- [ ] Add Prometheus metrics for:
  - Request/response time
  - Error count per upstream
  - Retry attempts
- [ ] Add tracing support (OpenTelemetry / Jaeger) to trace calls across services

---

## Security

- [ ] Use HTTPS in production (via reverse proxy or direct TLS support)
- [ ] Sanitize any downstream responses used in output
- [ ] Validate inputs and headers more strictly if user input is added later

---

## Deployment & Scaling

- [ ] Add readiness and liveness probes to deployment spec
- [ ] Use a production-grade logger that writes to stdout with timestamps for container logs
- [ ] Add horizontal pod autoscaler config (HPA) based on CPU or latency

---

## Config and Secrets

- [ ] Support environment variable overrides (e.g. `CONFIG_PATH`, or config via ENV)
- [ ] Add secret handling for future private APIs (if required)

---

## CI/CD & Testing

- [ ] Add GitHub Actions workflow for:
  - Linting (golangci-lint)
  - Running tests on PR
  - Build & push Docker image
- [ ] Add integration tests for `/api` handler using `httptest`
- [ ] Add load test profile using `k6` or `wrk`

---

## Frontend UX (if extended)

- [ ] Improve weather icon logic to prioritize first sentence and confidence words
- [ ] Handle empty or malformed forecasts gracefully on frontend

---

### Notes

I've intentionally written this codebase with modularity and testability in mind, so extending any of these should be relatively straightforward. Prioritization would depend on expected load, failure tolerance, and security needs in production.

