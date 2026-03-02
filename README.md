# go-service

A minimal, production-ready HTTP service written in Go, demonstrating clean code structure, configuration management, reliability, and observability.

---

## Overview

`go-service` exposes four HTTP endpoints:

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness check – always returns `{"status":"ok"}` |
| GET | `/config` | Returns the current runtime configuration |
| POST | `/echo` | Echoes a `message` field back with an RFC3339 timestamp |
| GET | `/metrics` | In-memory request counters (requires `ENABLE_METRICS=true`) |

---

## Prerequisites

| Tool | Minimum version |
|------|----------------|
| [Go](https://go.dev/dl/) | 1.21 |
| [Docker](https://docs.docker.com/get-docker/) | 24 (for containerised run) |

---

## Environment Variables

| Variable | Required | Default | Validation |
|----------|----------|---------|-----------|
| `APP_HOST` | ✅ | – | Non-empty string |
| `APP_PORT` | ✅ | – | Integer 1–65535 |
| `LOG_LEVEL` | ✅ | – | `debug` \| `info` \| `warn` \| `error` |
| `ALLOWED_ORIGINS` | ❌ | `""` | Comma-separated URLs |
| `ENABLE_METRICS` | ❌ | `false` | Boolean (`true`/`false`) |

---

## Running Locally (without Docker)

```bash
# 1. Install dependencies (none external – stdlib only)
go mod download

# 2. Export required environment variables
export APP_HOST=0.0.0.0
export APP_PORT=8080
export LOG_LEVEL=info
export ENABLE_METRICS=true          # optional
export ALLOWED_ORIGINS="https://a.com,https://b.com"  # optional

# 3. Run
go run .
```

On Windows (PowerShell):

```powershell
$env:APP_HOST="0.0.0.0"
$env:APP_PORT="8080"
$env:LOG_LEVEL="info"
$env:ENABLE_METRICS="true"
go run .
```

---

## Running with Docker

```bash
# Build the image
docker build -t go-service .

# Run (minimal)
docker run --rm -p 8080:8080 \
  -e APP_HOST=0.0.0.0 \
  -e APP_PORT=8080 \
  -e LOG_LEVEL=info \
  go-service

# Run (with all options)
docker run --rm -p 8080:8080 \
  -e APP_HOST=0.0.0.0 \
  -e APP_PORT=8080 \
  -e LOG_LEVEL=debug \
  -e ENABLE_METRICS=true \
  -e ALLOWED_ORIGINS="https://a.com,https://b.com" \
  go-service
```

---

## Sample curl Commands

```bash
# Health check
curl http://localhost:8080/health

# View configuration
curl http://localhost:8080/config

# Echo a message
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, World!"}'

# Bad echo – empty message (returns 400)
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": ""}'

# View metrics (requires ENABLE_METRICS=true)
curl http://localhost:8080/metrics
```

---

## Testing

```bash
# Run all tests with verbose output
go test -v ./...

# Run with race detector
go test -race ./...

# Run only config tests
go test -v -run TestLoadConfig ./...

# Run only handler tests
go test -v -run Test.*Handler ./...
```

### Test coverage summary

| File | Tests |
|------|-------|
| `config_test.go` | 12 – covers all required fields, validation edge cases, optional fields |
| `handlers_test.go` | 12 – covers all 4 handlers, validation errors, counter behaviour |

---

## Project Structure

```
go-service/
├── main.go           # Entry point; HTTP server with graceful shutdown
├── config.go         # Config struct + LoadConfig() with validation
├── handlers.go       # HTTP handler functions (App receiver pattern)
├── metrics.go        # Atomic in-memory request counters
├── config_test.go    # Configuration loading tests
├── handlers_test.go  # HTTP handler tests (httptest)
├── go.mod            # Module definition
├── Dockerfile        # Multi-stage build (golang:alpine → alpine)
├── .dockerignore     # Files excluded from Docker build context
└── README.md         # This file
```

---

## Design & Assumptions

1. **No external dependencies** – uses only the Go standard library (`net/http`, `encoding/json`, `os`, `sync/atomic`, etc.).
2. **ALLOWED_ORIGINS format** – comma-separated string; leading/trailing whitespace around each origin is trimmed.
3. **Timestamp format** – RFC3339 / ISO 8601 UTC (e.g. `2024-03-02T10:30:45Z`).
4. **Metrics persistence** – in-memory only; counters reset on restart or container recreation.
5. **Error responses** – all errors (400, 404, 500) are returned as `{"error": "..."}` JSON.
6. **Graceful shutdown** – the server drains active connections on SIGINT / SIGTERM with a 10-second deadline.
7. **CORS** – `ALLOWED_ORIGINS` is stored and exposed via `/config` but CORS headers are not automatically injected (adding a middleware is straightforward if needed).
8. **Log level** – the `LOG_LEVEL` value is validated and stored; structured log-level filtering can be layered on top using the standard `slog` package (Go 1.21+).
9. **Non-root Docker user** – the container runs as `appuser` to follow the principle of least privilege.

---

## Troubleshooting

| Symptom | Likely cause | Fix |
|---------|-------------|-----|
| `configuration error: APP_HOST is required` | Missing env var | Export `APP_HOST` |
| `configuration error: APP_PORT must be between 1 and 65535` | Port out of range | Use a valid port |
| `404` on `/metrics` | `ENABLE_METRICS` not set | Set `ENABLE_METRICS=true` |
| `400` on `/echo` | Invalid JSON or empty message | Provide `{"message":"..."}` |
| Port already in use | Another process on the port | Change `APP_PORT` or stop the other process |
