# AppStore Appmax App Integration

Appmax AppStore integration built with Go, Goravel, PostgreSQL, and Redis. Handles
merchant installation (OAuth + health check), checkout (credit card, Pix, boleto),
and webhook processing for order status updates.

---

## Tech Stack

| Component | Version |
|-----------|---------|
| Go | 1.25 |
| Goravel | 1.17 |
| PostgreSQL | 16 (Alpine) |
| Redis | 7 (Alpine) |
| HTTP Router | Gin |
| ORM | GORM |
| Infrastructure | Docker Compose, ngrok, nginx, Air (hot reload) |

---

## API Documentation

| Resource | Link |
|----------|------|
| Postman Documentation and Collection | <a href="https://www.postman.com/geovannegallinati/appmax-full-integration-suite/documentation/52585908-1b3d2ef6-e083-456c-b4c1-fdc46a37f771?sideView=agentMode" target="_blank"><img src="https://run.pstmn.io/button.svg" alt="Run in Postman"></a> |
| Collection JSON | [Appmax — Full Integration Suite.postman_collection.json](docs/postman/Appmax%20-%20Full%20Integration%20Suite.postman_collection.json) |

The collection covers two top-level folders:

- **Appmax Endpoints** — calls the Appmax sandbox APIs directly (OAuth token, installation authorization, customers, orders, payments, webhooks)
- **Localhost Endpoints** — calls the backend server at `localhost:8080` (install flow, merchant token sync, checkout, webhook simulation)

Before running requests, set the required variables in the collection. See [Postman Variables](docs/postman/postman-variables.md) for the full variable reference.

---

## Prerequisites

Before you begin, you need three things: Docker, a ngrok account, and Appmax credentials. Follow each step below.

### Step 1: Install Docker

<details>
<summary><strong>macOS</strong></summary>

You have two options:

**Option A: OrbStack (recommended)** — lighter, faster, uses fewer resources

1. Go to https://orbstack.dev/
2. Download the `.dmg` for your chip (Apple Silicon or Intel)
3. Open the `.dmg` and drag OrbStack to Applications
4. Launch OrbStack from Applications
5. OrbStack is a drop-in Docker Desktop replacement. All `docker` and `docker compose` commands work identically

**Option B: Docker Desktop (official)**

1. Go to https://www.docker.com/products/docker-desktop/
2. Download the `.dmg` for your chip (Apple Silicon or Intel)
3. Open the `.dmg` and drag Docker to Applications
4. Launch Docker Desktop and wait for the whale icon to appear in the menu bar

**Verify installation** (both options):

```bash
docker --version
# Docker version 27.x.x or similar

docker compose version
# Docker Compose version v2.x.x or similar
```

</details>

<details>
<summary><strong>Windows</strong></summary>

1. **Install WSL 2** (if not already installed):
   - Open PowerShell **as Administrator**
   - Run: `wsl --install`
   - Restart your machine when prompted

2. **Install Docker Desktop**:
   - Go to https://www.docker.com/products/docker-desktop/
   - Download the Windows installer
   - Run the installer. When asked, ensure **"Use WSL 2 instead of Hyper-V"** is checked
   - Restart your machine when prompted
   - Launch Docker Desktop from the Start menu

3. **Verify installation** (open a new PowerShell window):

```powershell
docker --version
# Docker version 27.x.x or similar

docker compose version
# Docker Compose version v2.x.x or similar
```

If `docker` is not found, close and reopen PowerShell after Docker Desktop is fully running.

</details>

<details>
<summary><strong>Linux (Ubuntu / Debian)</strong></summary>

```bash
# Update package index
sudo apt update

# Install Docker and Compose plugin
sudo apt install -y docker.io docker-compose-plugin

# Add your user to the docker group (so you don't need sudo)
sudo usermod -aG docker $USER

# IMPORTANT: Log out and log back in for the group change to take effect
# Or run: newgrp docker

# Verify installation
docker --version
docker compose version
```

</details>

---

### Step 2: Create a ngrok Account and Get Your Auth Token

ngrok creates a public URL that tunnels traffic to your local machine. Appmax needs this to reach your app during the install flow and to send webhooks.

