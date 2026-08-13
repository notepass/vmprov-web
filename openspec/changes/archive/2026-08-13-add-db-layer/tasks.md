## 1. Dependencies

- [x] 1.1 Add `github.com/Masterminds/squirrel` dependency to `go.mod`
- [x] 1.2 Add `github.com/jmoiron/sqlx` dependency to `go.mod`
- [x] 1.3 Add `github.com/pressly/goose/v3` dependency to `go.mod`
- [x] 1.4 Add `github.com/jackc/pgx/v5/stdlib` PostgreSQL driver dependency
- [x] 1.5 Add `github.com/go-sql-driver/mysql` MySQL/MariaDB driver dependency

## 2. Config updates

- [x] 2.1 Add `DBMaxOpenConns`, `DBMaxIdleConns`, `DBConnMaxLifetime` fields to config struct
- [x] 2.2 Bind new config fields to env vars in config loader
- [x] 2.3 Update config tests to cover new pool fields
- [x] 2.4 Add default values for new pool configuration fields

## 3. Database adaptor

- [x] 3.1 Create `internal/db/adaptor.go` with `Adaptor` interface definition
- [x] 3.2 Create `internal/db/connection.go` with database connection initialization
- [x] 3.3 Implement `NewAdaptor` constructor with connection string parsing
- [x] 3.4 Implement `HealthCheck` method on the adaptor
- [x] 3.5 Implement `Close` method for graceful shutdown
- [x] 3.6 Add unit tests for adaptor initialization and error handling

## 4. Migrations

- [x] 4.1 Create `migrations/` directory structure
- [x] 4.2 Create initial PostgreSQL migration: `001_create_base_tables.up.postgres.sql`
- [x] 4.3 Create initial PostgreSQL rollback: `001_create_base_tables.down.postgres.sql`
- [x] 4.4 Create initial MySQL migration: `001_create_base_tables.up.mysql.sql`
- [x] 4.5 Create initial MySQL rollback: `001_create_base_tables.down.mysql.sql`
- [x] 4.6 Embed migration files using Go `embed` directive
- [x] 4.7 Wire goose migration runner into adaptor initialization

## 5. Repositories

- [x] 5.1 Create `internal/repository/repository.go` with repository interfaces
- [x] 5.2 Implement template repository with Squirrel-based queries
- [x] 5.3 Implement user repository with Squirrel-based queries
- [x] 5.4 Implement audit log repository with Squirrel-based queries
- [x] 5.5 Add repository unit tests with mock implementations

## 6. Application wiring

- [x] 6.1 Wire database initialization into `cmd/server/main.go` startup sequence
- [x] 6.2 Add database health check before HTTP server starts listening
- [x] 6.3 Register database `Close` in the graceful shutdown handler
- [x] 6.4 Pass repositories to Echo handlers via dependency injection
- [x] 6.5 Add integration test verifying full startup and shutdown lifecycle

## 7. Verification

- [x] 7.1 Run `go vet` and `go test ./...` to verify no regressions
- [x] 7.2 Verify application starts and connects to a local PostgreSQL instance
- [x] 7.3 Verify application starts and connects to a local MySQL instance
- [x] 7.4 Verify migrations run correctly and can be rolled back
