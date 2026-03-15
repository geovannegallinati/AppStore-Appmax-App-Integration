# Contracts

## Appmax gateway contract
- Interface: `app/gateway/appmax/contracts/gateway.go`
- Implementation: `app/gateway/appmax/*.go`
- Contract tests: `app/gateway/appmax/contracts_test.go`

## Service contracts
- Main interfaces live in `app/services/*.go`.
- DTOs in services are now service-owned (no direct type alias to gateway DTOs).
- Mapping from service DTO -> gateway DTO is explicit inside `app/services/appmax_service.go`.

## Repository contracts
- Interfaces in `app/repositories/contracts/*.go`
- Implementations in `app/repositories/*.go`

## Backward compatibility
- Constructors now prioritize dependency injection over URL/config-based object creation in services.
- Dependency-aware token manager constructors:
  - `NewTokenManagerWithGatewayDeps(...)`
  - `NewTokenManagerWithGatewayDepsAndClock(...)`
- Appmax service constructor is contract-first:
  - `NewAppmaxServiceWithGateway(...)`
