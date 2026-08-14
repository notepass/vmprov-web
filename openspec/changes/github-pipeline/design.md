## Context

The project currently has no CI/CD automation. Tests are run manually via `make test`. Integration tests require PostgreSQL and MySQL to be running locally. There is no `Dockerfile` for containerization.

## Goals / Non-Goals

**Goals:**
- Provide a manual CI workflow that builds binaries and runs all tests
- Containerize database dependencies (PostgreSQL 18, MariaDB 12) for reproducible integration testing
- Add a Dockerfile for the application itself
- Cross-compile linux/amd64 and linux/arm64 binaries

**Non-Goals:**
- Linting or static analysis
- Docker image publishing to registries
- Automatic triggers on push/PR
- macOS or Windows builds

## Decisions

### GitHub Services vs. docker-compose
GitHub Actions `services` key was considered but `docker-compose.yml` is preferred because:
- The same compose file works locally for developers
- Easier to add more services later
- Explicit healthcheck configuration

### Integration test connection strings
The existing integration tests hardcode connection strings pointing to `localhost`. In GitHub Actions, services connect via `localhost` by default (GitHub maps service ports to the runner's localhost). No code changes needed.

### Makefile integrate target
Adding `make integrate` that runs `go test ./...` (without `-short`) provides a consistent local command for running integration tests, matching the CI behavior.

### Dockerfile multi-stage build
Use `golang:1.26` for building and `alpine` for the final image to minimize image size.

## Risks / Trade-offs

[Database startup timing] → Use docker-compose healthcheck and `depends_on: condition: service_healthy` to ensure databases are ready before tests run.

[Connection string mismatch] → MariaDB uses the MySQL driver (`github.com/go-sql-driver/mysql`) with the same DSN format, so existing tests work without changes.

[Manual-only workflow] → No automated gate on PRs. Developers must remember to trigger the workflow. Acceptable for now; can add automatic triggers later.

## Open Questions

- Should we add a `workflow_dispatch` input to select which database(s) to test?
