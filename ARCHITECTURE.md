# Architecture

## Layers
- Transport/Framework adapters: `app/http`, `app/adapters/goravel`, `routes`, `bootstrap`
- Application services: `app/services`
- Persistence adapters: `app/repositories`
- External gateway adapter: `app/gateway/appmax`
- Domain models: `app/models`

## Dependency direction
- `services` depend on interfaces/contracts, never directly on Goravel facades.
- Goravel integration (cache/log/orm) is wired in `bootstrap` through adapters.
- Appmax HTTP details stay in `app/gateway/appmax`.

## Current composition root
- `bootstrap/http_dependencies.go`
- `bootstrap/service_module.go`
- `bootstrap/repository_module.go`
- `bootstrap/controller_module.go`

## Design rules
- No framework imports in `app/services`.
- Constructors validate required dependencies/config and return `error`.
- Controllers map request/response DTOs and delegate business rules to services.
