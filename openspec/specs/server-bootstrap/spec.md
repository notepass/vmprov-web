## Purpose

Establish the base Go server scaffold including project structure, HTTP server, structured logging, and build tooling.

## Requirements

### Requirement: Go module initialization
The project SHALL have a valid `go.mod` file declaring the module path `github.com/notepass/vmprov-web` with Go version constraints.

#### Scenario: Module exists
- **WHEN** the developer runs `go mod verify`
- **THEN** the command succeeds with no errors

### Requirement: Project directory structure
The project SHALL follow standard Go project layout with `cmd/server/` for the entry point and `internal/` for private packages.

#### Scenario: Standard layout present
- **WHEN** the developer inspects the root directory
- **THEN** directories `cmd/server` and `internal` exist with appropriate files

### Requirement: HTTP server startup
The application SHALL start an Echo HTTP server on a configurable port and listen for incoming requests.

#### Scenario: Server starts on default port
- **WHEN** the application is run without configuration
- **THEN** the HTTP server listens on port 8080

#### Scenario: Server starts on configured port
- **WHEN** `SERVER_PORT` is set to `9090`
- **THEN** the HTTP server listens on port 9090

### Requirement: Graceful shutdown
The application SHALL handle SIGINT and SIGTERM signals and gracefully shut down the HTTP server.

#### Scenario: Graceful shutdown on SIGTERM
- **WHEN** the process receives SIGTERM
- **THEN** the server stops accepting new connections and completes in-flight requests within a timeout

### Requirement: Structured logging
The application SHALL use slog for structured JSON log output with configurable log levels.

#### Scenario: JSON log output
- **WHEN** the server starts
- **THEN** log entries are valid JSON objects with timestamp, level, and message fields

### Requirement: Makefile build target
The project SHALL provide a Makefile with a `build` target that compiles the binary.

#### Scenario: Build produces binary
- **WHEN** the developer runs `make build`
- **THEN** a compiled binary is produced

### Requirement: Makefile run target
The project SHALL provide a `run` target that builds and starts the server.

#### Scenario: Run starts server
- **WHEN** the developer runs `make run`
- **THEN** the server starts and listens on the configured port

### Requirement: Makefile test target
The project SHALL provide a `test` target that runs all tests.

#### Scenario: Test runs tests
- **WHEN** the developer runs `make test`
- **THEN** all Go test files are executed
