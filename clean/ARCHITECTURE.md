# Architecture

Go service template based on Clean Architecture. Designed to be copied and adapted for any service.

The template includes a **Notification Service** as a working example that demonstrates all architectural patterns: REST + gRPC API, Kafka consumer, periodic worker, external API integration (Telegram).

## Layers

```
api  →  contracts  ←  usecase  →  contracts  ←  infra
                        ↓
                      domain
```

### domain

Pure business logic. Entities, value objects, domain events, domain errors. Minimal imports: only stdlib, `uuid`, and `pkg/optional`. No infra, no frameworks.

Flat package: `internal/domain/`. One file per entity. Entity, its value objects, and its events all live in the same file.

*Example: `domain/notification.go` contains `Notification`, `NotificationStatus`, `NotificationChannel`, `NotificationCreatedEvent`.*

### contracts

Interfaces and DTOs that define boundaries between layers. Three sub-packages:

- **`contracts/infra/`** — interfaces for infrastructure: repositories, message producers, external API clients. One file per infra client. If you have `infra/postgres/`, there's a corresponding `contracts/infra/postgres.go`.
- **`contracts/service/`** — interfaces for shared services (reusable business logic used by multiple use cases).
- **`contracts/usecase/`** — interfaces for use cases + their input DTOs. Handlers import these, never concrete implementations.

### usecase

One use case = one file = one struct with `Execute` method. Named `<action>_<entity>.go`.

Depends on interfaces from `contracts/infra` and `contracts/service`. Never on concrete infra.

*Example: `usecase/create_notification.go` depends on `infra.NotificationRepository` and `infra.EventProducer` interfaces.*

### service

Shared business logic reused across multiple use cases. If you see the same code in two use cases — extract a service.

*Example: `service/notification_sender.go` — send notification + update status. Used by both `SendNotification` (consumer) and `RetryFailed` (worker) use cases.*

### infra

Concrete implementations of `contracts/infra` interfaces. One sub-package per external system:

- `infra/postgres/` — repositories (sqlc + pgx), connection, type converters
- `infra/kafka/` — producer, consumer, OpenTelemetry trace propagation through headers
- `infra/telegram/` — Telegram Bot API sender

### api

Transport/delivery layer. Translates external input into use case calls. Four sub-packages by transport type:

- **`api/rest/`** — gin HTTP handlers. Each endpoint in its own folder: `endpoint.go` + `schemas.go`.
- **`api/grpc/`** — gRPC handlers. Proto definitions in `proto/`, generated code in `gen/go/`.
- **`api/consumer/`** — Kafka event handlers. Each handler in its own folder.
- **`api/worker/`** — asynq job handlers. Each job in its own folder.

All handlers accept **interfaces** from `contracts/usecase`, not concrete types.

### app

Uber fx dependency injection modules. One file per concern (`config.go`, `postgres.go`, `kafka.go`, etc.). Binds concrete implementations to interfaces via `fx.As`. Manages lifecycle (connection pools, graceful shutdown).

## Directory structure

```
clean/
├── cmd/                              # Entry points (thin — just fx.New + modules)
│   ├── api/main.go                   # REST + gRPC server
│   ├── consumer/main.go              # Kafka consumer
│   └── worker/main.go                # Asynq periodic worker
├── config/                           # YAML config files (layered)
├── gen/go/                           # Generated protobuf code (buf)
├── internal/
│   ├── api/                          # Transport layer
│   │   ├── rest/<endpoint>/          # endpoint.go + schemas.go per route
│   │   ├── grpc/                     # gRPC handlers
│   │   ├── consumer/<handler>/       # Kafka event handlers
│   │   └── worker/<job>/             # Asynq job handlers
│   ├── app/                          # fx DI modules
│   ├── config/                       # Config struct + viper loader
│   ├── contracts/                    # Interfaces + DTOs (boundaries)
│   │   ├── infra/                    # Repo, producer, sender interfaces
│   │   ├── service/                  # Shared service interfaces
│   │   └── usecase/                  # Use case interfaces + input DTOs
│   ├── domain/                       # Entities, value objects, events
│   ├── infra/                        # Infra implementations
│   │   ├── postgres/                 # sqlc queries, repos, converters
│   │   ├── kafka/                    # Producer, consumer, propagation
│   │   └── <client>/                 # Any external system
│   ├── service/                      # Shared business logic
│   └── usecase/                      # Use case implementations
├── migrations/                       # SQL migrations (goose)
├── pkg/                              # Reusable packages (not service-specific)
│   ├── logger/                       # Zap init + trace context enrichment
│   ├── optional/                     # Generic Optional[T]
│   └── otel/                         # OpenTelemetry provider init
└── proto/                            # Protobuf source files
```

## Config

Layered loading via viper: `common.yaml` -> `{env}.yaml` -> `sensitive.yaml` -> env vars. Each layer overrides only what it sets. `sensitive.yaml` is gitignored for local secrets.

## Code generation

Two codegen pipelines, both idempotent:

- **sqlc** — SQL queries in `infra/postgres/query/*.sql` -> Go code in `infra/postgres/sqlc/`. Schema derived from `migrations/`.
- **buf** — proto files in `proto/` -> Go code in `gen/go/`.

Run `task generate` to execute both.

## DI wiring

All dependency injection is in `internal/app/`. Each file is an `fx.Module`. Concrete types are bound to interfaces:

