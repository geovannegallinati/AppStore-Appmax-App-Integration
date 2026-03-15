#!/usr/bin/env bash
set -euo pipefail

ENV_FILE="$(dirname "$0")/../.env"
if [ ! -f "$ENV_FILE" ]; then
  echo "ERROR: .env not found at $ENV_FILE"
  exit 1
fi

getenv() {
  grep -E "^$1=" "$ENV_FILE" | head -1 | cut -d= -f2- | sed "s/^['\"]//;s/['\"]$//"
}

NGROK_URL=$(getenv NGROK_URL)
APP_URL=$(getenv APP_URL)
APP_ID_UUID=$(getenv APP_ID_UUID)
APPMAX_APP_ID_NUMERIC=$(getenv APPMAX_APP_ID_NUMERIC)

APP_PUBLIC_URL="${NGROK_URL:-${APP_URL:-}}"
APP_ID="${APPMAX_APP_ID_NUMERIC:-${APP_ID_UUID}}"

if [ -z "$APP_PUBLIC_URL" ] || [ -z "$APP_ID" ]; then
  echo "ERROR: APP_PUBLIC_URL or APP_ID not set in .env"
  exit 1
fi

ENDPOINT="${APP_PUBLIC_URL}/integrations/appmax/callback/install"
CLIENT_ID="${1:-merchant-client-id-sandbox-test}"
CLIENT_SECRET="${2:-merchant-client-secret-sandbox-test}"
EXTERNAL_KEY="${3:-simulated-store-$(date +%s)}"

echo ""
echo "Simulando POST da Appmax → health check de instalação"
echo "Endpoint: $ENDPOINT"
echo ""
echo "Payload:"
echo "  app_id:        $APP_ID"
echo "  external_key:  $EXTERNAL_KEY"
echo "  client_id:     $CLIENT_ID"
echo "  client_secret: $CLIENT_SECRET"
echo ""

PAYLOAD="{\"app_id\":\"${APP_ID}\",\"client_id\":\"${CLIENT_ID}\",\"client_secret\":\"${CLIENT_SECRET}\",\"external_key\":\"${EXTERNAL_KEY}\"}"

RESPONSE=$(curl -s \
  --write-out "\n__HTTP__:%{http_code} __TIME__:%{time_total}s" \
  -X POST "$ENDPOINT" \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD" 2>/dev/null)

HTTP_CODE=$(echo "$RESPONSE" | grep -o '__HTTP__:[^ ]*' | cut -d: -f2)
TIME=$(echo "$RESPONSE" | grep -o '__TIME__:[^ ]*' | cut -d: -f2)
BODY=$(echo "$RESPONSE" | grep -v "__HTTP__:")

echo "HTTP $HTTP_CODE | $TIME"
echo "Response: $BODY"
echo ""

if echo "$BODY" | python3 -c "import sys,json; d=json.load(sys.stdin); print('external_id:', d['external_id'])" 2>/dev/null; then
  echo "Instalação registrada com sucesso."
else
  echo "Falha. Ver logs: docker compose logs app --tail=10"
fi
