# Frontend & Install Flow Reference

Developer-facing frontend pages, content negotiation, the full install redirect flow,
and environment-specific behavior.

---

## Overview

The app serves 5 frontend pages as Go HTML templates on the **same routes** as the
JSON API. A browser receives HTML; an API client receives JSON. The pages exist as
developer tools for validating install flow, health checks, callback reachability,
and webhook configuration.

- Templates: `resources/views/frontend/`
- Styles: `/public/css/frontend.css` (dark theme, responsive)
- Scripts: `/public/js/frontend.js` (install form handler, copy-to-clipboard)
- Static assets served via `facades.Route().Static("public", "./public")`

---

## Content Negotiation

Every shared route checks the `Accept` header to decide between HTML and JSON:

```go
func requestWantsHTML(ctx contractshttp.Context) bool {
    return strings.Contains(strings.ToLower(ctx.Request().Header("Accept")), "text/html")
}
```

File: `app/http/controllers/frontend_page.go:56-58`

| Scenario | Accept Header | Response |
|----------|---------------|----------|
| Browser navigation | `text/html, ...` | HTML page |
| API client / cURL (default) | `application/json` or absent | JSON |
| Postman (explicit) | Set `Accept: text/html` | HTML page |

Routes that apply content negotiation:
- `GET /health` (`HealthController.Check`)
- `GET /install/start` (`InstallController.Start`)
- `GET /integrations/appmax/callback/install` (`InstallController.CallbackGuide`)

Routes that always return HTML (GET only):
- `GET /` (`HealthController.RootFrontend`)
- `GET /webhooks/appmax` (`WebhookController.Guide`)

---

## Frontend Pages

| Route | Page Kind | Template | Controller Method | Purpose |
|-------|-----------|----------|-------------------|---------|
| `GET /` | root | `page.tmpl` | `HealthController.RootFrontend` | Welcome page with Install / Webhook buttons |
| `GET /health` | health | `page.tmpl` | `HealthController.HealthFrontend` | Health check visual validation |
| `GET /install/start` | install | `page.tmpl` | `InstallController.InstallStartFrontend` | Install form with auto-generated `external_key` |
| `GET /integrations/appmax/callback/install` | callback or completed | `page.tmpl` or `install_completed.tmpl` | `InstallController.CallbackGuide` | Callback readiness page or success page |
| `GET /webhooks/appmax` | webhook | `page.tmpl` | `WebhookController.Guide` | Webhook setup guide with Apphook link |

All pages share the same visual structure from `page.tmpl`:
1. Appmax logo
2. Status badge (Install, Health, Callback, Webhook)
3. Headline and messages
4. Action buttons (conditional)
5. Install form or Webhook guide (conditional)
6. Available endpoints list (current route highlighted)
7. Tips section

The `install_completed.tmpl` is a standalone compact page with only the success message.

---

## Templates

### `page.tmpl` (Main Template)

File: `resources/views/frontend/page.tmpl`

Receives a `frontendPageData` struct:

```go
type frontendPageData struct {
    Title               string
    Badge               string
    Headline            string
    Message             string
    Submessage          string
    ActiveRoute         string
    PageKind            string
    Endpoints           []frontendEndpoint
    Tips                []string
    Buttons             []frontendAction
    ShowInstallForm     bool
    InstallFormAction   string
    DefaultInstallAppID string
    DefaultExternalKey  string
    ShowWebhookGuide    bool
    AppmaxEnvironment   string
    AppmaxAdminURL      string
    AppmaxApphookURL    string
    WebhookEndpointURL  string
}
```

Conditional sections controlled by:
- `ShowInstallForm` — renders the install button and hidden form fields
- `ShowWebhookGuide` — renders the Apphook registration guide with copy button
- `Buttons` — renders action buttons (e.g., "Install Appmax", "Open Apphook")
- `Tips` — renders quick-tips list

### `install_completed.tmpl` (Success Page)

File: `resources/views/frontend/install_completed.tmpl`

Minimal page with:
- Appmax logo
- "Installation completed" badge
- "Success. Installation is complete." headline
- "The AppMax integration token was confirmed successfully."
- "You can safely close this tab now."

No endpoints list, no tips, no scripts.

---

## Install Flow — Step by Step

