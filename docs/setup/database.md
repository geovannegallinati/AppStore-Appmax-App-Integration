# Database Reference

PostgreSQL schema, migrations, connection pooling, and Redis cache configuration.

---

## Overview

- **PostgreSQL 16** (Alpine) via Docker Compose
- **3 tables**: `installations`, `orders`, `webhook_events`
- **Migrations**: Managed by Goravel framework in `database/migrations/`
- **ORM**: GORM (via Goravel's database provider)
- **Cache**: Redis 7 with key prefix `appmax_checkout_`

---

## Running Migrations

```bash
# Inside Docker (recommended)
make migrate

# Directly via docker compose
docker compose exec app ./tmp/server artisan migrate
```

Migrations run automatically during the full install command (`make install` on macOS/Linux or `.\install.ps1` on Windows PowerShell). The framework tracks executed
migrations in the `migrations` table.

---

## Migrations

| File | Table | Purpose |
|------|-------|---------|
| `20260313000001_create_installations_table.go` | `installations` | Merchant installations with OAuth credentials |
| `20260313000002_create_orders_table.go` | `orders` | Checkout orders with payment details |
| `20260313000003_create_webhook_events_table.go` | `webhook_events` | Webhook event storage and deduplication |

Registered in `database/migrations/migrations.go`:

```go
func All() []schema.Migration {
    return []schema.Migration{
        &CreateInstallationsTable{},
        &CreateOrdersTable{},
        &CreateWebhookEventsTable{},
    }
}
```

---

## Schema

### installations

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | BIGSERIAL | PK | Auto-increment ID |
| `external_key` | VARCHAR(255) | UNIQUE, NOT NULL | Merchant identifier from Appmax |
| `app_id` | VARCHAR(255) | NOT NULL | App ID (numeric form from health check, UUID form from OAuth) |
| `merchant_client_id` | VARCHAR(512) | NOT NULL | Merchant OAuth client ID |
| `merchant_client_secret` | VARCHAR(512) | NOT NULL | Merchant OAuth client secret |
| `external_id` | UUID | UNIQUE, NOT NULL | Auto-generated UUID exposed to external systems |
| `installed_at` | TIMESTAMP | DEFAULT NOW() | When the installation was confirmed |
| `created_at` | TIMESTAMP | — | Record creation time |
| `updated_at` | TIMESTAMP | — | Last update time |

Key design decisions:
- `external_key` is the upsert key — both install paths (OAuth GET and health check POST) use it
- `external_id` is a UUID generated on first creation, never changes after that
- `merchant_client_id` and `merchant_client_secret` are 512 chars to accommodate long OAuth tokens

### orders

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | BIGSERIAL | PK | Auto-increment ID |
| `installation_id` | BIGINT | FK → `installations.id` | Owning installation |
| `appmax_customer_id` | INTEGER | NOT NULL | Appmax customer ID |
| `appmax_order_id` | INTEGER | UNIQUE, NOT NULL | Appmax order ID (dedup key) |
| `status` | VARCHAR(64) | DEFAULT `'pendente'` | Order status (updated by webhooks) |
| `payment_method` | VARCHAR(32) | NULLABLE | `credit_card`, `pix`, or `boleto` |
| `total_cents` | INTEGER | DEFAULT 0 | Order total in cents |
| `pix_qr_code` | TEXT | NULLABLE | Pix QR code image data |
| `pix_emv` | TEXT | NULLABLE | Pix EMV (copy-paste code) |
| `boleto_pdf_url` | TEXT | NULLABLE | Boleto PDF download URL |
| `boleto_digitavel` | TEXT | NULLABLE | Boleto barcode digits |
| `upsell_hash` | VARCHAR(255) | NULLABLE | Upsell transaction hash |
| `created_at` | TIMESTAMP | — | Record creation time |
| `updated_at` | TIMESTAMP | — | Last update time |

Key design decisions:
- `appmax_order_id` is unique — prevents duplicate order records
- Payment-specific fields (pix, boleto) are nullable — only populated for their payment method
- `status` starts as `'pendente'` and is updated by webhook events

### webhook_events

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | BIGSERIAL | PK | Auto-increment ID |
| `event` | VARCHAR(64) | NOT NULL | Event name (e.g., `OrderPaid`, `order_paid`) |
| `event_type` | VARCHAR(32) | NOT NULL | Event category (e.g., `order`, `subscription`) |
| `appmax_order_id` | INTEGER | NULLABLE | Associated order ID (null for non-order events) |
| `payload` | JSON | NOT NULL | Full webhook payload |
| `processed` | BOOLEAN | DEFAULT FALSE | Whether the event has been handled |
| `processed_at` | TIMESTAMP | NULLABLE | When the event was processed |
| `error_message` | TEXT | NULLABLE | Error details if processing failed |
| `created_at` | TIMESTAMP | DEFAULT NOW() | When the event was received |

Key design decisions:
- No `updated_at` — events are write-once, then marked processed
- `appmax_order_id` is nullable because some events (customer, subscription) have no order
- Deduplication: the service checks for existing processed events with the same `event` + `appmax_order_id`

---

## Connection Pooling

Configured in `config/database.go`:

| Setting | Value | Purpose |
|---------|-------|---------|
| `max_idle_conns` | 5 | Idle connections kept open between requests |
| `max_open_conns` | 25 | Maximum simultaneous database connections |
| `conn_max_lifetime` | 300s (5 min) | Maximum time a connection can be reused |

These values are suitable for the current single-instance deployment. For multiple
instances, divide `max_open_conns` by the number of instances to stay within
PostgreSQL's `max_connections` (default 100).

---

## Redis Configuration

Configured in `config/cache.go` and `config/database.go`:

| Setting | Value |
|---------|-------|
| Driver | `redis` (custom via goravel/redis) |
| Key prefix | `appmax_checkout_` |
| Connection | `default` (from `config/database.go` redis section) |

### Cache Keys Used

| Key Pattern | TTL | Purpose |
|-------------|-----|---------|
| `install:{hash}` | 1 hour | Install state (app_id + external_key) during OAuth flow |
| `merchant_token:{installation_id}` | Token expiry minus 60s | Cached merchant bearer token per installation |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_CONNECTION` | `postgres` | Database driver |
| `DB_HOST` | `localhost` | Database host (`postgres` inside Docker) |
| `DB_PORT` | `5432` | Database port |
| `DB_DATABASE` | `appmax_checkout` | Database name |
| `DB_USERNAME` | `appmax` | Database user |
| `DB_PASSWORD` | — | Database password (required) |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `DB_SCHEMA` | `public` | PostgreSQL schema |
| `DB_TIMEZONE` | `UTC` | Connection timezone |
| `REDIS_HOST` | `localhost` | Redis host (`redis` inside Docker) |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_USERNAME` | — | Redis username (optional) |
| `REDIS_PASSWORD` | — | Redis password |
| `REDIS_DB` | `0` | Redis database number |
| `CACHE_DRIVER` | `redis` | Cache backend (`redis` or `memory`) |

---

## Inspecting the Database

```bash
# Connect to PostgreSQL inside Docker
docker compose exec postgres psql -U appmax -d appmax_checkout

# List tables
\dt

# Check installations
SELECT id, external_key, app_id, external_id, installed_at FROM installations;

# Check orders
SELECT id, installation_id, appmax_order_id, status, payment_method FROM orders;

# Check webhook events
SELECT id, event, event_type, appmax_order_id, processed FROM webhook_events ORDER BY created_at DESC LIMIT 10;

# Check migration status
SELECT * FROM migrations;
```
