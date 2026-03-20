# Postman Collection Variables

This document describes all 36 variables in the **AppMax — Full Integration Suite** collection, grouped by who is responsible for them.

---

## 1. Fill in Manually

These variables are environment-specific and must be set before running any requests.

Use them in two phases:

1. **Initial bootstrap** — set `NGROK_URL` and `BASE_URL`; leave `APPMAX_*` blank until Appmax emails them.
2. **Full Appmax flow** — after Appmax sends the credential email, fill `APPMAX_CLIENT_ID`, `APPMAX_CLIENT_SECRET`, `APPMAX_APP_ID_UUID`, and `APP_ID_NUMERIC`.

| Variable | Mapped from `.env` | Notes |
| --- | --- | --- |
| `NGROK_URL` | `NGROK_URL` | Changes every ngrok session. Update whenever you restart ngrok. Include the `https://` scheme (e.g., `https://foo.ngrok-free.app`). The collection pre-request script adds `https://` automatically if the scheme is missing and strips trailing slashes. |
| `BASE_URL` | `APP_HOST` + `APP_PORT` | Default `http://localhost:8080` — change only if the server runs on a different port. |
| `APPMAX_CLIENT_ID` | `APPMAX_CLIENT_ID` | Leave blank until Appmax sends the credential email after the first URL registration. |
| `APPMAX_CLIENT_SECRET` | `APPMAX_CLIENT_SECRET` | Leave blank until Appmax sends the credential email after the first URL registration. |
| `APPMAX_APP_ID_UUID` | `APPMAX_APP_ID_UUID` | Leave blank until Appmax sends the credential email after the first URL registration. |
| `APP_ID_NUMERIC` | `APPMAX_APP_ID_NUMERIC` | Leave blank until Appmax sends the credential email after the first URL registration. |

---

## 2. Pre-configured by Postman — Do Not Remove

These variables are hardcoded with stable sandbox infrastructure URLs and shared defaults. Do not delete or modify them unless the sandbox environment changes.

| Variable | Value | Purpose |
| --- | --- | --- |
| `AUTH_URL` | `https://auth.sandboxappmax.com.br` | Appmax Keycloak — app OAuth2 token endpoint |
| `API_URL` | `https://api.sandboxappmax.com.br` | Appmax REST API base |
| `REDIRECT_BASE` | `https://breakingcode.sandboxappmax.com.br` | Merchant browser redirect base |
| `CUSTOMER_ID` | `6933` | Default customer for the Appmax direct flow — overwritten by Create Customer requests. |
| `AUTO_SYNC_MERCHANT_TOKEN` | `true` | Controls whether `MERCHANT_TOKEN` is refreshed automatically before each request that uses it. Set to `false` to disable. |

### NGROK_URL Normalization

A collection-level pre-request script runs before every request and normalizes `NGROK_URL`:

- Strips trailing slashes
- Adds `https://` if no scheme is present (e.g., `foo.ngrok-free.app` → `https://foo.ngrok-free.app`)

This means you can set `NGROK_URL` with or without the scheme — the script handles it.
`http://` URLs are left as-is (not automatically upgraded to `https://`).

---

## 3. Generated Automatically by Scripts

These variables are populated by pre-request and test scripts during request execution. Do not set them manually — they will be overwritten.

### 3a. Installation and Authentication

Set once at the beginning of each flow run.

| Variable | Set by | When |
| --- | --- | --- |
| `EXTERNAL_KEY` | Appmax Step 1 / Localhost Step 1 (pre-request) | Random key generated at flow start (`postman-<timestamp>` or `local-install-<timestamp>`) |
| `INSTALLATION_KEY` | Appmax Step 1 / Localhost Step 1 (pre-request) | Same value as `EXTERNAL_KEY` — used to identify the installation in localhost endpoints |
| `APP_TOKEN` | Appmax Step 1 (test) | Bearer token for the app, obtained via OAuth2 |
| `HASH` | Appmax Step 2 / Localhost Step 1 (test) | Installation hash returned by Appmax — used to build the browser authorization URL |
| `_BROWSER_URL` | Step 3 (pre-request) | Full browser URL for merchant authorization (internal helper, not used in requests) |
| `MERCHANT_TOKEN` | Appmax Step 6 / Localhost Step 5 (test) + global auto-sync pre-request | Merchant bearer token, refreshed automatically before each request that uses `{{MERCHANT_TOKEN}}` |
| `MERCHANT_CLIENT_ID` | Localhost Step 5 — Sync Merchant Token (test) | Pulled from the local DB via the sync endpoint |
| `MERCHANT_CLIENT_SECRET` | Localhost Step 5 — Sync Merchant Token (test) | Pulled from the local DB via the sync endpoint |

### 3b. Payment Flows and Orders

Set progressively as you run through the payment flows.

| Variable | Set by |
| --- | --- |
| `CUSTOMER_ID` | Create Customer requests |
| `ORDER_ID` | Create Order requests |
| `CARD_TOKEN` | Tokenize Card requests |
| `TOKEN` | Checkout tokenize steps |
| `APPROVED_ORDER_ID` | Approved payment flow |
| `APPROVED_ORDER_ID_2` | Second approved order (split flows) |
| `ORDER_ID_DECLINED` | Declined payment flow |
| `ORDER_ID_PIX` | PIX payment flow |
| `PIX_ORDER_ID` | PIX order creation |
| `ORDER_ID_BOLETO` | Boleto payment flow |
| `BOLETO_ORDER_ID` | Boleto order creation |
| `ORDER_ID_PARTIAL` | Partial payment flow |
| `ORDER_ID_SUB_CC` | Subscription credit card flow |
| `ORDER_ID_SUB_PIX` | Subscription PIX flow |
| `UPSELL_ORDER_ID` | Upsell order |
| `UPSELL_HASH` | Upsell request |
| `RECIPIENT_HASH` | Recipient management flow |
| `ABANDONED_CUSTOMER_ID` | Abandoned cart flow |