This section describes the **full Appmax installation flow after `APPMAX_*` has been filled in**. During the initial bootstrap described in the setup docs, the project is started first to obtain the public ngrok URLs for Appmax registration, not to complete the full Appmax install/auth flow yet.

### Step 1: User Opens Install Page

**URL**: `/install/start` (browser)

The browser sends `Accept: text/html`. `InstallController.Start()` detects both
params are missing and the client wants HTML, so it calls `InstallStartFrontend()`.

The form is rendered with:
- `app_id` pre-filled from `APPMAX_APP_ID_UUID` (hidden input)
- `external_key` left empty (JavaScript generates it on submit)

### Step 2: User Clicks "Install Appmax"

`frontend.js` intercepts the form submit:

```js
var externalKeyField = installForm.querySelector('input[name="external_key"]');
externalKeyField.value = Date.now() + '-' + crypto.randomUUID();
```

The browser navigates to:
```
/install/start?app_id={APPMAX_APP_ID_UUID}&external_key={timestamp}-{uuid}
```

### Step 3: Backend Processes Install Start

Now both `app_id` and `external_key` are present. The API flow runs:

1. Calls `appmaxSvc.Authorize(appID, externalKey, callbackURL)` where
   `callbackURL = {APP_URL}/integrations/appmax/callback/install`
2. Receives a `hash` token from Appmax
3. Caches `installState{AppID, ExternalKey}` at key `install:{hash}` with **1-hour TTL**
4. Returns **HTTP 302** redirect

### Step 4: Redirect to Appmax Admin

The redirect URL is environment-dependent:

| Environment | Redirect URL |
|-------------|-------------|
| **Sandbox** | `https://breakingcode.sandboxappmax.com.br/appstore/integration/{hash}` |
| **Production** | `https://admin.appmax.com.br/appstore/integration/{hash}` |

The user lands on the Appmax admin panel and completes the OAuth authorization flow.

### Step 5: Appmax Redirects Back (GET Callback)

Appmax redirects the browser to:
```
{APP_URL}/integrations/appmax/callback/install?token={hash}
```

`CallbackGuide()` processes the callback (exact flow in `install_controller.go:91-157`):

1. Reads Redis cache for `install:{hash}`
2. **Cache miss** (consumed, expired, or invalid token): returns `200 "installation confirmed"`
   immediately. Does not check the database. Trusts that Path B will handle it.
   Renders `install_completed.tmpl` if browser, JSON otherwise.
3. **Cache hit**: consumes the entry (deletes it — each hash works once), then validates
   `AppID` matches `APPMAX_APP_ID_UUID`
4. Calls `appmaxSvc.GenerateMerchantCreds(token)` to obtain merchant `client_id` and `client_secret`
5. **If generate fails**: checks database for existing installation by `external_key`.
   If found with credentials (Path B already completed), shows success page.
   If not found, returns 502 error.
6. **If generate succeeds**: calls `installSvc.Upsert()` to save the installation,
   renders `install_completed.tmpl` — the user sees the success page

### Step 6: Appmax Health Check POST (Concurrent)

Appmax sends a POST to the same callback URL with credentials:

```json
{
  "app_id": "123",
  "external_key": "install-1710...",
  "client_key": "install-1710...",
  "client_id": "merchant_id_value",
  "client_secret": "merchant_secret_value"
}
```

`Callback()` validates:
- All 5 fields present
- `app_id == APPMAX_APP_ID_NUMERIC` (numeric form, not UUID)
- `client_key == external_key` (security check)

Then upserts the installation. The upsert is idempotent by `external_key`, so it is
safe regardless of whether Step 5 ran first.

Returns: `{"external_id": "uuid-value"}`

### Race Condition Handling

Steps 5 and 6 run concurrently. Three possible orderings:

| Order | What Happens |
|-------|-------------|
| GET first, then POST | GET consumes cache, generates creds, creates installation. POST upserts (updates same record). |
| POST first, then GET | POST creates installation. GET finds cache miss, returns 200 "installation confirmed" immediately (does not check DB). |
| Simultaneous | Both paths may call `Upsert` — idempotent by `external_key`, last write wins with same credentials. |

---

## Environment Detection & URLs

### Environment Name

```go
func appmaxEnvironmentName(adminURL string) string {
    if strings.Contains(strings.ToLower(adminURL), "sandbox") {
        return "Sandbox"
    }
    return "Production"
}
```

