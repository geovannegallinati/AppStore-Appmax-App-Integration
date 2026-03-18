# App Endpoints Reference

Detailed reference for all endpoints exposed by **our app** (not Appmax upstream endpoints).
For Appmax API endpoints, see [appmax/api-reference.md](../appmax/api-reference.md).

Routes: `routes/api.go`

---

## Public Endpoints

### GET /

Root frontend page. Returns HTML when `Accept: text/html` is present.

**Controller**: `HealthController.RootFrontend`

---

### GET /health

Health check endpoint.

**Response** (200):
```json
{"status": "ok"}
```

Returns HTML health page when `Accept: text/html` is present.

**Controller**: `HealthController.Check`

---

### GET /install/start

Initiate the installation flow. Validates the app ID, calls Appmax `/app/authorize`,
caches state in Redis, and redirects the merchant to the Appmax admin panel.

**Query Parameters**:

| Parameter      | Type   | Required | Description                                   |
|----------------|--------|----------|-----------------------------------------------|
| `app_id`       | string | Yes      | Must be non-empty. Passed to Appmax `/app/authorize` which expects UUID format |
| `external_key` | string | Yes      | Developer-provided merchant identifier         |

**Success Response**: `302 Redirect` to `{ADMIN_URL}/appstore/integration/{hash}`

**Error Responses**:

| Status | Condition             | Body                                |
|--------|-----------------------|-------------------------------------|
| 400    | Missing `app_id`      | `{"message": "app_id is required"}` |
| 400    | Missing `external_key`| `{"message": "external_key is required"}` |
| 502    | Appmax authorize fails| `{"message": "failed to initiate installation"}` |

When `Accept: text/html` is present and parameters are missing, returns an HTML form instead of JSON errors.

**Controller**: `InstallController.Start`

---

### GET /integrations/appmax/callback/install

OAuth callback. Called when Appmax redirects the merchant back after authorization.
Retrieves cached install state, generates merchant credentials via Appmax API, and upserts
the installation.

**Query Parameters**:

| Parameter | Type   | Required | Description                    |
|-----------|--------|----------|--------------------------------|
| `token`   | string | Yes      | Installation hash from Appmax  |

**Success Response** (200):
```json
{"external_id": "generated-uuid"}
```

**Other Responses**:

| Status | Condition                                         | Body                                                        |
|--------|---------------------------------------------------|-------------------------------------------------------------|
| 200    | Cache miss (token invalid, expired, or consumed)  | `{"message": "installation confirmed"}` (graceful fallback) |
| 200    | Generate fails but installation exists in DB      | `{"external_id": "uuid"}` (Path B already completed)        |
| 400    | Missing `token`                                   | `{"message": "token is required"}`                          |
| 400    | `app_id` mismatch in cached state                 | `{"message": "invalid app_id"}`                             |
| 500    | Cached state JSON parse failure                   | `{"message": "internal server error"}`                      |
| 502    | Merchant cred generation fails + no existing      | `{"message": "failed to generate merchant credentials"}`    |

**Cache miss handling**: If the cache entry is missing (consumed, expired, or invalid token),
the controller returns 200 immediately without checking the database. It trusts that
Path B (health check POST) will handle or has already handled the installation.
The cache entry is consumed (deleted) on first successful read — each hash can only be
used once. Returns HTML completion page when `Accept: text/html` is present.

**Controller**: `InstallController.CallbackGuide`

---

### POST /integrations/appmax/callback/install

Health check callback from Appmax. Called by Appmax to confirm the installation and deliver
merchant credentials directly.

**Request Body** (JSON):

| Field           | JSON Key        | Type   | Required | Description                              |
|-----------------|-----------------|--------|----------|------------------------------------------|
| AppID           | `app_id`        | string | Yes      | Must match `APPMAX_APP_ID_NUMERIC`       |
| ExternalKey     | `external_key`  | string | Yes      | Merchant identifier                      |
| ClientKey       | `client_key`    | string | Yes      | Must equal `external_key`                |
| MerchantClientID| `client_id`     | string | Yes      | OAuth client ID for merchant API access  |
| MerchantClientSecret | `client_secret` | string | Yes | OAuth client secret                      |

**Success Response** (200):
```json
{"external_id": "generated-uuid"}
```

**Error Responses**:

| Status | Condition                           | Body                                                                           |
|--------|-------------------------------------|--------------------------------------------------------------------------------|
| 400    | Missing required fields             | `{"message": "app_id, external_key, client_key, merchant_client_id and merchant_client_secret are required"}` |
| 400    | `app_id` != `APPMAX_APP_ID_NUMERIC` | `{"message": "invalid app_id"}`                                                |
| 400    | `client_key` != `external_key`      | `{"message": "invalid client_key"}`                                            |
| 500    | Database upsert failure             | `{"message": "internal server error"}`                                         |

