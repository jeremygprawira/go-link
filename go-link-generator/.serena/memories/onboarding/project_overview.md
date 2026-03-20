# Project Overview: go-link-generator

## Purpose
A production-ready Go REST API generator/link service built with Clean Architecture. It provides a standard structure for building scalable Go applications with features like tracing, metrics, and structured logging.

## Core Features
* **Framework**: Uses Echo (v4) for HTTP routing.
* **Database**: PostgreSQL with raw SQL queries (using pgx).
* **Migrations**: Goose for managing database schema changes.
* **Logging**: Structured logging using a custom Zap-based wrapper that supports "wide events" and contextual enrichment.
* **Tracing**: OpenTelemetry integration.
* **Metrics**: Prometheus support.
* **Config**: Viper for managing configuration (YAML-based).
* **Validation**: Go-playground/validator for input validation.
* **Documentation**: Swagger/OpenAPI documentation (Scalar UI).
* **Graceful Lifecycle**: Handles setup/teardown cleanly.
