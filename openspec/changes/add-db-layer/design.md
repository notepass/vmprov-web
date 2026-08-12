## Context

The project is a Go-based VM provisioning web server. Database configuration fields exist in the config struct but no actual database integration is implemented. The application needs to store cloud-init templates, users, and audit logs. Most state lives in libvirt; the DB is auxiliary.

## Goals / Non-Goals

**Goals:**
- Support PostgreSQL, MySQL, and MariaDB as database backends
- Provide cross-database query support via Squirrel query builder
- Enable schema migrations with up/down/jump-to-version support via goose
- Keep data access decoupled from business logic via repository pattern
- Embed migration files in the application binary for single-binary deployment
- Allow the adaptor to be swapped without changing consumer code

**Non-Goals:**
- Multi-database support beyond PostgreSQL, MySQL, and MariaDB in this change
- Query caching or advanced connection pooling beyond standard `database/sql` defaults
- Admin UI or database management interfaces
- Real-time schema reflection or auto-migration

## Decisions

### Decision 1: Squirrel + sqlx for queries
Squirrel is a query builder that handles dialect differences for DML operations. Combined with sqlx for struct mapping, this gives us:
- Cross-database SQL without an ORM's implicit behavior
- Familiar SQL-like syntax
- Minimal abstraction overhead
- No N+1 traps or callback hooks to learn

### Decision 2: goose for migrations
goose is an embeddable migration library that supports:
- Per-dialect SQL files (full copies, not templated)
- Up, down, and jump-to-version operations
- Embedded migration files via Go 1.16+ `embed`
- Both SQL and Go-based migrations

Future work: build lightweight tooling around migration templates to reduce per-dialect duplication.

### Decision 3: Repository pattern with interface-based adaptor
Define repository interfaces in `internal/repository/` with concrete implementations in `internal/repository/postgres/`, `internal/repository/mysql/`, etc. Business logic depends on interfaces, not concrete types.

### Decision 4: Embedded migrations
Use Go's `embed` package to include migration files in the binary. This enables single-binary deployment and avoids file system dependencies at runtime.

## Risks / Trade-offs

[Per-dialect migration files duplicate content] → Mitigation: DDL structure is nearly identical; only column types differ. Minimal duplication. Plan for template-based generator in the future.
[Squirrel's DDL support is limited] → Mitigation: Migrations use raw SQL files (goose), Squirrel handles runtime DML only
[Tight coupling to specific databases] → Mitigation: Repository interfaces define the API surface; swapping backends requires only new implementations

## Migration Plan

No existing data to migrate. Steps:
1. Add Squirrel, sqlx, goose, and database driver dependencies
2. Create `internal/db/` with connection management
3. Create `migrations/` with initial per-dialect SQL files
4. Create `internal/repository/` with repository interfaces and implementations
5. Wire goose initialization in `cmd/server/main.go`
6. Add database close to graceful shutdown handler
7. Update config with pool settings
8. Add integration tests

Rollback: Revert dependency changes and remove new packages. No data loss risk.
