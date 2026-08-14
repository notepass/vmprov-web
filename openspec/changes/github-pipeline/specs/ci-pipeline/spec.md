## ADDED Requirements

### Requirement: CI workflow triggers manually
The project SHALL provide a GitHub Actions workflow that is triggered only by `workflow_dispatch`.

#### Scenario: Workflow dispatch available
- **WHEN** a user triggers the workflow via GitHub Actions UI or `gh workflow run`
- **THEN** the workflow starts and executes all configured jobs

#### Scenario: No automatic triggers
- **WHEN** a push or pull request event occurs
- **THEN** the CI workflow does not run

### Requirement: Unit tests execute successfully
The workflow SHALL run `go test ./... -short` to execute all unit tests.

#### Scenario: Unit tests pass
- **WHEN** unit tests are executed in the workflow
- **THEN** all tests in packages without database dependencies pass

#### Scenario: Unit tests fail
- **WHEN** a unit test fails
- **THEN** the workflow job fails and reports the error

### Requirement: Integration tests run with containerized databases
The workflow SHALL start PostgreSQL 18 and MariaDB 12 containers before running integration tests via `go test ./...`.

#### Scenario: PostgreSQL integration tests pass
- **WHEN** the PostgreSQL 18 container is running with `POSTGRES_PASSWORD=password`
- **THEN** integration tests connecting to `postgres://postgres:password@localhost:5432/postgres?sslmode=disable` succeed

#### Scenario: MariaDB integration tests pass
- **WHEN** the MariaDB 12 container is running with `MYSQL_ROOT_PASSWORD=password` and `MYSQL_DATABASE=testdb`
- **THEN** integration tests connecting to `root:password@tcp(localhost:3306)/testdb` succeed

#### Scenario: Database startup failure
- **WHEN** a database container fails to start
- **THEN** the workflow job fails with a clear error message

### Requirement: Binaries build for target platforms
The workflow SHALL build the Go binary for linux/amd64 and linux/arm64.

#### Scenario: linux/amd64 binary builds
- **WHEN** the build job runs with `GOOS=linux GOARCH=amd64`
- **THEN** a functional binary is produced at `bin/server`

#### Scenario: linux/arm64 binary builds
- **WHEN** the build job runs with `GOOS=linux GOARCH=arm64`
- **THEN** a functional binary is produced at `bin/server`

#### Scenario: Build artifacts uploaded
- **WHEN** binaries are built successfully
- **THEN** they are uploaded as workflow artifacts for download

### Requirement: Dockerfile produces working image
The project SHALL include a Dockerfile that builds the Go application into a minimal container image.

#### Scenario: Dockerfile builds
- **WHEN** `docker build -t vmprov-web .` is run
- **THEN** a Docker image is produced containing the Go binary

#### Scenario: Multi-stage build
- **WHEN** the Dockerfile is inspected
- **THEN** it uses a multi-stage build with a Go builder stage and a minimal runtime stage
