## MODIFIED Requirements

### Requirement: Configuration struct
The application SHALL define a configuration struct containing server port, DB connect string, DB username, DB password, log level, and database pool configuration fields.

#### Scenario: Config struct exists
- **WHEN** the config package is imported
- **THEN** a config struct is available with fields for port, DB connect string, DB username, DB password, log level, and database pool settings (max open conns, max idle conns, conn max lifetime)

### Requirement: Environment variable overrides
All config fields SHALL be overridable via environment variables: `SERVER_PORT`, `DB_CONN_STRING`, `DB_USERNAME`, `DB_PASSWORD`, `LOG_LEVEL`, `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`.

#### Scenario: Port overridden by env var
- **WHEN** `SERVER_PORT` is set to `3000`
- **THEN** the server binds to port 3000 regardless of config file value

#### Scenario: Log level overridden by env var
- **WHEN** `LOG_LEVEL` is set to `DEBUG`
- **THEN** the logger outputs at DEBUG level

#### Scenario: DB pool settings overridden by env var
- **WHEN** `DB_MAX_OPEN_CONNS` is set to `25`
- **THEN** the database connection pool allows a maximum of 25 open connections