1. Go to https://dashboard.ngrok.com/signup
2. Create a free account (sign up with email or GitHub)
3. After login, you land on the dashboard. Go to **Your Authtoken**: https://dashboard.ngrok.com/get-started/your-authtoken
4. Click **Copy** to copy your auth token. It looks like a long string (e.g., `2abc123def456...`)
5. Save this token somewhere — you will paste it into `.env` in a later step

**Optional: Claim a free static domain**

By default, ngrok assigns a random URL every time it starts (e.g., `https://a1b2c3d4.ngrok-free.app`). This means you need to re-register your callback URL in Appmax every time you restart.

To avoid this, claim a free static domain:

1. Go to https://dashboard.ngrok.com/domains
2. Click **Create Domain** (free accounts get 1 static domain)
3. You get a permanent URL like `your-name.ngrok-free.app`
4. Save this domain — you will paste it into `NGROK_URL` in `.env`

---

### Step 3: Create the Appmax App Store Application

As of **March 20, 2026**, the Appmax App Store flow requires you to define your public URLs during the first application setup, before Appmax emails your app credentials.

Because App Store changes are **not automatically replicated into the Appmax sandbox**, this repository recommends the following order:

1. Clone this repository and bootstrap it locally with only `NGROK_AUTHTOKEN` and `NGROK_URL`.
2. Create the App Store application shell in Appmax.
3. In `Meus aplicativos` > select the app > click `Desenvolver` on Appstore, copy the generated public URLs from your project terminal into Appmax on the first setup.
4. Wait for Appmax to email the 4 app credentials.
5. Fill `APPMAX_CLIENT_ID`, `APPMAX_CLIENT_SECRET`, `APPMAX_APP_ID_UUID`, and `APPMAX_APP_ID_NUMERIC`.
6. Reinstall the project to activate the full Appmax installation flow.

If you change those URLs later, contact `desenvolvimento@appmax.com.br` and ask Appmax to replicate the update. Changing them only in your local environment is not enough.

Quick process summary:

1. Open <https://appstore.appmax.com.br/>, click `Desenvolva seu aplicativo`, and register the account with the same CNPJ/company data from Receita Federal.
2. Create the app (`Criar aplicativo`), choose `Aplicativo público` (listed for Appmax users) or `Aplicativo privado` (not publicly listed; only owner/invitees).
3. Fill app data carefully:
   - support email should match the company/CNPJ registration email
   - define name, description, billing model (`Cobrança via Appmax` or `Cobrança via Plataforma Externa`), and fee/commission terms according to Appmax rules
4. Upload app icon with Appmax UI requirements (`1200x1200`, PNG/JPG, square, no rounded corners).
5. Configure permissions for required events and save.
6. After the first local bootstrap in this repository, go to `Meus aplicativos` > select the app > click `Desenvolver` and configure:
   - `Host`: the printed **Install URL** (`https://<your-public-domain>/install/start`)
   - `URL do sistema`: the printed **Frontend URL** (`https://<your-public-domain>/`)
   - `URL de validação` (callback): the printed **Callback URL** (`https://<your-public-domain>/integrations/appmax/callback/install`)
7. If `Desenvolver` does not appear, request enablement at `desenvolvedores@appmax.com.br` (the app may still be under analysis).
8. Callback URL is mandatory. Without it, installation does not complete.
9. Appmax still needs to validate company/CNPJ/shareholder structure and then sends sandbox credentials by email.

References:

- Detailed guide: [How to Get Appmax Credentials and Configure Your App](docs/setup/how-to-get-appmax-credentials-and-configure-your-app.md)
- Official docs: <https://appmax.readme.io/> (for this project, prefer the docs in this repository for implementation details).

---

## Setup

### Step 1: Clone the Repository

```bash
git clone https://github.com/geovannegallinati/AppStore-Appmax-App-Integration.git
cd AppStore-Appmax-App-Integration
```

### Step 2: Create Your `.env` File

<details>
<summary><strong>macOS / Linux</strong></summary>

```bash
cp .env.example .env
```

</details>

<details>
<summary><strong>Windows (PowerShell)</strong></summary>

```powershell
Copy-Item .env.example .env
```

</details>

<details>
<summary><strong>Windows (Command Prompt)</strong></summary>

```cmd
copy .env.example .env
```

