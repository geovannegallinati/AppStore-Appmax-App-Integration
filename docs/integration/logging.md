# Logging Reference

Logging architecture, configuration, log prefixes, and debugging guides.

---

## Configuration

Logging is configured in `config/logging.go` using the Goravel framework.

| Env Variable   | Default                       | Description                              |
|----------------|-------------------------------|------------------------------------------|
| `LOG_CHANNEL`  | `stack`                       | Log channel (uses `single` driver)       |
| `LOG_PATH`     | `storage/logs/goravel.log`    | Log file path                            |
| `LOG_LEVEL`    | `debug`                       | Minimum log level                        |
| `LOG_PRINT`    | `true`                        | Also print to console (stdout)           |

Log retention: 14 days.

---

## Log Levels

| Level      | Usage                                                                       | Count |
|------------|-----------------------------------------------------------------------------|-------|
| `Debugf`   | Low-level tracing: HTTP calls, webhook payloads, cache misses, parsed fields | 6     |
| `Infof`    | Significant events: raw health check payloads, existing installation found   | 2     |
| `Warningf` | Recoverable issues: best-effort failures, unmapped events, cache write errors| 6     |
| `Errorf`   | Failures requiring attention: validation errors, upstream failures, DB errors | 26    |

---

## Logger Architecture

Two logging mechanisms coexist:

### 1. `facades.Log()` (Goravel Framework)

Used directly in controllers and the gateway HTTP client. Supports all levels:
`Debugf`, `Infof`, `Warningf`, `Errorf`.

### 2. `Logger` Interface (Service Layer)

Defined in `app/services/ports.go`:
```go
type Logger interface {
    Warningf(format string, args ...any)
    Errorf(format string, args ...any)
}
```

Only `Warningf` and `Errorf` are available at the service layer (no `Debugf`/`Infof`).
Services accept an optional Logger dependency; if nil, a no-op logger is used.

**Adapter**: `app/adapters/goravel/logger.go` wraps `facades.Log()` to implement
the `Logger` interface.

---

## Log Prefixes Quick Reference

| Prefix                       | Component                  | Level(s)         | What It Logs                                        |
|------------------------------|----------------------------|------------------|-----------------------------------------------------|
| `[install:debug]`            | InstallController.Callback | Info             | Raw health check request body and headers            |
| `[install]`                  | InstallController.Callback | Debug            | Parsed health check fields, JSON bind errors         |
| `install_controller:`        | InstallController          | Error/Warning/Info | OAuth failures, cache errors, app_id mismatches, upsert errors |
| `checkout_controller:`       | CheckoutController         | Error            | All checkout operation failures                      |
| `webhook_controller:`        | WebhookController          | Error/Debug      | Handle failures, incoming webhook event payloads     |
| `merchant_auth_controller:`  | MerchantAuthController     | Error            | Merchant token fetch failures                        |
| `webhook_service:`           | WebhookService             | Warning/Error    | Unmapped events, order not found, mark-processed failures |
| `checkout_service:`          | CheckoutService            | Warning          | Best-effort order persistence failures               |
| `token_manager:`             | TokenManager               | Warning          | Redis cache write failures for merchant tokens       |
| `[appmax]`                   | Gateway HTTP Client        | Debug/Error      | Outbound HTTP requests, responses, retries, status errors |
| `routes api:`                | routes/api.go              | Error            | Bootstrap dependency initialization failures         |

---

## Component-by-Component Log Inventory

### Install Controller

File: `app/http/controllers/install_controller.go` (16 log statements)

#### Start Method (GET /install/start)

| Level   | Prefix                | Message                                           | Fields Logged                 |
|---------|-----------------------|---------------------------------------------------|-------------------------------|
| Error   | `install_controller:` | `authorize failed for key %s: %v`                 | external_key, error           |
| Error   | `install_controller:` | `marshal state failed for hash %s: %v`            | hash, error                   |
| Error   | `install_controller:` | `cache put failed for hash %s: %v`                | hash, error                   |

#### CallbackGuide Method (GET /integrations/appmax/callback/install)

| Level   | Prefix                | Message                                           | Fields Logged                 |
|---------|-----------------------|---------------------------------------------------|-------------------------------|
| Debug   | `install_controller:` | `no cached state for token %s — installation will be confirmed via health check POST` | token |
| Error   | `install_controller:` | `unmarshal state failed for token %s: %v`         | token, error                  |
| Error   | `install_controller:` | `app_id mismatch in state for token %s: got %s`   | token, app_id                 |
| Warning | `install_controller:` | `generate merchant creds failed for token %s: %v — checking for existing installation` | token, error |
| Info    | `install_controller:` | `installation already exists for key %s (created by POST callback)` | external_key |
| Error   | `install_controller:` | `generate merchant creds failed for token %s: %v` | token, error                  |
| Error   | `install_controller:` | `upsert failed for key %s: %v`                    | external_key, error           |

#### Callback Method (POST /integrations/appmax/callback/install — Health Check)