**Controller**: `InstallController.Callback`

---

### GET /webhooks/appmax

Webhook setup guide. Returns an HTML page with instructions for configuring webhooks
in the Appmax admin panel.

**Controller**: `WebhookController.Guide`

---

### POST /webhooks/appmax

Webhook event handler. Receives events from Appmax, persists them, and updates order statuses.

**Request Body** (JSON):

| Field       | JSON Key     | Type            | Description                              |
|-------------|--------------|-----------------|------------------------------------------|
| Event       | `event`      | string          | Event name (e.g., `OrderPaid`, `order_paid`) |
| EventType   | `event_type` | string          | Event category (`order`, `customer`, `subscription`) |
| Data        | `data`       | object (varies) | Payload — structure depends on webhook model |

**Success Responses**:

| Status | Body                            | Condition                        |
|--------|---------------------------------|----------------------------------|
| 200    | `{"message": "ok"}`             | Event processed successfully     |
| 200    | `{"message": "already processed"}` | Duplicate event detected      |

**Error Responses**:

| Status | Condition             | Body                                   |
|--------|-----------------------|----------------------------------------|
| 400    | Invalid request body  | `{"message": "invalid request body"}`  |
| 500    | Processing failure    | `{"message": "internal server error"}` |

For webhook payload models and order ID extraction logic, see
[webhooks.md](../appmax/webhooks.md) and the Integration Guide.

**Controller**: `WebhookController.Handle`

---

## Protected Endpoints

All endpoints below require a valid `{key}` route parameter (the `external_key` of an installation).
The `MerchantContext` middleware (`app/http/middleware/merchant_context.go`) validates the key
by looking up the installation in the database.

**Middleware error**: `404 {"message": "installation not found"}` if `{key}` does not match any installation.

---

### GET /installations/{key}/merchant/token

Fetch and return the merchant OAuth token for the installation. Useful for external systems
that need to call Appmax directly.

**Response** (200):
```json
{
  "merchant_bearer_token": "eyJhbGciOiJS...",
  "external_key": "merchant-key",
  "merchant_client_id": "client-id",
  "merchant_client_secret": "client-secret"
}
```

**Controller**: `MerchantAuthController.SyncToken`

---

### POST /checkout/{key}/order

Create a customer and order in Appmax.

**Request Body** (JSON):
```json
{
  "customer": {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "phone": "11999999999",
    "document_number": "52998224725",
    "ip": "192.168.1.1",
    "address": {
      "postcode": "01001000",
      "street": "Rua Example",
      "number": "123",
      "complement": "Apt 1",
      "district": "Centro",
      "city": "Sao Paulo",
      "state": "SP"
    }
  },
  "order": {
    "products_value": 10000,
    "discount_value": 0,
    "shipping_value": 500,
    "products": [
      {"sku": "PROD-001", "name": "Product A", "quantity": 1, "unit_value": 10000, "type": "physical"}
    ]
  }
}
```

**Response** (200):
```json
{
  "customer_id": 12345,
  "order_id": 67890
}
```

**Notes**:

- `document_number` must be a valid CPF (11-digit Brazilian tax ID). Sequential or placeholder values (e.g., `12345678900`) are rejected by Appmax. Use a CPF that passes the verification algorithm (e.g., `52998224725`).

**Controller**: `CheckoutController.CreateOrder`

---

### POST /checkout/{key}/pay/credit-card

Process a credit card payment. If `customer_id` and `order_id` are omitted, creates customer
and order from the `customer` and `order` fields first.

**Request Body** (JSON):
```json
{
  "customer_id": 12345,
  "order_id": 67890,
  "customer": { "..." : "same as create order" },
  "order": { "..." : "same as create order" },
  "payment": {
    "token": "card-token-if-tokenized",
    "number": "4000000000000010",
    "cvv": "123",
    "expiration_month": "12",
    "expiration_year": "2030",
    "holder_document_number": "52998224725",
    "holder_name": "JOHN DOE",
    "installments": 1,
    "soft_descriptor": "MYSTORE",
    "upsell_hash": ""
  },
  "subscription": {
    "interval": "monthly",
    "interval_count": 1
  }
}
```

**Response** (200):
```json
{
  "order_id": 67890,
  "status": "aprovado",
  "upsell_hash": "hash-for-upsell"
}
```

| Status | Condition        | Body                                  |
|--------|------------------|---------------------------------------|
| 422    | Payment declined | `{"message": "payment declined"}`     |
| 502    | Appmax error     | `{"message": "payment processing failed"}` |

