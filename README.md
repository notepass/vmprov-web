# vmprov-web

VM provisioning web server built with Go.

## Prerequisites

- Go 1.21+

## Quick Start

```bash
make build   # Compile the binary
make run     # Build and start the server
make test    # Run all tests
```

## Configuration

Create a `config.yaml` in the project root:

```yaml
server_port: 8080
db_conn_string: "postgres://localhost:5432/vmprov"
db_username: "admin"
db_password: "secret"
log_level: "INFO"
```

All fields can be overridden via environment variables:

| Field | Environment Variable |
|---|---|
| `server_port` | `SERVER_PORT` |
| `db_conn_string` | `DB_CONN_STRING` |
| `db_username` | `DB_USERNAME` |
| `db_password` | `DB_PASSWORD` |
| `log_level` | `LOG_LEVEL` |

Default values: port `8080`, log level `INFO`.

## Testing

The integration tests require a running PostgreSQL and MariaDB instance. By default the tests expect:

| Database | Host | Port | User | Password | Database |
|---|---|---|---|---|---|
| PostgreSQL | `localhost` | `5432` | `postgres` | `password` | `postgres` |
| MariaDB | `localhost` | `3306` | `root` | `password` | `testdb` |

The `testdb` database will be created automatically if it does not exist. To skip integration tests, run:

```bash
go test -short ./...
```

## Project Structure

```
cmd/server/          # Application entry point
internal/config/     # Configuration loading
internal/logger/     # Slog setup
internal/server/     # Echo HTTP server
```
