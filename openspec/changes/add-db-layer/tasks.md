## 1. Dependencies

- [ ] 1.1 Add `github.com/Masterminds/squirrel` dependency to `go.mod`
- [ ] 1.2 Add `github.com/jmoiron/sqlx` dependency to `go.mod`
- [ ] 1.3 Add `github.com/pressly/goose/v3` dependency to `go.mod`
- [ ] 1.4 Add `github.com/jackc/pgx/v5/stdlib` PostgreSQL driver dependency
- [ ] 1.5 Add `github.com/go-sql-driver/mysql` MySQL/MariaDB driver dependency

## 2. Config updates

- [ ] 2.1 Add `DBMaxOpenConns`, `DBMaxIdleConns`, `DBConnMaxLifetime` fields to config struct
- [ ] 2.2 Bind new config fields to env vars in config loader
- [ ] 2.3 Update config tests to cover new pool fields
- [ ] 2.4 Add default values for new pool configuration fields

## 3. Database adaptor

- [ ] 3.1 Create `internal/db/adaptor.go` with `Adaptor` interface definition
- [ ] 3.2 Create `internal/db/connection.go` with database connection initialization
- [ ] 3.3 Implement `NewAdaptor` constructor with connection string parsing
- [ ] 3.4 Implement `HealthCheck` method on the adaptor
- [ ] 3.5 Implement `Close` method for graceful shutdown
- [ ] 3.6 Add unit tests for adaptor initialization and error handling

## 4. Migrations

- [ ] 4.1 Create `migrations/` directory structure
- [ ] 4.2 Create initial PostgreSQL migration: `001_create_base_tables.up.postgres.sql`
- [ ] 4.3 Create initial PostgreSQL rollback: `001_create_base_tables.down.postgres.sql`
- [ ] 4.4 Create initial MySQL migration: `001_create_base_tables.up.mysql.sql`
- [ ] 4.5 Create initial MySQL rollback: `001_create_base_tables.down.mysql.sql`
- [ ] 4.6 Embed migration files using Go `embed` directive
- [ ] 4.7 Wire goose migration runner into adaptor initialization

## 5. Repositories

- [ ] 5.1 Create `internal/repository/repository.go` with repository interfaces
- [ ] 5.2 Implement template repository with Squirrel-based queries
- [ ] 5.3 Implement user repository with Squirrel-based queries
- [ ] 5.4 Implement audit log repository with Squirrel-based queries
- [ ] 5.5 Add repository unit tests with mock implementations

## 6. Application wiring

- [ ] 6.1 Wire database initialization into `cmd/server/main.go` startup sequence
- [ ] 6.2 Add database health check before HTTP server starts listening
- [ ] 6.3 Register database `Close` in the graceful shutdown handler
- [ ] 6.4 Pass repositories to Echo handlers via dependency injection
- [ ] 6.5 Add integration test verifying full startup and shutdown lifecycle

## 7. Verification

- [ ] 7.1 Run `go vet` and `go test ./...` to verify no regressions
- [ ] 7.2 Verify application starts and connects to a local PostgreSQL instance
- [ ] 7.3 Verify application starts and connects to a local MySQL instance
- [ ] 7.4 Verify migrations run correctly and can be rolled back