**Notes**:

- `customer_id` and `order_id` are optional; if omitted, `customer` and `order` are required.
- `subscription` is optional; include for recurring payments.
- Provide `token` OR raw card fields, not both.
- `holder_document_number` must be a valid CPF (11-digit Brazilian tax ID). Sequential or placeholder values (e.g., `12345678900`) are rejected by Appmax. Use a CPF that passes the verification algorithm (e.g., `52998224725`).

**Controller**: `CheckoutController.PayCreditCard`

---

### POST /checkout/{key}/pay/pix

Process a Pix payment. Same customer/order auto-creation behavior as credit card.

**Request Body** (JSON):
```json
{
  "customer_id": 12345,
  "order_id": 67890,
  "customer": { "..." : "same as create order" },
  "order": { "..." : "same as create order" },
  "document_number": "52998224725",
  "subscription": {
    "interval": "monthly",
    "interval_count": 1
  }
}
```

**Response** (200):
```json
{
  "order_id": 67890,
  "qr_code": "base64-encoded-qr-image",
  "emv": "copy-paste-pix-code"
}
```

**Notes**:

- `document_number` must be a valid CPF (11-digit Brazilian tax ID). Sequential or placeholder values (e.g., `12345678900`) are rejected by Appmax. Use a CPF that passes the verification algorithm (e.g., `52998224725`).

**Controller**: `CheckoutController.PayPix`

---

### POST /checkout/{key}/pay/boleto

Process a Boleto payment. Same customer/order auto-creation behavior as credit card.

**Request Body** (JSON):
```json
{
  "customer_id": 12345,
  "order_id": 67890,
  "customer": { "..." : "same as create order" },
  "order": { "..." : "same as create order" },
  "document_number": "52998224725"
}
```

**Response** (200):
```json
{
  "order_id": 67890,
  "pdf_url": "https://...",
  "digitavel": "12345.67890..."
}
```

**Notes**:

- `document_number` must be a valid CPF (11-digit Brazilian tax ID). Sequential or placeholder values (e.g., `12345678900`) are rejected by Appmax. Use a CPF that passes the verification algorithm (e.g., `52998224725`).

**Controller**: `CheckoutController.PayBoleto`

---

### GET /checkout/{key}/status/{order_id}

Get order status from the local database (not from Appmax directly).

**Path Parameters**:

| Parameter  | Type    | Description                          |
|------------|---------|--------------------------------------|
| `order_id` | integer | Appmax order ID (must be > 0)        |

**Response** (200):
```json
{"status": "aprovado"}
```

| Status | Condition        | Body                              |
|--------|------------------|-----------------------------------|
| 400    | Invalid order_id | `{"message": "invalid order_id"}` |
| 404    | Order not found  | `{"message": "order not found"}`  |

**Controller**: `CheckoutController.Status`

---

### GET /checkout/{key}/installments

Get installment options for a given total value.

**Query Parameters**:

| Parameter      | Type    | Required | Default | Description                       |
|----------------|---------|----------|---------|-----------------------------------|
| `total_value`  | integer | Yes      | -       | Order total in cents              |
| `installments` | integer | No       | 12      | Maximum number of installments    |

**Response** (200):
```json
[
  {"installments": 1, "value": 10000.00, "total_value": 10000.00},
  {"installments": 2, "value": 5050.00, "total_value": 10100.00},
  {"installments": 3, "value": 3400.00, "total_value": 10200.00}
]
```

**Controller**: `CheckoutController.Installments`

---

### POST /checkout/{key}/refund

Request a refund for an order.

**Request Body** (JSON):
```json
{
  "order_id": 67890,
  "type": "total",
  "value": 10000
}
```

**Response** (200):
```json
{"message": "Refund request accepted"}
```

**Notes**:
- `type`: `"total"` or `"partial"`
- `value`: amount in cents (relevant for partial refunds)
- `order_id` must be > 0

**Controller**: `CheckoutController.Refund`

---

### POST /checkout/{key}/tokenize

Tokenize a credit card for future use.

**Request Body** (JSON):
```json
{
  "number": "4000000000000010",
  "cvv": "123",
  "expiration_month": "12",
  "expiration_year": "2030",
  "holder_name": "JOHN DOE"
}
```

**Response** (200):
```json
{"token": "card-token-string"}
```

**Controller**: `CheckoutController.Tokenize`

---

### POST /checkout/{key}/tracking

Add a shipping tracking code to an order.

**Request Body** (JSON):
```json
{
  "order_id": 67890,
  "shipping_tracking_code": "BR123456789"
}
```

