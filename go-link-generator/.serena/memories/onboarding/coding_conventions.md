# Coding Conventions: go-link-generator

## General Patterns
* **Clean Architecture**: Strong separation between delivery (HTTP), service (Business Logic), and repository (Data Access).
* **Interface-Driven**: High-level layers depend on interfaces defined in lower layers (e.g., `UserService` interface, `UserRepository` interface).
* **Dependency Injection**: Dependencies are passed via structs (e.g., `service.Dependencies`).

## Naming & Style
* **Context**: Always pass `ctx context.Context` as the first argument to all methods in services, repositories, and logic packages.
* **Error Handling**: Use the `internal/pkg/errorc` package for standardizing API errors.
* **Logging**: 
    * Enrich context with `logger.Add(ctx, key, value)` or `logger.AddMap(ctx, map)`.
    * Log errors with `logger.AddError(ctx, &logger.ErrorContext{...})`.
    * Avoid direct use of `fmt.Println` or raw `zap`.

## Repository Implementation
* **SQL Queries**: Store SQL strings in variables within `internal/repository/pgsql/*_query.go` to keep implementation files clean.
* **Database Access**: Use `pgxpool.Pool` for executing queries.

## Service Implementation
* **Business Logic**: Services should orchestrate repository calls, validation, and transformations.
* **Enrichment**: Always add relevant business identifiers (like `user_id`, `account_number`) to the context logger for better observability.
