## Why

The project has database configuration placeholders but no persistence layer. We need to store cloud-init templates, users, and audit logs in a database. Most operational state lives in libvirt; the DB is auxiliary but necessary for configuration management and auditability.

## What Changes

- Add Squirrel query builder for cross-database SQL queries (PostgreSQL, MySQL, MariaDB)
- Add sqlx as the database access layer
- Add goose for embedded, versioned schema migrations with up/down support
- Create repository interfaces for data access decoupled from storage details
- Wire database initialization and graceful shutdown into the application lifecycle

## Capabilities

### New Capabilities
- `db-adaptor`: Database connection management, repository pattern, cross-database query support, and embedded schema migrations

### Modified Capabilities
- `config-management`: Add database pool settings (max open connections, max idle connections, connection max lifetime) to existing config

## Impact

- `go.mod`: New dependencies: Squirrel, sqlx, goose, PostgreSQL driver, MySQL driver
- `internal/config/config.go`: New database pool configuration fields
- `internal/`: New `db/` package for connection management and `repository/` package for data access
- `cmd/server/main.go`: Database initialization and graceful shutdown wiring
- `migrations/`: New directory for per-dialect SQL migration files
- Existing config tests updated for new fields
