#!/usr/bin/env bash
set -euo pipefail

ENV_FILE="$(dirname "$0")/../.env"
if [ ! -f "$ENV_FILE" ]; then
  echo "ERROR: .env not found at $ENV_FILE"
  exit 1
fi

set -a
# shellcheck disable=SC1090
source "$ENV_FILE"
set +a

APP_PUBLIC_URL="${NGROK_URL:-${APP_URL:-}}"

for var in APP_ID_UUID APP_PUBLIC_URL; do
  if [ -z "${!var:-}" ]; then
    echo "ERROR: $var is not set (set NGROK_URL or APP_URL in .env)"
    exit 1
  fi
done

sep() { echo ""; echo "────────────────────────────────────────────────"; }
ts()  { date '+%H:%M:%S'; }

run=1
while true; do
  sep
  echo "[$(ts)] ITERACAO $run"
  sep

  EXTERNAL_KEY="test-curl-$(date +%s)"
  START_URL="${APP_PUBLIC_URL}/install/start?app_id=${APP_ID_UUID}&external_key=${EXTERNAL_KEY}"

  echo ""
  echo "[$(ts)] STEP 1 — /install/start (salva estado no Redis)"
  echo "  GET $START_URL"

  TMPHEADERS=$(mktemp)
  START_HTTP=$(curl -s \
    --write-out "%{http_code}" \
    --max-redirs 0 \
    --output /dev/null \
    -D "$TMPHEADERS" \
    "$START_URL" 2>/dev/null || true)

  LOCATION=$(grep -i "^location:" "$TMPHEADERS" | awk '{print $2}' | tr -d '\r' || true)
  rm -f "$TMPHEADERS"

  echo "  HTTP $START_HTTP"

  if [ "$START_HTTP" != "302" ]; then
    echo "  ERROR: esperado 302, recebido $START_HTTP"
    echo ""
    read -rp "Tentar novamente? [s/N] " again
    [[ "$again" =~ ^[sS]$ ]] || break
    ((run++))
    continue
  fi

  if [ -z "$LOCATION" ]; then
    echo "  ERROR: Location header não encontrado na resposta 302"
    echo ""
    read -rp "Tentar novamente? [s/N] " again
    [[ "$again" =~ ^[sS]$ ]] || break
    ((run++))
    continue
  fi

  HASH=$(echo "$LOCATION" | sed 's|.*/appstore/integration/||' | tr -d '[:space:]')
  echo "  Location: $LOCATION"
  echo "  Hash: $HASH"
  echo "  external_key: $EXTERNAL_KEY"
  echo "  Estado salvo no Redis: install:$HASH"

  sep
  echo ""
  echo "  ABRA NO BROWSER E AUTORIZE:"
  echo ""
  echo "  $LOCATION"
  echo ""
  echo "  Depois de autorizar, o Appmax vai redirecionar para:"
  echo "  ${APP_PUBLIC_URL}/integrations/appmax/callback/install?token=${HASH}"
  echo "  Nossa app vai chamar /app/client/generate automaticamente."
  echo ""
  read -rp "  Pressione ENTER depois de autorizar no browser..."

  echo ""
  echo "[$(ts)] STEP 2 — Aguardando callback processar..."
  sleep 3

  echo ""
  echo "[$(ts)] STEP 3 — Logs da app (últimas 25 linhas, excluindo health check):"
  echo ""
  docker compose logs app --tail=25 2>/dev/null | grep -v 'GET.*"/health"' | tail -25 || true

  echo ""
  CALLBACK_HTTP=$(docker compose logs app --tail=25 2>/dev/null \
    | grep "callback/install" \
    | grep -oP 'HTTP\] \K\d+' \
    | tail -1 || true)

  CALLBACK_OK=false
  if [[ "$CALLBACK_HTTP" == "200" ]]; then
    echo "[$(ts)] CALLBACK: OK (HTTP 200)"
    CALLBACK_OK=true
  elif [ -n "$CALLBACK_HTTP" ]; then
    echo "[$(ts)] CALLBACK: HTTP $CALLBACK_HTTP (ver logs acima)"
  else
    echo "[$(ts)] CALLBACK: sem log encontrado (pode ainda estar processando ou falhou antes)"
  fi

  sep
  echo ""
  echo "[$(ts)] STEP 4 — HEALTH CHECK (simula POST da breaking code para nossa app)"

  if $CALLBACK_OK; then
    MERCHANT_CLIENT_ID=$(docker compose logs app --tail=50 2>/dev/null \
      | grep -oP "merchant_client_id.*" | head -1 || echo "")
    HC_CLIENT_ID="real-from-appmax"
    HC_SECRET="real-from-appmax"
    echo "  Usando external_key real da instalação: $EXTERNAL_KEY"
  else
    HC_CLIENT_ID="fake-client-id-test"
    HC_SECRET="fake-client-secret-test"
    echo "  WARN: callback não completou — usando dados fictícios para testar endpoint isolado"
  fi

  HC_PAYLOAD="{\"app_id\":\"${APP_ID_UUID}\",\"client_id\":\"${HC_CLIENT_ID}\",\"client_secret\":\"${HC_SECRET}\",\"external_key\":\"hc-probe-$(date +%s)\"}"
  echo "  POST ${APP_PUBLIC_URL}/integrations/appmax/callback/install"
  echo "  Body: $HC_PAYLOAD"

  HC_TMPFILE=$(mktemp)
  HC_RESPONSE=$(curl -s \
    --write-out "\n__META__:%{http_code} %{time_total}s" \
    -X POST "${APP_PUBLIC_URL}/integrations/appmax/callback/install" \
    -H "Content-Type: application/json" \
    -d "$HC_PAYLOAD" \
    -D "$HC_TMPFILE" 2>/dev/null || true)

  HC_META=$(echo "$HC_RESPONSE" | grep "__META__:" || true)
  HC_JSON=$(echo "$HC_RESPONSE" | grep -v "__META__:" | tr -d '\n' || true)
  HC_HTTP=$(echo "$HC_META" | sed 's/.*__META__:\([^ ]*\).*/\1/')
  HC_TIME=$(echo "$HC_META" | sed 's/.*__META__:[^ ]* \(.*\)/\1/')
  rm -f "$HC_TMPFILE"

  echo "  HTTP $HC_HTTP | Time: $HC_TIME"
  echo "  Response: $HC_JSON"

  EXTERNAL_ID=$(echo "$HC_JSON" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('external_id',''))" 2>/dev/null || true)
  if [ -n "$EXTERNAL_ID" ]; then
    echo ""
    echo "  OK — external_id: $EXTERNAL_ID"
  fi

  sep
  echo ""
  echo "[$(ts)] RESUMO"
  echo "  1. /install/start:  HTTP $START_HTTP | Hash: ${HASH:0:16}..."
  echo "  2. Callback:        HTTP ${CALLBACK_HTTP:-?}"
  echo "  3. Health check:    HTTP $HC_HTTP | Time: $HC_TIME"
  echo ""

  read -rp "Repetir? [s/N] " again
  [[ "$again" =~ ^[sS]$ ]] || break
  ((run++))
done

echo ""
echo "Fim."
