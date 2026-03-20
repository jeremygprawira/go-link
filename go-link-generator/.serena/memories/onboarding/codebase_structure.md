# Codebase Structure: go-link-generator

```
.
├── cmd/
│   └── http/               # Entry point (main.go)
├── builds/                 # Binary outputs and Dockerfiles
├── config/                 # YAML configuration files (.local.yaml, .example.yaml)
├── docs/                   # Generated Swagger documentation
├── internal/
│   ├── config/             # Configuration loading logic (Viper)
│   ├── core/               # Application setup and dependency injection (Setup/Teardown)
│   ├── deliveries/http/    # HTTP handlers, router definition, and middleware
│   ├── models/             # API DTOs and Domain Models
│   ├── pkg/                # Shared internal packages (database, logger, validator, etc.)
│   ├── repository/         # Data access layer (PostgreSQL implementation)
│   └── service/            # Business logic layer
└── migration/              # SQL migration scripts
```
