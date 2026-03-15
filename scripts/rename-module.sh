#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$ROOT_DIR"

OLD=$(grep '^module ' go.mod | awk '{print $2}')

NEW="${1:-}"
if [ -z "$NEW" ] && [ -f .env ]; then
  NEW=$(grep -E '^MODULE_PATH=' .env | cut -d= -f2 | tr -d '[:space:]')
fi

if [ -z "$NEW" ]; then
  echo "Usage: $0 <new-module-path>"
  echo "       or set MODULE_PATH in .env"
  exit 1
fi

if [ "$OLD" = "$NEW" ]; then
  echo "Module path already is '$NEW' — nothing to do."
  exit 0
fi

echo "Renaming module:"
echo "  old: $OLD"
echo "  new: $NEW"

go mod edit -module "$NEW"

find . -name '*.go' \
  -not -path './.git/*' \
  -not -path './.gocache/*' \
  -not -path './.gomodcache/*' \
  -not -path './vendor/*' \
  | xargs sed -i.bak "s|\"${OLD}|\"${NEW}|g"

find . -name '*.go.bak' -delete

echo "Done. Run 'go mod tidy' to verify."
