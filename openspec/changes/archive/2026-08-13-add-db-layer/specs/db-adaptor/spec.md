## ADDED Requirements

### Requirement: Database adaptor interface
The system SHALL define a `db.Adaptor` interface in `internal/db/` that provides connection management and health checking.

#### Scenario: Adaptor interface exists
- **WHEN** the `db` package is imported
- **THEN** an `Adaptor` interface is available with methods for health checking and closing the connection

#### Scenario: Adaptor can be implemented by multiple backends
- **WHEN** a new struct implements all `Adaptor` methods
- **THEN** it can be used in place of the default implementation without changing consumer code

### Requirement: Cross-database query support
The system SHALL use Squirrel query builder to generate SQL queries compatible with PostgreSQL, MySQL, and MariaDB.

#### Scenario: Query works across databases
- **WHEN** a repository method executes a SELECT query via Squirrel
- **THEN** the generated SQL syntax is correct for the target database dialect

### Requirement: Database connection initialization
The application SHALL initialize the database connection and run pending migrations during startup.

#### Scenario: Database initializes on startup
- **WHEN** the application starts with valid database configuration
- **THEN** the database connection pool is established and migrations are applied before the HTTP server begins listening

#### Scenario: Application fails startup on database error
- **WHEN** the application starts and the database is unreachable
- **THEN** the application logs the error and exits with a non-zero code

### Requirement: Graceful database shutdown
The application SHALL close the database connection pool during graceful shutdown.

#### Scenario: Database closes on shutdown
- **WHEN** the application receives a termination signal
- **THEN** the database connection pool is closed before the process exits

### Requirement: Repository pattern for data access
The system SHALL provide repository interfaces in `internal/repository/` that abstract CRUD operations over domain entities.

#### Scenario: Repository interface exists
- **WHEN** the `repository` package is imported
- **THEN** repository interfaces for domain entities are available

#### Scenario: Repository is testable with mocks
- **WHEN** a unit test needs to test business logic depending on a repository
- **THEN** a mock implementation of the repository interface can be substituted

### Requirement: Embedded schema migrations
The system SHALL embed migration SQL files in the application binary using Go's `embed` package.

#### Scenario: Migrations are embedded
- **WHEN** the application is built
- **THEN** migration files from `migrations/` are included in the binary and accessible at runtime

### Requirement: Migration up/down support
The system SHALL support forward migrations, rollback migrations, and jumping to any migration version.

#### Scenario: Forward migration applies
- **WHEN** the application starts and pending migrations exist
- **THEN** migrations are applied in order and recorded in the migration tracking table

#### Scenario: Migration rollback works
- **WHEN** a rollback to a specific version is requested
- **THEN** down scripts are executed in reverse order until the target version is reached

### Requirement: Database health check
The adaptor SHALL provide a health check method that verifies database connectivity.

#### Scenario: Health check reports healthy
- **WHEN** the database is reachable
- **THEN** the health check method returns no error

#### Scenario: Health check reports unhealthy
- **WHEN** the database is unreachable
- **THEN** the health check method returns an error