| Level   | Prefix                | Message                                           | Fields Logged                           |
|---------|-----------------------|---------------------------------------------------|-----------------------------------------|
| Info    | `[install:debug]`     | `raw callback from %s — headers: %v \| body: %s`  | client IP, all headers, raw JSON body   |
| Debug   | `[install]`           | `healthcheck POST from %s — bind error: %v`       | client IP, bind error                   |
| Debug   | `[install]`           | `healthcheck POST from %s — app_id=%s external_key=%s client_key=%s client_id=%s` | client IP, all parsed fields |
| Error   | `install_controller:` | `app_id mismatch for key %s: got %s`              | external_key, received app_id           |
| Error   | `install_controller:` | `client_key mismatch for key %s: got %s`          | external_key, received client_key       |
| Error   | `install_controller:` | `upsert failed for key %s: %v`                    | external_key, error                     |

---

### Checkout Controller

File: `app/http/controllers/checkout_controller.go` (10 log statements)

All statements are `Errorf` with prefix `checkout_controller:`:

| Method         | Message                          | Fields Logged |
|----------------|----------------------------------|---------------|
| CreateOrder    | `create order failed: %v`        | error         |
| PayCreditCard  | `credit card failed: %v`         | error         |
| PayPix         | `pix failed: %v`                 | error         |
| PayBoleto      | `boleto failed: %v`              | error         |
| Status         | `get status failed: %v`          | error         |
| Installments   | `installments failed: %v`        | error         |
| Refund         | `refund failed: %v`              | error         |
| Tokenize       | `tokenize failed: %v`            | error         |
| AddTracking    | `tracking failed: %v`            | error         |
| Upsell         | `upsell failed: %v`              | error         |

---

### Webhook Controller

File: `app/http/controllers/webhook_controller.go` (3 log statements)

| Level   | Prefix                  | Message                                           | Fields Logged                                   |
|---------|-------------------------|---------------------------------------------------|-------------------------------------------------|
| Error   | `webhook_controller:`   | `handle failed for event %s: %v`                  | event name, error                               |
| Debug   | `webhook_controller:`   | `received event=%s event_type=%s order_id=%s model=%s payload_unmarshalable=true payload_data=%s` | event, event_type, order_id, model, raw data |
| Debug   | `webhook_controller:`   | `received event=%s event_type=%s order_id=%s model=%s payload=%s` | event, event_type, order_id, model, full JSON |

---

### Merchant Auth Controller

File: `app/http/controllers/merchant_auth_controller.go` (1 log statement)

| Level   | Prefix                        | Message                                           | Fields Logged           |
|---------|-------------------------------|---------------------------------------------------|-------------------------|
| Error   | `merchant_auth_controller:`   | `merchant token fetch failed for key %s: %v`      | external_key, error     |

---

### Gateway HTTP Client

File: `app/gateway/appmax/http.go` (4 log statements)

| Level   | Prefix     | Message                                                | Fields Logged                                      |
|---------|------------|--------------------------------------------------------|----------------------------------------------------|
| Debug   | `[appmax]` | `→ %s %s`                                              | HTTP method, endpoint URL (outbound request)       |
| Debug   | `[appmax]` | `← %s %s status=%d`                                    | HTTP method, endpoint URL, status code (response)  |
| Debug   | `[appmax]` | `retrying after status %d (attempt %d/%d)`              | status code, attempt number, max attempts          |
| Error   | `[appmax]` | `unexpected status %d [trace-headers]: [body]`          | status code, CF-Ray, X-Request-Id, X-Trace-Id, response body |

The error log for unexpected statuses includes Cloudflare and Appmax trace headers
when available, which are useful for contacting Appmax support.

---

### Webhook Service

File: `app/services/webhook_service.go` (3 log statements, via `Logger` interface)

| Level   | Prefix              | Message                                                 | Fields Logged       |
|---------|---------------------|---------------------------------------------------------|---------------------|
| Warning | `webhook_service:`  | `event %s has no status mapping or no order_id, marking processed` | event name   |
| Warning | `webhook_service:`  | `order %d not found for event %s`                       | order_id, event name|
| Error   | `webhook_service:`  | `failed to mark event %d as processed: %v`              | event ID, error     |

---

### Checkout Service

File: `app/services/checkout_service.go` (1 log statement, via `Logger` interface)

| Level   | Prefix               | Message                                   | Fields Logged       |
|---------|----------------------|-------------------------------------------|---------------------|
| Warning | `checkout_service:`  | `failed to persist order %d: %v`          | appmax_order_id, error |

This is a **best-effort** log: the payment has already succeeded, but the local
database write failed. The payment response is still returned to the client.

---

### Token Manager

File: `app/services/token_manager.go` (1 log statement, via `Logger` interface)

| Level   | Prefix            | Message                                                  | Fields Logged           |
|---------|-------------------|----------------------------------------------------------|-------------------------|
| Warning | `token_manager:`  | `failed to cache merchant token for installation %d: %v` | installation_id, error  |

