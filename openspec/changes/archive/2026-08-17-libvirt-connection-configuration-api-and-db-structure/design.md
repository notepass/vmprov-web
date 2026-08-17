## Context

The project is a Go-based VM provisioning web server (Echo, viper, sqlx + Squirrel, goose embedded migrations). The DB layer supports PostgreSQL and MySQL/MariaDB via per-dialect migration files and dialect-aware Squirrel placeholder formats. Repositories for users, templates, and audit logs exist but are not yet wired into HTTP handlers — the Echo server currently registers no routes (`cmd/server/main.go` holds the repository instances in `_ =` variables).

RESEARCH.md selected `github.com/digitalocean/go-libvirt` as the libvirt binding: pure Go, speaks libvirt's RPC/XDR protocol directly, no C dependencies. Its API surface used here: `dialers.NewSSH` (options: `UseSSHUsername`, `UseKeyFile`, `UseKnownHostsFile`, `WithAcceptUnknownHostKey`, `WithSSHAuthMethods`), `dialers.NewLocal` (options: `WithSocket`, `WithLocalTimeout`), `libvirt.NewWithDialer` + `Connect`, `ConnectGetLibVersion`, hypervisor-type lookup, `ConnectListAllDomains`, `Disconnect`. SSH auth via `(&dialers.SSHAuthMethods{}).PrivKey()` reads the key file and calls `ssh.ParsePrivateKey`, which rejects passphrase-protected keys.

The application must address multiple libvirt hosts (a fleet). Two transports are in scope: SSH to remote hosts (username + host + SSH private key stored on the app server) and the local libvirt Unix socket (socket path on the app server). Connection endpoints, including credentials, need to be persisted in the database and manageable at runtime through the REST API.

## Goals / Non-Goals

**Goals:**
- Persist named, typed libvirt connections (ssh, socket) with their credentials in the database for both supported dialects
- Provide a REST API for full CRUD plus a connection-test endpoint, with per-type validation
- Manage SSH credentials by server-side key path (key contents never stored in the DB); only non-passphrase-protected keys
- Make the SSH known_hosts file location configurable with a sensible default
- Wrap go-libvirt dialers behind an injectable interface so handlers are testable with fakes
- Bound connection-test dial time with a configurable timeout
- Establish the first HTTP route registration and handler dependency-injection pattern for future endpoints

**Non-Goals:**
- Passphrase-protected SSH keys (explicitly unsupported for now)
- TCP and TLS libvirt transports (SSH + local socket only for now)
- Storing key contents or passwords in the database
- Overriding the remote libvirt socket path for SSH connections (the remote default socket is used)
- VM management (listing/creating/destroying domains) — follow-up change
- Persistent libvirt sessions or connection pooling
- Authentication, authorization, or per-user scoping of connections
- Frontend UI

## Decisions

### Decision 1: go-libvirt behind an injectable `LibvirtClient` interface
`internal/libvirt` defines a `Client` interface (connection test returning libvirt version, hypervisor type, and domain counts, given a typed connection) with a concrete go-libvirt implementation built on the dialer API (`libvirt.NewWithDialer` + `Connect`). Handlers depend on the interface; unit tests substitute a fake.

Alternatives considered:
- Calling go-libvirt directly in handlers — rejected: not mockable, couples transport to the library.
- Persistent pooled sessions keyed by endpoint — rejected: libvirt sessions are cheap to establish; per-request connect/disconnect avoids stale sessions and is simplest. Pooling can be revisited when VM operations (a follow-up change) need it.

### Decision 2: Typed connection model with explicit credentials
Each connection has a `type` (`ssh` or `socket`) and per-type fields:
- `ssh`: `host`, `username`, `ssh_key_path` (private key on the app server), `accept_unknown_host_key` (boolean, default false)
- `socket`: `socket_path` (libvirt Unix socket on the app server)
Plus `name` (unique human-facing label), `description` (optional), and last-check status fields.

Rationale: go-libvirt's dialers take explicit parameters (host, user, key file, socket path), so a typed schema maps 1:1 onto them, enables per-type validation at the API, and keeps credentials out of opaque URI strings.

Alternatives considered:
- Free-form connection URI (`qemu+ssh://user@host/system`) — rejected: hides credentials inside a string, prevents per-field validation and key checks at save time.
- Storing key contents in the DB — rejected: secret management (encryption at rest, rotation) is out of scope; a server-side path reference is sufficient for an internal tool.

