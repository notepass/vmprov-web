## MODIFIED Requirements

### Requirement: Configuration struct
The application SHALL define a configuration struct containing server port, DB connect string, DB username, DB password, log level, database pool configuration fields, libvirt connect timeout, and libvirt known hosts file path.

#### Scenario: Config struct exists
- **WHEN** the config package is imported
- **THEN** a config struct is available with fields for port, DB connect string, DB username, DB password, log level, database pool settings (max open conns, max idle conns, conn max lifetime), libvirt connect timeout, and libvirt known hosts file path

### Requirement: Environment variable overrides
All config fields SHALL be overridable via environment variables: `SERVER_PORT`, `DB_CONN_STRING`, `DB_USERNAME`, `DB_PASSWORD`, `LOG_LEVEL`, `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`, `LIBVIRT_CONNECT_TIMEOUT`, `LIBVIRT_KNOWN_HOSTS_FILE`.

#### Scenario: Port overridden by env var
- **WHEN** `SERVER_PORT` is set to `3000`
- **THEN** the server binds to port 3000 regardless of config file value

#### Scenario: Log level overridden by env var
- **WHEN** `LOG_LEVEL` is set to `DEBUG`
- **THEN** the logger outputs at DEBUG level

#### Scenario: DB pool settings overridden by env var
- **WHEN** `DB_MAX_OPEN_CONNS` is set to `25`
- **THEN** the database connection pool allows a maximum of 25 open connections

#### Scenario: Libvirt connect timeout overridden by env var
- **WHEN** `LIBVIRT_CONNECT_TIMEOUT` is set to `30`
- **THEN** libvirt connection tests use a 30 second timeout regardless of config file value

#### Scenario: Libvirt known hosts file overridden by env var
- **WHEN** `LIBVIRT_KNOWN_HOSTS_FILE` is set to `/etc/vmprov/known_hosts`
- **THEN** SSH connections verify host keys against `/etc/vmprov/known_hosts` regardless of config file value

### Requirement: Default configuration values
The application SHALL provide sensible defaults when no config file or env vars are set.

#### Scenario: Default port
- **WHEN** no configuration is provided
- **THEN** the server binds to port 8080

#### Scenario: Default log level
- **WHEN** no configuration is provided
- **THEN** the log level defaults to INFO

#### Scenario: Default libvirt connect timeout
- **WHEN** no configuration is provided
- **THEN** libvirt connection tests use a 10 second timeout

#### Scenario: Default libvirt known hosts file
- **WHEN** no configuration is provided
- **THEN** the known hosts file defaults to `~/.ssh/known_hosts` resolved from the home directory of the user running the server
