## Why

The application will manage VMs on libvirt hosts, but there is no way to configure which libvirt endpoints to talk to. Connection settings need to be persisted in the database and manageable at runtime through a REST API, so the app can address multiple libvirt hosts (e.g. a fleet) without redeployment.

## What Changes

- Add a `libvirt_connections` table storing named, typed libvirt endpoints via per-dialect goose migration: SSH connections (host, username, SSH key path, accept-unknown-host-key flag) and local socket connections (socket path), plus description and last check status
- Add credential management for SSH connections: the private key is referenced by a path on the application server, must not be passphrase-protected, and is validated (exists, readable, parseable) when a connection is saved
- Make the SSH known_hosts file location configurable (`libvirt_known_hosts_file` / `LIBVIRT_KNOWN_HOSTS_FILE`), defaulting to `~/.ssh/known_hosts` of the user running the server
- Add `LibvirtConnection` domain model and a Squirrel-based repository (CRUD + unique-name lookups) following the existing repository pattern
- Add REST API under `/api/v1/remotes/libvirt/connections` for listing, creating, getting, updating, and deleting connections, with per-type field validation
- Add `POST /api/v1/remotes/libvirt/connections/:id/test` endpoint that dials the configured libvirt endpoint (SSH or local socket) and reports libvirt version, hypervisor type, and domain counts
- Add a libvirt client package wrapping `github.com/digitalocean/go-libvirt` dialers (`dialers.NewSSH`, `dialers.NewLocal`) behind an injectable interface for testability
- Add `LIBVIRT_CONNECT_TIMEOUT` configuration (env var + default) bounding connection-test dial time
- Register the first real HTTP routes in the Echo server (server currently exposes none)

## Capabilities

### New Capabilities
- `libvirt-connections`: Persistence and REST management of typed libvirt connections (SSH and local socket) with server-side SSH key credential management, plus connection verification against a live libvirt host

### Modified Capabilities
- `config-management`: Add `LIBVIRT_CONNECT_TIMEOUT` and `LIBVIRT_KNOWN_HOSTS_FILE` environment variable overrides, config file keys, and default values

## Impact

- `go.mod`: New dependency `github.com/digitalocean/go-libvirt`
- `internal/domain/models.go`: New `LibvirtConnection` model (typed fields)
- `internal/db/migrations/`: New `002_create_libvirt_connections.postgres.sql` and `002_create_libvirt_connections.mysql.sql`
- `internal/repository/`: New `LibvirtConnectionRepository` interface and SQL implementation, new mock
- `internal/libvirt/`: New package with `LibvirtClient` interface, go-libvirt dialer-based implementation (SSH + local socket), and SSH key validation
- `internal/api/`: New package with libvirt-connection handlers (first HTTP routes)
- `internal/server/server.go`: Route registration for `/api/v1/remotes/libvirt/connections`
- `internal/config/`: New `LibvirtConnectTimeout` and `LibvirtKnownHostsFile` fields, env bindings, defaults
- `cmd/server/main.go`: Wire repository, client, and handlers into the server; ensure known_hosts file exists at startup
