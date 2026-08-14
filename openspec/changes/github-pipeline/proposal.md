## Why

The project lacks CI automation. Every build and test run requires manual execution, risking regressions going undetected. A CI pipeline ensures consistent verification on demand before merging changes.

## What Changes

- Add a GitHub Actions workflow (`.github/workflows/ci.yml`) triggered manually via `workflow_dispatch`
- Add `docker-compose.yml` defining PostgreSQL 18 and MariaDB 12 services for integration testing
- Add `Dockerfile` for building the Go binary
- Update `Makefile` with an `integrate` target to run integration tests against the dockerized databases
- The workflow runs unit tests, starts database containers, runs integration tests, then builds linux/amd64 and linux/arm64 binaries

## Capabilities

### New Capabilities
- `ci-pipeline`: GitHub Actions workflow that builds binaries and runs all tests (unit and integration) with containerized databases

### Modified Capabilities
(none)

## Impact

- New files: `.github/workflows/ci.yml`, `docker-compose.yml`, `Dockerfile`
- Modified: `Makefile` (new `integrate` target)
- Integration tests remain gated by `testing.Short()`; CI runs them without the flag
- No changes to application code or existing test logic
