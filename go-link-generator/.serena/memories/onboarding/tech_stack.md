# Tech Stack: go-link-generator

## Languages & Core
* **Language**: Go 1.23.4
* **HTTP Framework**: [Echo v4](https://github.com/labstack/echo)

## Persistence
* **Database**: PostgreSQL
* **Driver**: `pgx/v5` (pool-based)
* **Migrations**: `goose`
* **Query Style**: Raw SQL (stored in dedicated `*_query.go` files)

## Observability
* **Logging**: `zap` (wrapped in `internal/pkg/logger`)
* **Tracing**: `opentelemetry` (OTEL SDK)
* **Metrics**: `prometheus`

## Utilities
* **Configuration**: `viper`
* **Validation**: `validator/v10`
* **Auth/Hashing**: `bcrypt` (for passwords), `jwt` (for tokens)
* **API Docs**: `swaggo` (Swagger)
