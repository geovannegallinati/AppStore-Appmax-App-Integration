# Appmax API Reference

Upstream API endpoints provided by Appmax. These are **not** our app's endpoints — they are
the external services our integration calls.

Source: Appmax official docs (appmax.readme.io), cross-referenced with our gateway
implementation in `app/gateway/appmax/`.

---

## Environments

| Environment | Auth Base URL                           | API Base URL                           | Admin / Redirect Base URL                       |
|-------------|----------------------------------------|----------------------------------------|-------------------------------------------------|
| Sandbox     | `https://auth.sandboxappmax.com.br`    | `https://api.sandboxappmax.com.br`     | `https://breakingcode.sandboxappmax.com.br`     |
| Production  | `https://auth.appmax.com.br`           | `https://api.appmax.com.br`            | `https://admin.appmax.com.br`                   |

Configured via env vars `APPMAX_AUTH_URL`, `APPMAX_API_URL`, `APPMAX_ADMIN_URL`.
Defaults to production if not set.

---

## Authentication

### POST {AUTH_URL}/oauth2/token

Obtain an OAuth2 access token using client credentials grant.

**Content-Type**: `application/x-www-form-urlencoded`

| Parameter       | Type   | Description                                      |
|-----------------|--------|--------------------------------------------------|
| grant_type      | string | Always `client_credentials`                      |
| client_id       | string | App-level or merchant-level client ID             |
| client_secret   | string | Corresponding client secret                       |

**Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJS...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

**Notes**:
- No refresh token mechanism. When expired, request a new token.
- `expires_in` is in seconds (typically 3600 = 1 hour).
- Two credential pairs exist (see Integration Guide for details):
  - **App credentials** (`APPMAX_CLIENT_ID` / `APPMAX_CLIENT_SECRET`): limited scope, installation only.
  - **Merchant credentials** (per-installation `client_id` / `client_secret`): full API scope.

---

## App Installation

### POST {API_URL}/app/authorize

Generate an installation hash for merchant authorization.

**Auth**: Bearer token (app-level).

**Body** (JSON):
```json
{
  "app_id": "uuid-of-the-app",
  "external_key": "merchant-identifier",
  "url_callback": "https://your-app.com/integrations/appmax/callback/install"
}
```

**Response** (200/201):
```json
{
  "data": {
    "token": "installation-hash-string"
  }
}
```

The returned `token` (hash) is used to build the redirect URL:
`{ADMIN_URL}/appstore/integration/{hash}`

---

### POST {API_URL}/app/client/generate

Exchange an installation hash for merchant credentials.

**Auth**: Bearer token (app-level).

**Body** (JSON):
```json
{
  "token": "installation-hash-string"
}
```

**Response** (200/201):
```json
{
  "data": {
    "client": {
      "client_id": "merchant-client-id",
      "client_secret": "merchant-client-secret"
    }
  }
}
```

---

## Customer Management

### POST {API_URL}/v1/customers

Create or update a customer. Returns the Appmax customer ID.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
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
  },
  "products": [
    {"sku": "PROD-001", "name": "Product A", "quantity": 1, "unit_value": 5000, "type": "physical"}
  ],
  "tracking": {
    "utm_source": "google",
    "utm_campaign": "summer-sale"
  }
}
```

**Response** (200/201):
```json
{
  "data": {
    "customer": {
      "id": 12345
    }
  }
}
```

**Notes**:
- `address`, `products`, and `tracking` are optional.
- If the customer already exists (matched by email), it is updated.
- The `data` field may be an empty array `[]` in some edge cases (our gateway handles this).
- `document_number` must be a valid CPF (11-digit Brazilian tax ID). Appmax validates the CPF checksum — sequential or placeholder values (e.g., `12345678900`) are rejected. Use a CPF that passes the verification algorithm (e.g., `52998224725`).

---

## Order Management

### POST {API_URL}/v1/orders

Create an order. Returns the Appmax order ID.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
  "customer_id": 12345,
  "products_value": 10000,
  "discount_value": 0,
  "shipping_value": 500,
  "products": [
    {"sku": "PROD-001", "name": "Product A", "quantity": 1, "unit_value": 10000, "type": "physical"}
  ]
}
```

**Response** (200/201):
```json
{
  "data": {
    "order": {
      "id": 67890,
      "status": "pendente"
    }
  }
}
```

**Notes**:
- Values are in cents (10000 = R$100.00).
- `data` may be an empty array `[]` in error scenarios.

---

### GET {API_URL}/v1/orders/{order_id}

Retrieve order details.

**Auth**: Bearer token (merchant-level).

**Response** (200):
```json
{
  "data": {
    "order": {
      "id": 67890,
      "status": "aprovado",
      "total_paid": 10500,
      "amounts": {
        "sub_total": 10000,
        "shipping_value": 500,
        "discount": 0,
        "installment_fee": 0
      },
      "created_at": "2026-03-15 10:00:00",
      "updated_at": "2026-03-15 10:05:00"
    },
    "customer": {
      "id": 12345,
      "name": "John Doe",
      "email": "john@example.com",
      "document_number": "52998224725"
    },
    "payment": {
      "method": "credit_card",
      "installments": 3,
      "installments_amount": 3500,
      "card": {
        "brand": "visa",
        "number": "****0010"
      },
      "paid_at": "2026-03-15 10:05:00"
    }
  }
}
```

---

## Payments

### POST {API_URL}/v1/payments/credit-card

