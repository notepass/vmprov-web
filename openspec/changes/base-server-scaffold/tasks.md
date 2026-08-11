## 1. Project Initialization

- [ ] 1.1 Initialize Go module with `go mod init github.com/notepass/vmprov-web`
- [ ] 1.2 Create directory structure (`cmd/server/`, `internal/config/`, `internal/server/`)
- [ ] 1.3 Create Makefile with `build`, `run`, `test` targets

## 2. Configuration Layer

- [ ] 2.1 Define config struct (port, DB conn string, DB username, DB password, log level)
- [ ] 2.2 Implement YAML config loading with viper
- [ ] 2.3 Wire env var overrides (`SERVER_PORT`, `DB_CONN_STRING`, `DB_USERNAME`, `DB_PASSWORD`, `LOG_LEVEL`)
- [ ] 2.4 Set defaults (port 8080, log level INFO)

## 3. Logging

- [ ] 3.1 Configure slog with JSON encoder
- [ ] 3.2 Wire log level from config

## 4. HTTP Server

- [ ] 4.1 Add Echo dependency
- [ ] 4.2 Implement Echo server creation with graceful shutdown (SIGINT/SIGTERM)
- [ ] 4.3 Wire server startup in `cmd/server/main.go`

## 5. Test Foundation

- [ ] 5.1 Add testify dependency
- [ ] 5.2 Create config loading tests
- [ ] 5.3 Create server startup tests

## 6. Verification

- [ ] 6.1 Run `make build` to confirm compilation
- [ ] 6.2 Run `make test` to verify tests pass
- [ ] 6.3 Run `make run` to verify server starts