Key constraint: only non-passphrase-protected keys are supported — `ssh.ParsePrivateKey` (used by go-libvirt's `PrivKey()` auth method) rejects encrypted keys. Passphrase support is a deliberate non-goal for now.

### Decision 3: SSH dialer configuration and host-key verification
The SSH dial is built as `dialers.NewSSH(host, UseSSHUsername(...), UseKeyFile(ssh_key_path), WithSSHAuthMethods(&dialers.SSHAuthMethods{}.PrivKey()), UseKnownHostsFile(configuredPath), WithAcceptUnknownHostKey() when flagged)`, then `libvirt.NewWithDialer(...)` + `Connect()`.

- `PrivKey()`-only auth methods: deterministic authentication with clear failure when the key is missing, unreadable, or passphrase-protected. go-libvirt's default (agent → key → password fallbacks) would make failures ambiguous and could silently use an unrelated agent key.
- Host keys are verified against the configured known_hosts file. When `accept_unknown_host_key` is true, unknown host keys are auto-appended to the file on the first successful dial (go-libvirt behavior); the flag is per-connection and defaults to false, so strict verification is the default.
- The socket dial is `dialers.NewLocal(WithSocket(socket_path), WithLocalTimeout(timeout))` + `libvirt.NewWithDialer(...)`.

### Decision 4: Known hosts file location — configurable, default `~/.ssh/known_hosts`
A `libvirt_known_hosts_file` config key (env `LIBVIRT_KNOWN_HOSTS_FILE`) sets the known_hosts path used for all SSH dials. The default is `~/.ssh/known_hosts` resolved from the server user's home directory at startup.

- go-libvirt's own default is `~/.config/libvirt/known_hosts` (XDG dir), so the path is always passed explicitly via `UseKnownHostsFile`.
- At startup, if the resolved file does not exist it is created empty (mode 0600, parent dir 0700) so strict mode fails with a host-key error, not a file-not-found error. This mirrors what OpenSSH does on demand.
- `x/crypto/ssh/knownhosts` reads the standard OpenSSH format including hashed entries, so the file can be shared with interactive SSH and `ssh-keyscan`.

### Decision 5: SSH key validated at save time
Create and update requests for `ssh`-type connections reject the request with 400 when the key file does not exist, is not readable by the application process, or does not parse via `ssh.ParsePrivateKey` (i.e. is passphrase-protected). `/test` re-validates at runtime, so a key that later rotates or is removed is caught there.

Rationale: fail fast with an actionable error at the point of configuration; prevents storing permanently broken credentials.

### Decision 6: Per-dialect goose migration `002`
Follows the migration-001 pattern exactly: full copies of `002_create_libvirt_connections.postgres.sql` and `002_create_libvirt_connections.mysql.sql` with `-- +goose Up` / `-- +goose Down` sections. The existing embedded-FS dialect filtering picks the right file at runtime; no migration runner changes needed.

Columns: `id` PK, `name` VARCHAR(255) NOT NULL UNIQUE, `type` VARCHAR(20) NOT NULL (CHECK IN `('ssh','socket')`), `host` VARCHAR(255) NULL, `username` VARCHAR(255) NULL, `ssh_key_path` VARCHAR(1024) NULL, `accept_unknown_host_key` BOOLEAN NOT NULL DEFAULT FALSE, `socket_path` VARCHAR(1024) NULL, `description` TEXT NULL, `last_status` VARCHAR(20) NULL, `last_checked_at` TIMESTAMP NULL, `created_at`, `updated_at`.

Alternative considered: single templated migration — rejected: repo precedent is full per-dialect copies; a template generator is a separate future effort.

### Decision 7: Last-check status recorded server-side
The test endpoint updates `last_status` (`ok` / `error`) and `last_checked_at` on every check. List/get responses include the last-known status so clients can display health without re-dialing. Failure detail is logged but not persisted, to avoid unbounded text growth in the row.

### Decision 8: New `internal/api` package; first route registration
Handlers live in a new `internal/api` package (handler struct holding repository, libvirt client, timeout, known_hosts path, and logger). `internal/server/server.go` registers the `/api/v1/remotes/libvirt/connections` routes and receives dependencies, replacing the `_ = repo` stubs in `main.go` for this capability. REST conventions: JSON bodies, 201 on create, 204 on delete, 400 for validation errors, 404 for missing rows, 409 for duplicate names.

### Decision 9: Config — `LIBVIRT_CONNECT_TIMEOUT` and `LIBVIRT_KNOWN_HOSTS_FILE`
Add `LibvirtConnectTimeout` (seconds, default 10) and `LibvirtKnownHostsFile` (string, default `~/.ssh/known_hosts` resolved from the server user's home) to the config struct with env bindings `LIBVIRT_CONNECT_TIMEOUT` and `LIBVIRT_KNOWN_HOSTS_FILE`. The timeout applies as a dial/RPC context deadline for both transports; the known-hosts path is passed to every SSH dial.

## Risks / Trade-offs

- [go-libvirt warns its API is not stable] → Mitigation: pin the version in `go.mod`; isolate usage behind the `LibvirtClient` interface so a replacement or upgrade is contained to `internal/libvirt`.
- [Test endpoint dials hosts stored in the database] → Mitigation: host, key path, and socket path come from persisted typed rows (not raw request input), the timeout bounds wait time, and no key material is echoed in responses. Network exposure is acceptable for an internal tool until an authentication change lands.
- [Auto-appending unknown host keys weakens verification] → Mitigation: only when the per-connection `accept_unknown_host_key` flag is explicitly true; default is strict verification against known_hosts.
- [App must have read access to key files on the server] → Mitigation: readability is validated at save time with an actionable 400; operators deploy keys to the app server as they would for any service account.
- [Concurrent appends to known_hosts are unguarded in go-libvirt] → Accepted: appends are rare (first test of a new host) and are harmless for this file at our scale.
- [Per-dialect migration duplication] → Mitigation: same accepted trade-off as migration 001; a template-based generator is planned future work.
- [Per-request connect adds latency to `/test`] → Inherent to on-demand health checks; bounded by the configured timeout.

## Migration Plan

No existing data to migrate. Steps:
1. Add `github.com/digitalocean/go-libvirt` dependency
2. Add config fields (`LibvirtConnectTimeout`, `LibvirtKnownHostsFile`), env bindings, defaults, and config tests
3. Add `LibvirtConnection` domain model (typed fields)
4. Add migration `002` for both dialects (up/down)
5. Add repository interface, Squirrel implementation, and mock
6. Add `internal/libvirt` client (SSH + socket dials), key validation helper, fake, and unit tests
7. Add `internal/api` handlers with per-type validation, register routes, add handler tests
8. Wire dependencies in `cmd/server/main.go`, including known_hosts ensure-file at startup
9. Verify: `go vet`, `make test`, `make integrate` (Postgres and MariaDB), manual API check against SSH and socket connections

Rollback: revert routes and handlers, then run goose down to `002`. The change only adds a new table; no existing data is affected.

## Open Questions

None blocking.