**Response** (200):
```json
{"message": "tracking accepted"}
```

**Controller**: `CheckoutController.AddTracking`

---

### POST /checkout/{key}/upsell

Create an upsell offer for an existing order.

**Request Body** (JSON):
```json
{
  "upsell_hash": "hash-from-credit-card-response",
  "products_value": 5000,
  "products": [
    {"sku": "UPSELL-001", "name": "Bonus Product", "quantity": 1, "unit_value": 5000}
  ]
}
```

**Response** (200):
```json
{
  "message": "Upsell created",
  "redirect_url": "https://..."
}
```

**Controller**: `CheckoutController.Upsell`

---

## Error Response Reference

Consolidated reference for all error responses across endpoints. Every error response
follows the format `{"message": "..."}`.

### Common Errors (All Protected Endpoints)

| Status | Condition                        | Body                                        | Source           |
|--------|----------------------------------|---------------------------------------------|-----------------|
| 404    | `{key}` not found in DB         | `{"message": "installation not found"}`     | MerchantContext  |
| 400    | Invalid/missing request body     | `{"message": "invalid request body"}`       | Controller       |
| 500    | Unexpected internal failure      | `{"message": "internal server error"}`      | Controller       |

### Installation Errors

| Endpoint | Status | Condition | Body |
|----------|--------|-----------|------|
| `GET /install/start` | 400 | Missing `app_id` | `{"message": "app_id is required"}` |
| `GET /install/start` | 400 | Missing `external_key` | `{"message": "external_key is required"}` |
| `GET /install/start` | 502 | Appmax `/app/authorize` fails | `{"message": "failed to initiate installation"}` |
| `GET .../callback/install` | 200 | Cache miss (token invalid/expired/consumed) | `{"message": "installation confirmed"}` (graceful) |
| `GET .../callback/install` | 200 | Generate fails but installation exists in DB | `{"external_id": "uuid"}` (Path B completed) |
| `GET .../callback/install` | 400 | Missing `token` param | `{"message": "token is required"}` |
| `GET .../callback/install` | 400 | `app_id` in cached state != UUID | `{"message": "invalid app_id"}` |
| `GET .../callback/install` | 502 | `/app/client/generate` fails + no existing | `{"message": "failed to generate merchant credentials"}` |
| `POST .../callback/install` | 400 | Any of 5 fields empty | `{"message": "app_id, external_key, client_key, merchant_client_id and merchant_client_secret are required"}` |
| `POST .../callback/install` | 400 | `app_id` != `APPMAX_APP_ID_NUMERIC` | `{"message": "invalid app_id"}` |
| `POST .../callback/install` | 400 | `client_key` != `external_key` | `{"message": "invalid client_key"}` |
| `POST .../callback/install` | 500 | Database upsert failure | `{"message": "internal server error"}` |

### Checkout Errors

| Endpoint | Status | Condition | Body |
|----------|--------|-----------|------|
| `POST .../pay/credit-card` | 422 | Payment declined by Appmax | `{"message": "payment declined"}` |
| `POST .../pay/credit-card` | 502 | Appmax API error | `{"message": "payment processing failed"}` |
| `POST .../pay/pix` | 502 | Appmax API error | `{"message": "payment processing failed"}` |
| `POST .../pay/boleto` | 502 | Appmax API error | `{"message": "payment processing failed"}` |
| `GET .../status/{order_id}` | 400 | Invalid or zero `order_id` | `{"message": "invalid order_id"}` |
| `GET .../status/{order_id}` | 404 | Order not in local DB | `{"message": "order not found"}` |
| `POST .../refund` | 502 | Appmax refund fails | `{"message": "refund request failed"}` |

### Webhook Errors

| Endpoint | Status | Condition | Body |
|----------|--------|-----------|------|
| `POST /webhooks/appmax` | 400 | Invalid request body | `{"message": "invalid request body"}` |
| `POST /webhooks/appmax` | 500 | Processing failure | `{"message": "internal server error"}` |

### Upstream Error Propagation

When Appmax returns an error with a parseable message, our app extracts it and returns
it in the response body. The HTTP status code is mapped as follows:

| Appmax Status | Our Status | Description |
|---------------|------------|-------------|
| 400           | 502        | Bad request to Appmax (likely our bug) |
| 401/403       | 502        | Token expired or invalid credentials |
| 402           | 422        | Payment declined |
| 404           | 502        | Resource not found at Appmax |
| 422           | 422        | Validation error (declined, invalid data) |
| 500/502/503   | 502        | Appmax server error |

Implementation: `app/http/controllers/upstream_error_message.go`