```go
fx.Provide(fx.Annotate(postgres.NewNotificationRepository, fx.As(new(infra.NotificationRepository))))
fx.Provide(fx.Annotate(usecase.NewCreateNotification, fx.As(new(uc.CreateNotification))))
```

Entry points in `cmd/` compose modules and invoke startup functions. Lifecycle hooks handle graceful shutdown.

## Conventions

- **Handlers depend on interfaces.** Import from `contracts/usecase`, never from `usecase`.
- **One use case per file.** `<action>_<entity>.go` with an `Execute` method.
- **One REST endpoint per folder.** `endpoint.go` (handler) + `schemas.go` (request/response).
- **Domain is pure.** No imports from infra, api, or external packages.
- **Domain validates.** Constructors return `(Entity, error)`. Invalid input yields `*ValidationError`.
- **Error handling is centralized.** REST middleware and gRPC interceptor map domain errors to status codes. Handlers do NOT map errors themselves.
- **Input DTOs live in contracts.** Shared between handlers and use case implementations.
- **Infra mapping stays in infra.** Domain <-> DB type conversion in `infra/postgres/convert.go`. pgx errors wrapped to domain errors in `infra/postgres/errors.go`.
- **Generated code is never edited.** Excluded from linting.
- **New infra client = new file in `contracts/infra/` + new folder in `infra/` + new fx module in `app/`.**

## Error model

Domain errors (`domain/errors.go`):
- `ErrNotFound` — entity not found (pgx.ErrNoRows)
- `ErrAlreadyExists` — unique constraint violated
- `ErrInvalidInput` — generic invalid input
- `ValidationError{Fields}` — per-field validation, `Is(ErrInvalidInput)` returns true

Transport mapping:
| Domain error | REST status | gRPC code |
|---|---|---|
| `ErrNotFound` | 404 | NotFound |
| `ErrAlreadyExists` | 409 | AlreadyExists |
| `ErrInvalidInput` / `ValidationError` | 422 | InvalidArgument |
| anything else | 500 | Internal |

Internal errors are logged but NOT leaked to the client.

## Testing

Reference test files — copy these when adding a new domain:

| Layer | Reference file | What it tests |
|---|---|---|
| Domain | `internal/domain/user_test.go` | Constructor valid/invalid, aggregate methods |
| Use case | `internal/usecase/create_user_test.go` | Mock repo, happy path / repo error / validation error |
| Transport | `internal/api/rest/create_user/endpoint_test.go` | httptest + mock use case, 201/422/500 |

## Adding to the template

### New entity

1. `domain/<entity>.go` — struct, validating constructor returning `(Entity, error)`, value objects, events
2. `migrations/` — goose SQL migration
3. `infra/postgres/query/<entity>.sql` — sqlc queries
4. `task generate:sqlc`
5. `infra/postgres/<entity>.go` — repo wrapping sqlc, domain mapping, error wrapping via `mapError()`
6. `contracts/infra/postgres.go` — add repository interface
7. `app/postgres.go` — register with `fx.As`

### New use case

1. `contracts/usecase/<entity>.go` — interface + input DTO
2. `usecase/<action>_<entity>.go` — implementation
3. `app/usecase.go` — register with `fx.As`

### New REST endpoint

1. `api/rest/<action>_<entity>/schemas.go` — request/response
2. `api/rest/<action>_<entity>/endpoint.go` — `func New(uc Interface) gin.HandlerFunc`. On error call `middleware.SetError(c, err)`, do NOT write JSON error directly.
3. `api/rest/router.go` — register route

### New gRPC method

1. `proto/<entity>.proto` — add rpc + messages
2. `task generate:proto`
3. `api/grpc/<entity>.go` — handler method

### New Kafka handler

1. `domain/` — event type + struct (if new event)
2. `api/consumer/<handler>/handler.go` — parse event, call use case
3. `cmd/consumer/main.go` — `consumer.Register(eventType, handler.New(uc))`

### New worker job

1. `api/worker/<job>/handler.go` — `TaskType`, `NewTask()`, `New()` handler
2. `cmd/worker/main.go` — register scheduler + handler

### New infra client

1. `contracts/infra/<client>.go` — interface + DTOs
2. `infra/<client>/` — implementation
3. `app/<client>.go` — fx module with lifecycle hooks
4. `cmd/*/main.go` — add module

### New shared service

1. `contracts/service/<name>.go` — interface
2. `service/<name>.go` — implementation
3. `app/service.go` — register with `fx.As`

## Tech stack

| Concern | Tool |
|---|---|
| HTTP | gin |
| gRPC | google.golang.org/grpc + buf |
| Database | PostgreSQL, pgx, sqlc |
| Migrations | goose |
| Messaging | Kafka (franz-go) |
| Jobs | asynq (Redis) |
| Config | viper (layered YAML + env) |
| Logging | zap |
| Telemetry | OpenTelemetry (OTLP) |
| DI | uber fx |
| Linting | golangci-lint, buf lint |
| Tasks | Taskfile |

## Commands

```sh
task infra:up              # Start docker infra
task infra:down            # Stop
task infra:destroy         # Stop + delete volumes
task migrate:up            # Apply migrations
task migrate:down          # Rollback last
task generate              # Run all codegen
task generate:sqlc         # sqlc only
task generate:proto        # buf only
task run:api               # Start API
task run:consumer          # Start consumer
task run:worker            # Start worker
task lint                  # Run all linters
task lint:go               # Go linters only
task lint:proto            # Proto linters only
go test ./... -race        # Run all tests
```