Process a credit card payment.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
  "order_id": 67890,
  "customer_id": 12345,
  "payment_data": {
    "credit_card": {
      "token": "card-token-if-tokenized",
      "number": "4000000000000010",
      "cvv": "123",
      "expiration_month": "12",
      "expiration_year": "2030",
      "holder_document_number": "52998224725",
      "holder_name": "JOHN DOE",
      "installments": 1,
      "soft_descriptor": "MYSTORE"
    },
    "subscription": {
      "interval": "monthly",
      "interval_count": 1
    }
  }
}
```

**Response** (200/201):
```json
{
  "data": {
    "payment": {
      "id": 99999,
      "pay_reference": "ref-string",
      "upsell_hash": "hash-for-upsell",
      "status": "aprovado"
    }
  }
}
```

**Notes**:
- Provide either `token` (from tokenization) OR raw card fields (`number`, `cvv`, etc.), not both.
- `upsell_hash` is present when upsell is available.
- `subscription` is optional; include for recurring payments.
- Empty `pay_reference` indicates a declined payment.
- `data` may be an empty array in error scenarios.
- `holder_document_number` must be a valid CPF. See the customer creation note above.

---

### POST {API_URL}/v1/payments/pix

Process a Pix payment.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
  "order_id": 67890,
  "payment_data": {
    "pix": {
      "document_number": "52998224725"
    },
    "subscription": {
      "interval": "monthly",
      "interval_count": 1
    }
  }
}
```

**Response** (200/201):
```json
{
  "data": {
    "payment": {
      "pix_qrcode": "base64-encoded-qr-image",
      "pix_emv": "copy-paste-pix-code"
    }
  }
}
```

**Notes**:
- `subscription` is optional; include for recurring Pix payments.
- Empty `pix_qrcode` may indicate the order already has a pending payment.
- `document_number` must be a valid CPF. See the customer creation note above.

---

### POST {API_URL}/v1/payments/boleto

Process a Boleto payment.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
  "order_id": 67890,
  "payment_data": {
    "boleto": {
      "document_number": "52998224725"
    }
  }
}
```

**Response** (200/201):
```json
{
  "data": {
    "payment": {
      "boleto_link_pdf": "https://...",
      "boleto_digitable_line": "12345.67890..."
    }
  }
}
```

**Notes**:

- `document_number` must be a valid CPF. See the customer creation note above.

---

## Utilities

### POST {API_URL}/v1/payments/installments

Calculate installment options.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
  "installments": 12,
  "total_value": 10000,
  "settings": true
}
```

**Response** (200):
```json
{
  "data": {
    "parcels": {
      "1": 10000.00,
      "2": 5050.00,
      "3": 3400.00
    }
  }
}
```

**Notes**:
- `total_value` in cents.
- `parcels` keys are installment counts (as strings), values are per-installment amounts.
- `settings: true` applies merchant-specific installment rules.

---

### POST {API_URL}/v1/payments/tokenize

Tokenize a credit card for future use.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
  "payment_data": {
    "credit_card": {
      "number": "4000000000000010",
      "cvv": "123",
      "expiration_month": "12",
      "expiration_year": "2030",
      "holder_name": "JOHN DOE"
    }
  }
}
```

**Response** (200/201):
```json
{
  "data": {
    "token": "card-token-string"
  }
}
```

---

### POST {API_URL}/v1/orders/refund-request

Request a refund for an order.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
  "order_id": 67890,
  "type": "total",
  "value": 10000
}
```

**Notes**:
- `type`: `"total"` for full refund, `"partial"` for partial.
- `value`: amount in cents (required for partial, ignored for total).

---

### POST {API_URL}/v1/orders/shipping-tracking-code

Add a shipping tracking code to an order.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
  "order_id": 67890,
  "shipping_tracking_code": "BR123456789"
}
```

---

### POST {API_URL}/v1/orders/upsell

Create an upsell offer for an existing order.

**Auth**: Bearer token (merchant-level).

**Body** (JSON):
```json
{
  "upsell_hash": "hash-from-credit-card-response",
  "products_value": 5000,
  "products": [
    {"sku": "UPSELL-001", "name": "Bonus Product", "quantity": 1, "unit_value": 5000}
  ]
}
```

**Response** (200/201):
```json
{
  "data": {
    "message": "Upsell created",
    "redirect_url": "https://...",
    "order": {
      "id": 67891,
      "status": "aprovado"
    }
  }
}
```

---

## Order Statuses

| Status                               | Description                                                     |
|--------------------------------------|-----------------------------------------------------------------|
| `pendente`                           | Payment pending (awaiting Pix/Boleto confirmation)              |
| `aprovado`                           | Payment approved, funds available                               |
| `autorizado`                         | Credit card authorized, anti-fraud analysis in progress         |
| `cancelado`                          | Payment not authorized / declined / expired                     |
| `estornado`                          | Refunded                                                        |
| `integrado`                          | Final status, order ready to ship                               |
| `pendente_integracao`                | Approved but pending integration issues                         |
| `pendente_integracao_em_analise`     | Approved, refund pending manual analysis                        |
| `recusado_por_risco`                 | Declined for fraud risk                                         |
| `chargeback_em_tratativa`            | Chargeback under analysis                                       |
| `chargeback_em_disputa`              | Appmax disputing chargeback                                     |
| `chargeback_perdido`                 | Chargeback lost                                                 |
| `chargeback_vencido`                 | Chargeback expired / recovered                                  |

---

## Test Credit Cards (Sandbox Only)

| Card Number          | Behavior                              |
|----------------------|---------------------------------------|
| `4000000000000010`   | Success (all operations succeed)      |
| `4000000000000028`   | Failure (not authorized)              |
| Any other number     | Not authorized                        |

Use any future expiration date, any 3-digit CVV, and any holder name.
