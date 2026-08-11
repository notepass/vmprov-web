## Purpose

Provide configuration management via YAML file and environment variables for all application settings.

## Requirements

### Requirement: Configuration struct
The application SHALL define a configuration struct containing server port, DB connect string, DB username, DB password, and log level fields.

#### Scenario: Config struct exists
- **WHEN** the config package is imported
- **THEN** a config struct is available with fields for port, DB connect string, DB username, DB password, and log level

### Requirement: YAML config file loading
The application SHALL load configuration from a `config.yaml` file in the project root.

#### Scenario: Config loaded from file
- **WHEN** `config.yaml` exists in the project root
- **THEN** configuration values are loaded from the file

### Requirement: Environment variable overrides
All config fields SHALL be overridable via environment variables: `SERVER_PORT`, `DB_CONN_STRING`, `DB_USERNAME`, `DB_PASSWORD`, `LOG_LEVEL`.

#### Scenario: Port overridden by env var
- **WHEN** `SERVER_PORT` is set to `3000`
- **THEN** the server binds to port 3000 regardless of config file value

#### Scenario: Log level overridden by env var
- **WHEN** `LOG_LEVEL` is set to `DEBUG`
- **THEN** the logger outputs at DEBUG level

### Requirement: Default configuration values
The application SHALL provide sensible defaults when no config file or env vars are set.

#### Scenario: Default port
- **WHEN** no configuration is provided
- **THEN** the server binds to port 8080

#### Scenario: Default log level
- **WHEN** no configuration is provided
- **THEN** the log level defaults to INFO
