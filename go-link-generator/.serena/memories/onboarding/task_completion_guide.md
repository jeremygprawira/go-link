# Task Completion Guide: go-link-generator

Before considering a task finished:

1. **Linting**: Run `make lint` and fix any issues.
2. **Formatting**: Ensure code matches project style (standard `go fmt` is usually sufficient).
3. **Testing**: Run `make test` to ensure no regressions and verify new functionality.
4. **Documentation**:
    * If new API endpoints were added, ensure Swagger annotations are updated and run `make swagger`.
    * Update `README.md` if significant architectural changes or new setup steps were introduced.
5. **Dependencies**: Run `make tidy` if any new packages were added or removed.
6. **Logging**: Verify that new logic includes appropriate context enrichment via the `logger` package.
