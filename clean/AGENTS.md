# AGENTS.md

Instructions for AI agents working on this codebase. Read ARCHITECTURE.md for full details.

## Layer structure

```
cmd/           Entry points (thin fx.New wiring)
internal/
  api/         Transport handlers (rest, grpc, consumer, worker)
  app/         Uber fx DI modules
  config/      Config struct + viper loader
  contracts/   Interfaces + DTOs (boundaries between layers)
    infra/     Repo, producer, sender interfaces
    service/   Shared service interfaces
    usecase/   Use case interfaces + input DTOs
  domain/      Entities, value objects, events, errors
  infra/       Infrastructure implementations (postgres, kafka, telegram)
  service/     Shared business logic
  usecase/     Use case implementations
gen/go/        Generated protobuf code (buf)
migrations/    SQL migrations (goose)
pkg/           Reusable packages (logger, optional, otel)
proto/         Protobuf source files
```

## Dependency rules

- `api/` imports `contracts/usecase` (interfaces), never `usecase/` (implementations)
- `usecase/` imports `contracts/infra` and `contracts/service` (interfaces), never `infra/` or `service/`
- `domain/` imports only stdlib, `github.com/google/uuid`, and `pkg/optional` — no infra, no frameworks
- `infra/` imports `domain/` and `contracts/infra` (to implement interfaces)
- DI wiring lives in `internal/app/`, binds concrete to interface via `fx.As`

## Error model

Domain errors are in `internal/domain/errors.go`:
- `ErrNotFound` — entity not found
- `ErrAlreadyExists` — unique constraint violated
- `ErrInvalidInput` — generic invalid input
- `ValidationError{Fields}` — per-field validation errors, `Is(ErrInvalidInput)` returns true

Transport mapping is centralized:
- REST: `internal/api/rest/middleware/error_handler.go` — maps domain errors to HTTP status codes (404/409/422/500)
- gRPC: `internal/api/grpc/interceptor/error.go` — maps domain errors to gRPC status codes

Handlers do NOT map errors themselves:
- REST handlers call `middleware.SetError(c, err)` for ALL errors including bind/parse failures
- gRPC handlers return raw domain errors; the interceptor maps them

Repositories wrap pgx errors into domain errors (`ErrNotFound`, `ErrAlreadyExists`) in `infra/postgres/errors.go`.

## Lifecycle patterns

- **Long-running processes** (consumer, worker) must create their own `context.WithCancel(context.Background())` — never use the fx OnStart context which cancels after startup.
- **HTTP server** uses `http.Server` with `Shutdown(ctx)` in OnStop for graceful drain.
- **gRPC server** uses `GracefulStop()` in OnStop.
- **Postgres pool** is created in the fx constructor, but `Ping` runs in OnStart with the fx startup-timeout context. `Close` runs in OnStop.

## Testing patterns

Reference tests exist — copy their structure:
- **Domain**: `internal/domain/user_test.go` — constructor valid/invalid, aggregate methods
- **Use case**: `internal/usecase/create_user_test.go` — mock repos via interface, test happy/error/validation paths
- **Transport**: `internal/api/rest/create_user/endpoint_test.go` — httptest + mock use case, test 201/422/500

## Checklist: new domain end-to-end

1. `internal/domain/<entity>.go` — entity + validating constructor + aggregate methods
2. `migrations/NNN_create_<table>.sql` — migration + indexes for query patterns
3. `internal/infra/postgres/query/<entity>.sql` — sqlc queries
4. `task generate:sqlc` — regenerate
5. `internal/contracts/infra/postgres.go` — repository interface
6. `internal/infra/postgres/<entity>.go` — repo implementation with hydrate function + error mapping
7. `internal/contracts/usecase/<entity>.go` — use case interfaces + input DTOs
8. `internal/usecase/<action>_<entity>.go` — use case implementations
9. `internal/app/postgres.go` — register repo with `fx.As`
10. `internal/app/usecase.go` — register use cases with `fx.As`
11. REST: `internal/api/rest/<action>_<entity>/endpoint.go` + `schemas.go`, register in `router.go`
12. gRPC: `proto/<entity>.proto` → `task generate:proto` → `internal/api/grpc/<entity>.go` → register in `server.go`
13. Consumer (if needed): `internal/api/consumer/<handler>/handler.go` → register in `cmd/consumer/main.go`
14. Worker (if needed): `internal/api/worker/<job>/handler.go` → register in `cmd/worker/main.go`
15. Tests: domain + usecase + at least one transport endpoint (follow user reference tests)

## Adding a new SQL query

1. Add query to `infra/postgres/query/<entity>.sql` with sqlc annotation
2. Run `task generate:sqlc`
3. Update repo wrapper in `infra/postgres/<entity>.go`
4. Update interface in `contracts/infra/postgres.go` if needed

## Adding a new gRPC method

1. Add rpc + messages to `proto/<entity>.proto`
2. Run `task generate:proto`
3. Add handler in `api/grpc/<entity>.go`

## Adding a new REST endpoint

1. Create `api/rest/<action>_<entity>/schemas.go` (request/response)
2. Create `api/rest/<action>_<entity>/endpoint.go` — `func New(uc Interface) gin.HandlerFunc`
3. On error: call `middleware.SetError(c, err)` — do NOT write JSON error or use `c.JSON(4xx, ...)` directly
4. Register in `api/rest/router.go`

## Commands

```sh
task generate          # sqlc + buf codegen
task lint              # golangci-lint + buf lint
go build ./...         # verify compilation
go test ./... -race    # run all tests
```