This is a **best-effort** log: the token was fetched successfully from Appmax, but
caching to Redis failed. The token is still returned to the caller.

---

### Routes Bootstrap

File: `routes/api.go` (1 log statement)

| Level   | Prefix        | Message                                    | Fields Logged |
|---------|---------------|--------------------------------------------|---------------|
| Error   | `routes api:` | `bootstrap http dependencies failed: %v`   | error         |

This log indicates a fatal startup error — the app will not serve any routes.

---

## Debugging Guides

### How do I debug a failed installation?

Search logs in this order:

1. **`[install:debug] raw callback from`** — Shows the raw HTTP body and headers that
   Appmax sent. Verify `app_id` is numeric, `client_key` matches `external_key`, and
   all 5 fields are present.
2. **`[install] healthcheck POST from`** — Shows parsed fields after JSON binding.
   If this line is missing but the raw log exists, the JSON body couldn't be parsed.
3. **`install_controller: app_id mismatch`** — `APPMAX_APP_ID_NUMERIC` env var doesn't
   match what Appmax sent.
4. **`install_controller: client_key mismatch`** — `client_key` != `external_key` in
   the Appmax payload.
5. **`install_controller: upsert failed`** — Database error during installation save.

If no health check logs appear at all, the callback URL is not reachable from Appmax.
Check `NGROK_URL` or `APP_URL` configuration.

---

### How do I trace an Appmax API call?

Every outbound HTTP call to Appmax is logged with the `[appmax]` prefix:

```
[appmax] → POST https://api.sandboxappmax.com.br/v1/customers    (outbound)
[appmax] ← POST https://api.sandboxappmax.com.br/v1/customers status=200  (response)
```

If the call fails:
```
[appmax] unexpected status 422 CF-Ray=abc123 X-Request-Id=xyz: {"message":"validation error","errors":{"email":["required"]}}
```

The `CF-Ray` and `X-Request-Id` values can be provided to Appmax support for
request tracing.

If retries are configured:
```
[appmax] retrying after status 502 (attempt 1/3)
```

---

### How do I debug a webhook?

1. **`webhook_controller: received event=OrderPaid event_type=order order_id=12345 model=standard payload={...}`**
   — Shows the full incoming webhook with detected payload model and extracted order_id.
2. **`webhook_service: event X has no status mapping`** — The event is unknown or has
   no order_id. Check `webhookStatusMap` in `app/services/webhook_service.go`.
3. **`webhook_service: order %d not found`** — The webhook references an order that
   doesn't exist locally (may not have been persisted during payment).
4. **`webhook_controller: handle failed`** — Processing error (DB failure, etc.).

---

### How do I debug a payment failure?

1. **`checkout_controller: credit card failed: %v`** (or `pix failed`, `boleto failed`)
   — The controller-level error, which wraps the underlying cause.
2. **`[appmax] → POST .../v1/payments/credit-card`** — The outbound request to Appmax.
3. **`[appmax] ← ... status=422`** or **`[appmax] unexpected status ...`** — The Appmax
   response with error details.
4. **`checkout_service: failed to persist order`** — (Warning) If this appears, the
   payment succeeded but local DB write failed. The client still got the success response.

---

### How do I debug token issues?

1. **`merchant_auth_controller: merchant token fetch failed`** — Token fetch from Appmax
   failed for this installation.
2. **`token_manager: failed to cache merchant token`** — (Warning) Token was fetched but
   Redis caching failed. Next request will re-fetch from Appmax.
3. **`[appmax] → POST .../oauth2/token`** — The actual OAuth token request to Appmax.
4. **`[appmax] unexpected status 401`** — Invalid credentials. Check `merchant_client_id`
   and `merchant_client_secret` in the `installations` table.

---

### How do I know if the app started successfully?

If `routes api: bootstrap http dependencies failed` appears in the logs, the app
failed to initialize. Common causes:
- Missing required env vars (`APPMAX_CLIENT_ID`, `APPMAX_CLIENT_SECRET`,
  `APPMAX_APP_ID_UUID`, `APPMAX_APP_ID_NUMERIC`, `APP_URL` or `NGROK_URL`)
- Database connection failure
- Redis connection failure

If this log does NOT appear, the app bootstrapped successfully.

---

## Log Grep Patterns

Quick grep patterns for common investigations:

```bash
# All errors
grep "local\.error:" storage/logs/goravel.log

# Health check debugging
grep "\[install" storage/logs/goravel.log

# All Appmax API calls
grep "\[appmax\]" storage/logs/goravel.log

# Failed Appmax API calls only
grep "\[appmax\] unexpected" storage/logs/goravel.log

# All webhook events received
grep "webhook_controller: received" storage/logs/goravel.log

# All payment failures
grep "checkout_controller:" storage/logs/goravel.log

# Best-effort warnings (non-critical)
grep "failed to persist order\|failed to cache merchant" storage/logs/goravel.log
```
