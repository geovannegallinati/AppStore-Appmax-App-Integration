# Quality Gates

## Mandatory checks
```bash
GOCACHE=$(pwd)/.gocache go test ./app/... ./bootstrap/... ./routes/... ./tests/unit/... ./tests/integration/...
GOCACHE=$(pwd)/.gocache go test ./tests/unit/services -coverpkg=./app/... -coverprofile=tests_unit_services.cover
go tool cover -func=tests_unit_services.cover | tail -n 1
```

## Optional local checks
```bash
go vet ./...
golangci-lint run ./...
```

## Gate policy
- Pull request cannot merge with failing unit/integration tests.
- Coverage must not decrease for `tests/unit/services -coverpkg=./app/...` baseline.
- New services/repositories/controllers must include constructor tests for nil/invalid dependencies.
