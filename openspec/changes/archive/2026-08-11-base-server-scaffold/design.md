## Context

The project starts from zero code. Go was selected as the language based on research (RESEARCH.md) for its superior libvirt binding, ISO creation ecosystem, and deployment model. This change establishes the foundational server framework and tooling before any domain-specific API endpoints are added.

## Goals / Non-Goals

**Goals:**
- Initialize Go module with `github.com/notepass/vmprov-web`
- Standard Go project layout (`cmd/server/`, `internal/`)
- Working HTTP server using Echo framework
- Configuration via `config.yaml` + environment variables (port, DB connect string, DB username, DB password)
- Structured logging with slog (stdlib)
- Testify as the test assertion framework for all backend components
- Makefile for `build`, `test`, `run` workflow

**Non-Goals:**
- API endpoints (future change)
- Libvirt integration (future change)
- Cloud-init ISO generation (future change)
- Frontend UI (future change)

## Decisions

### Echo for web framework
Chosen over Gin for its cleaner middleware design and more elegant routing API. Echo is lightweight with built-in middleware for common concerns.

### Standard Go project layout
`cmd/server/` for the entry point, `internal/` for private application packages. Follows Go community conventions.

### Viper for configuration
YAML config file + environment variable overrides via viper. All config fields have corresponding env vars for Kubernetes deployments.

### slog for logging
Standard library structured logging (Go 1.21+). Zero dependencies, supports log levels (DEBUG/INFO/WARN/ERROR), JSON output. Global log level configurable via `LOG_LEVEL` env var. Per-package level control is possible but not implemented in this baseline.

### Testify for tests
`testify/assert` and `testify/require` for clean assertions. Standard `testing` package as the test runner.

### Makefile for build workflow
Provides `make build`, `make run`, `make test`, `make lint` targets for consistent development.

## Risks / Trade-offs

- [Echo lock-in] → Echo's API is stable and widely used; can fall back to standard `net/http` if needed
- [slog per-package config not built-in] → Can be added later via `pkg/logger` if required
- [No database integration yet] → Scaffold prepares config but doesn't connect; DB setup is a future change

## Migration Plan

Not applicable — initial scaffold with no prior code.

## Open Questions

- Reverse proxy (traefik/nginx) vs embedded TLS termination
- Container registry for deployment images
