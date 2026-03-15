# Testing Strategy

## Test suites
- Unit tests: `tests/unit/...`
- Integration tests: `tests/integration/...`
- Live E2E tests (tagged): `tests/end_to_end/appmax` with `//go:build appmax_live`

## Recommended commands
```bash
GOCACHE=$(pwd)/.gocache go test ./app/gateway/appmax
GOCACHE=$(pwd)/.gocache go test ./tests/unit/...
GOCACHE=$(pwd)/.gocache go test ./tests/integration/...
```

## Coverage commands
```bash
GOCACHE=$(pwd)/.gocache go test ./tests/unit/services -coverpkg=./app/... -coverprofile=tests_unit_services.cover
go tool cover -func=tests_unit_services.cover | tail -n 1

GOCACHE=$(pwd)/.gocache go test ./tests/unit/... -coverpkg=./app/... -coverprofile=tests_unit_all.cover
go tool cover -func=tests_unit_all.cover | tail -n 1
```

## Live E2E
```bash
GOCACHE=$(pwd)/.gocache go test -tags=appmax_live ./tests/end_to_end/appmax -v
```
Requires Appmax sandbox credentials/env vars.