</details>

### Step 3: Fill in the `.env` File for the First Bootstrap

Open `.env` in any text editor (VS Code, nano, vim, Notepad++, etc.).

The file is pre-configured with sensible defaults. For the **first bootstrap**, only ngrok needs to be configured.

**Fill in now:**

| Variable | What to Put | Where to Get It |
|----------|------------|-----------------|
| `NGROK_AUTHTOKEN` | Your ngrok auth token | [ngrok dashboard > Your Authtoken](https://dashboard.ngrok.com/get-started/your-authtoken) |

**Recommended to fill in now:**

| Variable | Default | When to Change |
|----------|---------|---------------|
| `NGROK_URL` | *(empty)* | Set this to your static ngrok domain (e.g., `your-name.ngrok-free.app`) if you claimed one. This is strongly recommended for the first Appmax registration, because later URL changes require Appmax-side replication |
| `DB_PASSWORD` | `secret` | Change if you want a different local database password |
| `REDIS_PASSWORD` | `redis-secret` | Change if you want a different local Redis password |

**Leave blank on purpose during the first bootstrap:**

```env
APPMAX_CLIENT_ID=
APPMAX_CLIENT_SECRET=
APPMAX_APP_ID_UUID=
APPMAX_APP_ID_NUMERIC=
```

You will fill those values only after Appmax emails the app credentials.

**Variables you do NOT need to change:**

These are pre-configured for local development with the Appmax sandbox environment:

| Variable | Default | Purpose |
|----------|---------|---------|
| `APP_NAME` | `AppStore Appmax App Integration` | Application name |
| `APP_ENV` | `local` | Environment (local/production) |
| `APP_DEBUG` | `true` | Enable debug mode |
| `APP_HOST` | `0.0.0.0` | Listen address inside Docker |
| `APP_PORT` | `8080` | HTTP port |
| `APP_URL` | `http://127.0.0.1` | Base URL (ngrok overrides this for public access) |
| `APP_KEY` | *(empty)* | Auto-generated on first `install` |
| `DB_HOST` | `127.0.0.1` | Database host (Docker overrides to container name) |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_DATABASE` | `appmax_checkout` | Database name |
| `DB_USERNAME` | `appmax` | Database user |
| `REDIS_HOST` | `127.0.0.1` | Redis host (Docker overrides to container name) |
| `REDIS_PORT` | `6379` | Redis port |
| `APPMAX_AUTH_URL` | `https://auth.sandboxappmax.com.br` | Appmax OAuth server (sandbox) |
| `APPMAX_API_URL` | `https://api.sandboxappmax.com.br` | Appmax API server (sandbox) |
| `APPMAX_ADMIN_URL` | `https://breakingcode.sandboxappmax.com.br` | Appmax admin panel (sandbox) |

> **Switching to production**: Change the three `APPMAX_*_URL` variables to their production equivalents: `https://auth.appmax.com.br`, `https://api.appmax.com.br`, `https://admin.appmax.com.br`.

### Step 4: Run the Project Once to Generate the Public URLs

<details>
<summary><strong>macOS / Linux</strong></summary>

```bash
make install
```

This single command does everything:

1. **env-init** — Copies `.env.example` to `.env` if it doesn't exist
2. **generate-key** — Generates a random 32-character `APP_KEY` if not already set
3. **teardown** — Removes any existing containers and volumes (clean slate)
4. **up** — Builds Docker images and starts all containers (app, postgres, redis, ngrok, nginx)
5. **migrate** — Runs database migrations inside the app container
6. **test** — Runs the full test suite inside the app container
7. **validate** — Waits for the app to be healthy, checks that ngrok is tunneling correctly, and verifies all public endpoints are reachable

For the first bootstrap, it is normal for Appmax credential warnings to appear because `APPMAX_*` is intentionally blank at this stage.

When it finishes successfully, you see:

```
All validations passed.
Stack is ready.
```

</details>

<details>
<summary><strong>Windows (PowerShell)</strong></summary>

Windows does not have `make` by default. Use the PowerShell install script instead.

**Important**: Scripts downloaded or cloned from GitHub are marked as "remote" by Windows (Zone.Identifier). PowerShell will block them even with `RemoteSigned` policy. You must unblock the scripts first.

```powershell
# Open PowerShell as Administrator

# Navigate to the project root
cd AppStore-Appmax-App-Integration

# Step 1: Allow running locally-created scripts (one-time setting)
Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned

# Step 2: Unblock the install script (Windows marks cloned files as "remote")
Unblock-File -Path .\install.ps1

# Step 3: Run the install script
.\install.ps1
```

**If you get "cannot be loaded because running scripts is disabled":**

```powershell
# Bypass the execution policy for this run only
powershell -ExecutionPolicy Bypass -File .\install.ps1
```

**If the script fails or you prefer manual steps:**

These commands run on your host machine. Commands prefixed with `docker compose exec` execute inside the running container:

```powershell
# 1. Copy env file (if not done yet) — runs on your host
Copy-Item .env.example .env

# 2. Remove old containers — runs on your host, talks to Docker
docker compose down -v --remove-orphans

# 3. Build and start all containers — runs on your host
docker compose up -d --build

# 4. Wait for the app container to be healthy (~60-90 seconds)
#    Open http://localhost:8080/health in your browser to check

# 5. Run database migrations — executes INSIDE the "app" container
docker compose exec -T app ./tmp/server artisan migrate

# 6. Run tests — executes INSIDE the "app" container
docker compose exec -T app sh -lc 'PATH=/usr/local/go/bin:/go/bin:$PATH GOCACHE=/tmp/.gocache go test ./...'
#    The -T flag disables pseudo-TTY allocation (required on Windows to avoid hanging)

# 7. Verify the app is healthy — runs on your host
Invoke-WebRequest -Uri http://localhost:8080/health -UseBasicParsing
# Or open http://localhost:8080/health in your browser
```

</details>

### Step 5: Register the Generated URLs in Appmax

After `make install` (or the PowerShell equivalent) finishes, the terminal prints the public URLs:

```
Frontend URL:  https://your-name.ngrok-free.app/
Health URL:    https://your-name.ngrok-free.app/health
Install URL:   https://your-name.ngrok-free.app/install/start
Callback URL:  https://your-name.ngrok-free.app/integrations/appmax/callback/install
Webhook URL:   https://your-name.ngrok-free.app/webhooks/appmax
```

Use them in Appmax like this:

- **Frontend URL** → `URL do sistema`
- **Install URL** → `Host`
- **Callback URL** → `URL de validação`

Example:

```
Frontend URL: https://your-name.ngrok-free.app/
Install URL: https://your-name.ngrok-free.app/install/start
Callback URL: https://your-name.ngrok-free.app/integrations/appmax/callback/install
```

As of **March 20, 2026**, App Store URL changes are not automatically replicated into the Appmax sandbox. If you change them later, contact `desenvolvimento@appmax.com.br` and ask Appmax to replicate the update.

### Step 6: Fill the Appmax Credentials After Appmax Emails Them

Appmax sends 4 app-level values after validating the application:

| Credential | Format | Example |
|-----------|--------|---------|
| `APPMAX_CLIENT_ID` | OAuth client ID | `abc123` |
| `APPMAX_CLIENT_SECRET` | OAuth client secret | `secret_xyz789` |
| `APPMAX_APP_ID_UUID` | UUID format | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |
| `APPMAX_APP_ID_NUMERIC` | Integer | `42` |

> **Important**: These are two different identifiers for the same app. The UUID is used in the OAuth authorize flow (`/app/authorize`). The numeric ID is used in the health check POST callback from Appmax. Both are required for the full Appmax installation flow.

When Appmax sends those values, fill them in `.env` and reinstall the project:

```bash
make install
```

On Windows PowerShell:

```powershell
.\install.ps1
```

### Step 7: Verify Everything Works

After reinstalling with the Appmax credentials filled in:

| What to Check | URL | Expected Result |
|--------------|-----|----------------|
| App is running | http://localhost:8080 | Frontend welcome page |
| Health check | http://localhost:8080/health | Health status page |
| ngrok inspector | http://localhost:4040 | ngrok web UI showing your public tunnel URL |
| Public health check | `https://<your-ngrok-url>/health` | Same health page, accessible from the internet |

At this stage, the printed **Install URL** and **Callback URL** are not just registered in Appmax; they are also backed by the full Appmax credential configuration in your local `.env`.

---

## Debugging ngrok

ngrok is the most common source of setup issues. Here is how to diagnose and fix problems.

### ngrok Container Exits Immediately

**Symptoms**: `docker compose ps` shows the ngrok container as "Exited".

**Cause**: `NGROK_AUTHTOKEN` in `.env` is empty, missing, or invalid.

**Fix**:
1. Open `.env` and check that `NGROK_AUTHTOKEN` has a value with no spaces around `=`:
   ```
   NGROK_AUTHTOKEN=2abc123def456ghi789
   ```
2. Verify your token is valid at https://dashboard.ngrok.com/get-started/your-authtoken
3. Restart: `docker compose up -d ngrok` (or rerun the full install command: `make install` on macOS/Linux or `.\install.ps1` on Windows PowerShell)

### ngrok Starts but No Tunnel URL

**Symptoms**: ngrok container is running but http://localhost:4040 shows no tunnels.

**Diagnosis**:
```bash
# Check ngrok logs for errors
docker compose logs ngrok
```

**Common errors in logs**:
- `ERR_NGROK_108` — Auth token is expired or revoked. Go to https://dashboard.ngrok.com/get-started/your-authtoken, copy a fresh token, update `.env`, restart
- `ERR_NGROK_120` — Too many simultaneous sessions. Free accounts allow 1 tunnel. Kill other ngrok processes:
  - macOS/Linux: `killall ngrok`
  - Windows: Close any ngrok terminal windows, or use Task Manager to end `ngrok.exe`
- `ERR_NGROK_105` — Invalid auth token format. Check for copy-paste errors (trailing spaces, missing characters)

### Static Domain Not Working

**Symptoms**: You set `NGROK_URL` in `.env` but ngrok uses a random URL instead, or fails to start.

**Fix**:
1. Go to https://dashboard.ngrok.com/domains and verify the domain exists
2. In `.env`, set the domain **without** `https://`:
   ```
   # Correct
   NGROK_URL=your-name.ngrok-free.app

   # Wrong
   NGROK_URL=https://your-name.ngrok-free.app
   ```
3. The domain must match exactly what ngrok assigned (case-sensitive)
4. Restart: `docker compose up -d ngrok`

### Appmax Cannot Reach Your App

**Symptoms**: The install flow fails because Appmax cannot POST to your callback URL.

**Step-by-step diagnosis**:

```bash
# 1. Is your app running?
curl http://localhost:8080/health
# Should return {"status":"ok"} or similar

# 2. Is ngrok running and tunneling?
curl http://localhost:4040/api/tunnels
# Should return JSON with your public_url

# 3. Is the public URL reachable?
curl https://your-name.ngrok-free.app/health
# Should return the same health response as step 1
```

If step 1 fails: the app container is down. Check `docker compose logs app`.

If step 2 fails: ngrok is not running. Check `docker compose logs ngrok`.

If step 3 fails but steps 1 and 2 pass: ngrok is not correctly tunneling to the app. Check the nginx container: `docker compose logs nginx`.

### How to Inspect All ngrok Traffic

The ngrok inspector at http://localhost:4040 shows every HTTP request that passes through the tunnel:

- Full request headers and body
- Full response headers and body
- Timing information
- **Replay button** — re-send any request (useful for debugging webhooks)

This is invaluable for debugging what Appmax is sending to your app and what your app responds with.

---

## Troubleshooting Windows Issues

| Problem | Cause | Fix |
|---------|-------|-----|
| `running scripts is disabled` | PowerShell execution policy blocks cloned scripts | Run `Set-ExecutionPolicy -Scope CurrentUser -ExecutionPolicy RemoteSigned` then `Unblock-File -Path .\install.ps1` |
| Script still blocked after RemoteSigned | Windows Zone.Identifier on cloned files | `Unblock-File -Path .\install.ps1` or run with `powershell -ExecutionPolicy Bypass -File .\install.ps1` |
| `make` not found | Windows doesn't ship GNU Make | Use the PowerShell script, or install Make via `choco install make` (Chocolatey) or `winget install GnuWin32.Make` |
| `docker compose exec` hangs | TTY allocation issue on Windows | Always use `-T` flag: `docker compose exec -T app <command>` |
| `.env` parsing errors | Windows Notepad adds `\r\n` line endings | Edit `.env` with VS Code, or configure git: `git config core.autocrlf input` and re-clone |
| `docker` command not found | Docker Desktop not running or not in PATH | Launch Docker Desktop, close and reopen PowerShell |
| WSL 2 not installed | Docker Desktop requires WSL 2 backend | Run `wsl --install` in admin PowerShell, restart machine |
| `curl` behaves differently | PowerShell aliases `curl` to `Invoke-WebRequest` | Use `Invoke-WebRequest -Uri <url> -UseBasicParsing`, or install real curl: `winget install curl.curl` |

---

## Common Commands

Most `docker compose` commands work natively in PowerShell without any extra tooling.

| Action | Unix/macOS | Windows (PowerShell) |
|--------|-----------|----------------------|
| Full install | `make install` | `.\install.ps1` |
| Start containers | `make up` | `docker compose up -d --build` |
| Stop containers | `make down` | `docker compose down -v --remove-orphans` |
| Restart containers | `make restart` | `docker compose restart` |
| Follow logs | `make logs` | `docker compose logs -f` |
| Run tests | `make test` | `docker compose exec -T app sh -lc 'PATH=/usr/local/go/bin:/go/bin:$PATH GOCACHE=/tmp/.gocache go test ./...'` |
| Run migrations | `make migrate` | `docker compose exec -T app ./tmp/server artisan migrate` |
| Wait for health check | `make health` | `Invoke-WebRequest -Uri http://localhost:8080/health -UseBasicParsing` |
| Verify endpoints via ngrok | `make validate` | _(included in `.\install.ps1`)_ |
| Remove all containers + volumes | `make teardown` | `docker compose down -v --remove-orphans` |
| Rename Go module path | `make rename-module NEW=<path>` | `.\scripts\rename-module.ps1 -NewPath <path>` |

---

## Project Structure

```
app/
  adapters/goravel/       # Framework adapters (Logger, Cache)
  gateway/appmax/         # Appmax HTTP client, retry logic, types
  http/
    controllers/          # HTTP handlers + frontend pages
    middleware/            # MerchantContext auth middleware
    requests/             # Request validation structs
    responses/            # Response structs
  models/                 # GORM models (Installation, Order, WebhookEvent)
  repositories/           # Database access layer
  services/               # Business logic (install, checkout, webhook, token)
bootstrap/                # Dependency injection modules
config/                   # Goravel config (app, database, cache, http, logging)
database/migrations/      # PostgreSQL migrations
docs/                     # Documentation
public/                   # Static assets (CSS, JS, images)
resources/views/frontend/ # Go HTML templates
routes/                   # Route definitions
scripts/                  # Setup and testing scripts
tests/
  unit/                   # Unit tests with hand-written mocks
  integration/            # Integration tests (real PostgreSQL)
```

---

## Documentation

### Setup
- [Local Development](docs/setup/local-development.md) -- Docker, Makefile, Air hot reload, ngrok, scripts
- [How to Get Appmax Credentials and Configure Your App](docs/setup/how-to-get-appmax-credentials-and-configure-your-app.md) -- App Store registration, app creation, permissions, callback setup, credential approval flow
- [Database](docs/setup/database.md) -- Migrations, schema, connection pooling, Redis

### Integration
- [Architecture Guide](docs/integration/guide.md) -- Install flow, checkout flow, webhooks, credentials
- [Endpoints](docs/integration/endpoints.md) -- All 19 HTTP endpoints with request/response details
- [Data Model](docs/integration/data-model.md) -- ER diagram, models, order status lifecycle
- [Frontend Pages](docs/integration/frontend.md) -- Content negotiation, install redirect flow, templates
- [Logging](docs/integration/logging.md) -- Log prefixes, debugging guides, grep patterns
- [Testing](docs/integration/testing.md) -- Test organization, running tests, mocks

### Appmax
- [API Reference](docs/appmax/api-reference.md) -- Upstream Appmax API endpoints and types
- [Webhook Events](docs/appmax/webhooks.md) -- Webhook payload models and event reference
- [Postman Variables](docs/postman/postman-variables.md) -- Postman collection variable reference
