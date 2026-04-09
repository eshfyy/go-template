# CLAUDE.md

@AGENTS.md

Read ARCHITECTURE.md for full project structure and AGENTS.md for AI-agent-specific instructions.

## Quick reference

- This is a Go service template using Clean Architecture
- The Notification Service is an example domain — treat it as a reference, not the goal
- Entry points: `cmd/api/` (REST+gRPC), `cmd/consumer/` (Kafka), `cmd/worker/` (asynq)

## Build & run

```sh
task infra:up          # docker compose (postgres, kafka, redis, jaeger)
task migrate:up        # apply migrations
task generate          # codegen (sqlc + buf)
task run:api           # start API locally
task run:consumer      # start consumer locally
task run:worker        # start worker locally
task lint              # golangci-lint + buf lint
go build ./...         # verify compilation
go test ./... -race    # run all tests
```

## Key rules

- **Handlers depend on interfaces** from `contracts/usecase/`, never on concrete `usecase.*` types — this applies to ALL transports: REST, gRPC, consumer, worker
- **Domain has minimal imports** — only stdlib, `uuid`, and `pkg/optional`. No infra, no frameworks.
- **Domain validates** — constructors return `(Entity, error)` with `ValidationError` for invalid input
- **Never edit generated code** — `gen/go/` (buf) and `internal/infra/postgres/sqlc/` (sqlc)
- **One use case per file** — `usecase/<action>_<entity>.go` with `Execute` method
- **One REST endpoint per folder** — `api/rest/<name>/endpoint.go` + `schemas.go`
- **Input DTOs live in contracts** — `contracts/usecase/`, not in `usecase/`
- **Error handling is centralized** — REST uses `middleware.SetError(c, err)` for ALL errors (including bind/parse), gRPC uses error interceptor. Handlers do NOT write error JSON directly.
- **New infra = new contract + new impl + new fx module** — `contracts/infra/<x>.go`, `infra/<x>/`, `app/<x>.go`
- **DI is in `internal/app/`** — bind concrete to interface with `fx.As`, add lifecycle hooks for cleanup
- **Config layers**: `common.yaml` → `{env}.yaml` → `sensitive.yaml` → env vars

## Error taxonomy

- `domain.ErrNotFound` → REST 404 / gRPC NotFound
- `domain.ErrAlreadyExists` → REST 409 / gRPC AlreadyExists
- `domain.ErrInvalidInput` / `domain.ValidationError` → REST 422 / gRPC InvalidArgument
- Everything else → REST 500 / gRPC Internal (error message NOT leaked to client)

## Testing

Reference tests to copy:
- Domain: `internal/domain/user_test.go`
- Use case: `internal/usecase/create_user_test.go`
- Transport: `internal/api/rest/create_user/endpoint_test.go`

## After making changes

1. `go build ./...` — must compile
2. `go test ./... -race` — must pass
3. `task lint` — must pass
4. If you changed `.sql` queries or migrations — run `task generate:sqlc`
5. If you changed `.proto` files — run `task generate:proto`
6. If you added a new dependency to DI — register it in `internal/app/` and add `fx.As` for interface binding
