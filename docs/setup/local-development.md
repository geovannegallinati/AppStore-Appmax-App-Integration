# Local Development

Docker setup, Makefile targets, hot reload, ngrok tunneling, and helper scripts.

---

## Prerequisites

- **Docker** and **Docker Compose** (v2)
- **ngrok account** — [sign up](https://dashboard.ngrok.com/signup) (free tier works)
- (Optional) **Go 1.25+** for running tests outside Docker

---

## Quick Start

```bash
make install
```

This runs the full setup pipeline:

1. `env-init` — copies `.env.example` to `.env` if missing
2. `generate-key` — generates a 32-character `APP_KEY` if not set
3. `teardown` — removes old containers and volumes
4. `up` — builds and starts all containers
5. `migrate` — runs database migrations
6. `test` — runs Go tests inside the app container
7. `validate` — verifies all endpoints are reachable via ngrok

---

## Docker Compose Services

| Service | Image | Purpose | Ports | Health Check |
|---------|-------|---------|-------|-------------|
| `postgres` | `postgres:16-alpine` | Database | `5432` | `pg_isready` every 5s, 10 retries |
| `redis` | `redis:7-alpine` | Cache (install state, merchant tokens) | `6379` | `redis-cli ping` every 5s, 10 retries |
| `app` | Dockerfile `dev` target | Go app with Air hot reload | `8080` | `wget /health` every 5s, 30 retries, 180s start period |
| `ngrok` | `ngrok/ngrok:alpine` | Public HTTPS tunnel | `4040` (inspector) | Checks `/api/tunnels` every 3s, 20 retries |
| `nginx` | `nginx:alpine` | Reverse proxy | `80` | — |

Startup order: `postgres` + `redis` → `app` (waits for both healthy) → `ngrok` (waits for app healthy) → `nginx`.

---

## Dockerfile Targets

The Dockerfile has two build targets:

### `builder` (production)
```dockerfile
FROM golang:1.25-alpine AS builder
# CGO_ENABLED=0 static binary → /app/server
```

### `dev` (development — used by docker-compose)
```dockerfile
FROM golang:1.25-alpine AS dev
# Installs Air for hot reload
# Source code mounted via volume
# CMD ["air"]
```

---

## Air Hot Reload

Config: `.air.toml`

| Setting | Value |
|---------|-------|
| Watches | `*.go` files, `.env` |
| Excludes | `tmp/`, `.git/`, `tests/`, `vendor/` |
| Build command | `go build -o ./tmp/server ./main.go` |
| Delay after change | 500ms |

When you edit a `.go` file, Air rebuilds and restarts the server automatically.

---

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make install` | Full setup: env → key → teardown → up → migrate → test → validate |
| `make up` | Build images and start all containers |
| `make down` | Stop containers, remove volumes and orphans |
| `make restart` | Restart all containers |
| `make logs` | Follow all container logs (`docker compose logs -f`) |
| `make health` | Wait for `/health` to respond (up to 60 attempts, 3s apart) |
| `make validate` | Verify all endpoints reachable via ngrok (checks frontend, health, install, callback, webhook) |
| `make test` | Run `go test ./...` inside the app container |
| `make migrate` | Run Goravel migrations (`artisan migrate`) |
| `make env-init` | Copy `.env.example` → `.env` if missing |
| `make generate-key` | Generate 32-char random `APP_KEY` if not set or invalid length |
| `make teardown` | Remove all containers and volumes |
| `make rename-module NEW=github.com/foo/bar` | Rename Go module path across all files |

---

## ngrok Setup

ngrok provides a public HTTPS URL that tunnels to your local app. Appmax needs this
to reach your callback and webhook endpoints.

### Configuration

In `.env`:
```
NGROK_AUTHTOKEN=your_token_here
NGROK_URL=                          # leave empty for dynamic URL
```

| Setting | Behavior |
|---------|----------|
| `NGROK_AUTHTOKEN` set, `NGROK_URL` empty | ngrok assigns a random URL (changes on restart) |
| `NGROK_AUTHTOKEN` set, `NGROK_URL` set | ngrok uses your reserved static domain |
| `NGROK_AUTHTOKEN` empty | ngrok container fails with instructions |

### Inspector

ngrok exposes a local web inspector at `http://localhost:4040` where you can see
all tunneled requests, replay them, and inspect headers/bodies.

### How `make validate` Works

1. Waits for the app health check to pass
2. Requires `NGROK_AUTHTOKEN`; warns for each missing Appmax credential (`APPMAX_CLIENT_ID`, `APPMAX_CLIENT_SECRET`, `APPMAX_APP_ID_UUID`, `APPMAX_APP_ID_NUMERIC`)
3. Queries ngrok's local API (`http://127.0.0.1:4040/api/tunnels`) to discover the active tunnel URL
4. Verifies all 5 public endpoints respond through the tunnel:
   - `GET /` (frontend)
   - `GET /health`
   - `GET /install/start`
   - `GET /integrations/appmax/callback/install`
   - `GET /webhooks/appmax`

---

## Scripts

| Script                    | Platform             | Purpose                                             |
|---------------------------|----------------------|-----------------------------------------------------|
| `install.ps1`             | Windows (PowerShell) | Full setup equivalent of `make install`             |
| `scripts/rename-module.sh`| Unix                 | Renames the Go module path across all project files |

---

## Environment Variables

See `.env.example` for all available variables. Key groups:

| Group | Variables | Notes |
|-------|-----------|-------|
| App | `APP_NAME`, `APP_ENV`, `APP_DEBUG`, `APP_HOST`, `APP_PORT`, `APP_URL`, `APP_KEY` | `APP_KEY` auto-generated by `make install` |
| ngrok | `NGROK_AUTHTOKEN`, `NGROK_URL` | Required for public tunnel |
| Database | `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD` | `DB_HOST=postgres` inside Docker |
| Redis | `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD` | `REDIS_HOST=redis` inside Docker |
| Appmax URLs | `APPMAX_AUTH_URL`, `APPMAX_API_URL`, `APPMAX_ADMIN_URL` | Defaults to production; set sandbox URLs for testing |
| Appmax Credentials | `APPMAX_CLIENT_ID`, `APPMAX_CLIENT_SECRET`, `APPMAX_APP_ID_UUID`, `APPMAX_APP_ID_NUMERIC` | From Appmax AppStore |

---

## Troubleshooting

### App container health check fails

```bash
make logs
```

Common causes:
- Missing required env vars (`APPMAX_CLIENT_ID`, `APPMAX_CLIENT_SECRET`, etc.)
- Database migration error (check postgres logs)
- Port conflict on 8080

### ngrok tunnel not found

- Verify `NGROK_AUTHTOKEN` is set in `.env`
- Check ngrok inspector: `http://localhost:4040`
- If using a reserved domain (`NGROK_URL`), verify it matches your ngrok account
- If running the **Appmax Endpoints flow** in Postman, ensure the `NGROK_URL` collection variable is set to your tunnel URL — see [Postman Variables](../postman/postman-variables.md)

### Database connection refused

- Wait for postgres health check to pass (can take 30-50s on first start)
- Inside Docker, the app uses `DB_HOST=postgres` (set automatically by docker-compose)
- Locally, use `DB_HOST=127.0.0.1`

### "Air" not rebuilding

- Check `.air.toml` excludes — files in `tmp/`, `tests/`, `vendor/` are ignored
- Ensure the file has a `.go` extension
- Check `make logs` for build errors
