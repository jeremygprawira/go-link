# Suggested Commands: go-link-generator

## Development
* `make run`: Starts the HTTP server locally (uses `go run cmd/http/main.go`).
* `make tidy`: Runs `go mod tidy` to clean up dependencies.
* `make swagger`: Regenerates the Swagger documentation from code annotations.

## Infrastructure & Database
* `make docker-up`: Starts PostgreSQL and other dependencies in the background.
* `make docker-down`: Stops all Docker containers.
* `make migrate-up`: Applies pending SQL migrations (requires `DATABASE_URL`).
* `make migrate-down`: Reverts the last SQL migration.

## Quality Control
* `make test`: Runs all tests with coverage and verbose output.
* `make lint`: Executes `golangci-lint` to check for code smells and errors.
* `make build`: Compiles the application into a binary in `./builds/bin/`.
