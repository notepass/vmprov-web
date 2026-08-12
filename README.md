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

## Project Structure

```
cmd/server/          # Application entry point
internal/config/     # Configuration loading
internal/logger/     # Slog setup
internal/server/     # Echo HTTP server
```
