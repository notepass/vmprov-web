## 1. Project Initialization

- [x] 1.1 Initialize Go module with `go mod init github.com/notepass/vmprov-web`
- [x] 1.2 Create directory structure (`cmd/server/`, `internal/config/`, `internal/server/`)
- [x] 1.3 Create Makefile with `build`, `run`, `test` targets

## 2. Configuration Layer

- [x] 2.1 Define config struct (port, DB conn string, DB username, DB password, log level)
- [x] 2.2 Implement YAML config loading with viper
- [x] 2.3 Wire env var overrides (`SERVER_PORT`, `DB_CONN_STRING`, `DB_USERNAME`, `DB_PASSWORD`, `LOG_LEVEL`)
- [x] 2.4 Set defaults (port 8080, log level INFO)

## 3. Logging

- [x] 3.1 Configure slog with JSON encoder
- [x] 3.2 Wire log level from config

## 4. HTTP Server

- [x] 4.1 Add Echo dependency
- [x] 4.2 Implement Echo server creation with graceful shutdown (SIGINT/SIGTERM)
- [x] 4.3 Wire server startup in `cmd/server/main.go`

## 5. Test Foundation

- [x] 5.1 Add testify dependency
- [x] 5.2 Create config loading tests
- [x] 5.3 Create server startup tests

## 6. Verification

- [x] 6.1 Run `make build` to confirm compilation
- [x] 6.2 Run `make test` to verify tests pass
- [x] 6.3 Run `make run` to verify server starts
