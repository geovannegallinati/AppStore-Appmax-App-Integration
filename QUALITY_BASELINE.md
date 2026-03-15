# Quality Baseline (Updated: 2026-03-14)

## Scope
- Repository root: `/Users/geovanne.gallinati/AppStoreAppDemo`
- Main packages: `app`, `bootstrap`, `routes`, `tests`

## Current quality snapshot
- `GOCACHE=$(pwd)/.gocache go test ./app/... ./bootstrap/... ./routes/... ./tests/unit/... ./tests/integration/...`
  - Status: **PASS**
- `GOCACHE=$(pwd)/.gocache go test ./tests/unit/services -coverpkg=./app/... -coverprofile=tests_unit_services.cover`
  - Status: **PASS**
  - Coverage (`go tool cover -func ... | tail -n1`): **57.3%**
- `GOCACHE=$(pwd)/.gocache go test ./tests/unit/... -coverpkg=./app/... -coverprofile=tests_unit_all.cover`
  - Status: **PASS**
  - Coverage (`go tool cover -func ... | tail -n1`): **49.2%**
- `GOCACHE=$(pwd)/.gocache go test ./app/gateway/appmax -coverprofile=cover_gateway.out`
  - Status: **PASS**
  - Coverage (`go tool cover -func ... | tail -n1`): **86.0%**

## Decoupling snapshot
- Removed direct framework coupling from services (`facades` no longer imported by `app/services`).
- Services no longer import the concrete gateway package; they depend only on `app/gateway/appmax/contracts`.
- Added adapters for framework logging/cache:
  - `app/adapters/goravel/logger.go`
  - `app/adapters/goravel/cache.go`
- `routes/api.go` no longer panics on dependency bootstrap failure; now logs and returns.

## Contracts snapshot
- Appmax gateway contract remains explicit in `app/gateway/appmax/contracts/gateway.go`.
- Gateway DTOs are defined in `app/gateway/appmax/contracts/types.go`.
- Added contract tests for gateway client in `app/gateway/appmax/contracts_test.go`.
- Service DTOs are service-owned (removed type aliases directly tied to gateway DTOs).
- Compile-time contract assertions added across implementations (`var _ Interface = (*Impl)(nil)`).

## Stages 1-4 status
- Stage 1 (baseline/hygiene): **completed**
- Stage 2 (decoupling): **completed**
- Stage 3 (contracts): **completed**
- Stage 4 (constructor/DI audit): **completed**

## Environment/IDE notes
- `golangci-lint` binary is not installed in this environment; `.golangci.yml` was added for CI/dev consistency.
- Live E2E tests use `appmax_live` build tag and external credentials.

## Remaining gap to strict 100%
- Full 100% coverage across all `app/...` is not achieved yet (current 43.6% in unit aggregate).
- Next increments should focus on controller/service integration flow and uncovered bootstrap paths.