Displayed on the webhook guide page as "Detected environment: **Sandbox**" or "**Production**".

### Apphook URL

```go
func apphookURLFromAdmin(adminURL string) string {
    return fmt.Sprintf("%s/v2/apphook", normalized)
}
```

| Environment | Apphook URL |
|-------------|-------------|
| Sandbox | `https://breakingcode.sandboxappmax.com.br/v2/apphook` |
| Production | `https://admin.appmax.com.br/v2/apphook` |

### Environment Variables

| Variable | Sandbox | Production (default) |
|----------|---------|---------------------|
| `APPMAX_ADMIN_URL` | `https://breakingcode.sandboxappmax.com.br` | `https://admin.appmax.com.br` |
| `APPMAX_AUTH_URL` | `https://auth.sandboxappmax.com.br` | `https://auth.appmax.com.br` |
| `APPMAX_API_URL` | `https://api.sandboxappmax.com.br` | `https://api.appmax.com.br` |
| `NGROK_URL` or `APP_URL` | Your tunnel URL (e.g., `https://abc123.ngrok.io`) | Your production URL |

### Base URL Resolution

`frontendBaseURL()` resolves the public base URL for generating absolute endpoint links.

Priority order:
1. Parsed from `ctx.Request().Url()`
2. Parsed from `ctx.Request().FullUrl()`
3. `Host` header + `X-Forwarded-Proto` header (proxy/tunnel scenario)
4. Fallback to `APP_URL` / `NGROK_URL`

---

## Webhook Setup Guide

The `GET /webhooks/appmax` page provides a visual guide for registering the webhook
endpoint in Appmax's Apphook system.

The page displays:
- **Detected environment**: Sandbox or Production (based on `APPMAX_ADMIN_URL`)
- **Admin base URL**: The configured `APPMAX_ADMIN_URL` value
- **Apphook link**: Direct link to `{APPMAX_ADMIN_URL}/v2/apphook`
- **Webhook endpoint**: The full URL to register (e.g., `https://abc123.ngrok.io/webhooks/appmax`)
- **Copy button**: One-click copy of the webhook endpoint URL
- **"Open Apphook" button**: Opens the Apphook registration page in a new tab

The POST route (`POST /webhooks/appmax`) continues to handle Appmax webhook events
independently of the frontend guide.

---

## Static Assets

| Asset | Path | Purpose |
|-------|------|---------|
| CSS | `/public/css/frontend.css` | Dark theme, responsive grid, buttons, forms, endpoint list |
| JS | `/public/js/frontend.js` | Install form submit handler, clipboard copy |
| Logo | `/public/images/appmax-logo.jpg` | Appmax brand logo |

### JavaScript Behavior

**Install form** (`frontend.js`):
- Listens for `.install-form` submit event
- Prevents default form submission
- Generates `external_key` as `{Date.now()}-{crypto.randomUUID()}`
- Builds URL with query params and navigates via `window.location.href`

**Copy button** (`frontend.js`):
- Listens for click on `[data-copy-target]` buttons
- Reads text from the target element
- Copies to clipboard via `navigator.clipboard.writeText()`
- Shows "Copied" feedback for 1.2 seconds

---

## Helper Functions Reference

| Function | File:Line | Purpose |
|----------|-----------|---------|
| `requestWantsHTML(ctx)` | `frontend_page.go:56` | Returns `true` if `Accept` header contains `text/html` |
| `frontendBaseURL(ctx, fallback)` | `frontend_page.go:60` | Resolves public base URL from request or env fallback |
| `parseBaseURL(raw)` | `frontend_page.go:85` | Extracts `scheme://host` from a URL string |
| `absoluteURL(baseURL, path)` | `frontend_page.go:99` | Joins base URL with an endpoint path |
| `frontendEndpoints(baseURL, activeRoute)` | `frontend_page.go:111` | Generates the 5-endpoint list with active highlighting |
| `appmaxEnvironmentName(adminURL)` | `frontend_page.go:153` | Returns "Sandbox" or "Production" based on admin URL |
| `apphookURLFromAdmin(adminURL)` | `frontend_page.go:161` | Appends `/v2/apphook` to the admin URL |

All helpers are in `app/http/controllers/frontend_page.go`.
