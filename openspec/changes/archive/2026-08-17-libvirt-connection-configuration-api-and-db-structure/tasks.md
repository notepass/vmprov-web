## 1. Dependencies

- [x] 1.1 Add `github.com/digitalocean/go-libvirt` dependency to `go.mod`

## 2. Config updates

- [x] 2.1 Add `LibvirtConnectTimeout` and `LibvirtKnownHostsFile` fields to the config struct in `internal/config/config.go`
- [x] 2.2 Bind `LIBVIRT_CONNECT_TIMEOUT` and `LIBVIRT_KNOWN_HOSTS_FILE` env vars in `internal/config/loader.go`
- [x] 2.3 Add defaults: `DefaultLibvirtConnectTimeout` (10s) and known hosts file default `~/.ssh/known_hosts` resolved from the server user's home directory
- [x] 2.4 Update config tests to cover the new fields, env overrides, and defaults

## 3. Domain model

- [x] 3.1 Add `LibvirtConnection` struct to `internal/domain/models.go` with `ID`, `Name`, `Type`, `Host`, `Username`, `SSHKeyPath`, `AcceptUnknownHostKey`, `SocketPath`, `Description`, `LastStatus`, `LastCheckedAt`, `CreatedAt`, `UpdatedAt` (nullable pointers for `Host`, `Username`, `SSHKeyPath`, `SocketPath`, `Description`, `LastStatus`, `LastCheckedAt`)

## 4. Migrations

- [x] 4.1 Create `internal/db/migrations/002_create_libvirt_connections.postgres.sql` with up script (`libvirt_connections` table: unique `name`, `type` CHECK IN ('ssh','socket'), per-type nullable columns, `accept_unknown_host_key` BOOLEAN NOT NULL DEFAULT FALSE) and down script
- [x] 4.2 Create `internal/db/migrations/002_create_libvirt_connections.mysql.sql` with up and down scripts
- [x] 4.3 Add integration test verifying migration `002` applies and rolls back on PostgreSQL
- [x] 4.4 Add integration test verifying migration `002` applies and rolls back on MySQL

## 5. Repository

- [x] 5.1 Add `LibvirtConnectionRepository` interface to `internal/repository/repository.go` (`Create`, `GetByID`, `GetByName`, `Update`, `Delete`, `List`)
- [x] 5.2 Implement `LibvirtConnectionRepo` with Squirrel queries in `internal/repository/libvirt_connection_repo.go` following the `TemplateRepo` pattern (dialect-aware placeholder format)
- [x] 5.3 Add `MockLibvirtConnectionRepository` to `internal/repository/mock_test.go`
- [x] 5.4 Add repository unit tests covering CRUD, unique-name lookup, and not-found cases

## 6. Libvirt client

- [x] 6.1 Define `Client` interface in `internal/libvirt/client.go` with a connection-test method taking a typed connection (type plus per-type fields, known_hosts file path, and timeout) and returning libvirt version, hypervisor type, and total/active domain counts
- [x] 6.2 Implement the SSH dial: `dialers.NewSSH` with `UseSSHUsername`, `UseKeyFile`, `WithSSHAuthMethods` (PrivKey only), `UseKnownHostsFile`, and `WithAcceptUnknownHostKey` (when flagged), then `libvirt.NewWithDialer` + `Connect`, version, hypervisor type, `ConnectListAllDomains`, `Disconnect`
- [x] 6.3 Implement the socket dial: `dialers.NewLocal` with `WithSocket` and `WithLocalTimeout`, then `libvirt.NewWithDialer` + `Connect`
- [x] 6.4 Add an SSH key validation helper (file exists, readable, parses via `ssh.ParsePrivateKey` — passphrase-protected keys fail)
- [x] 6.5 Add unit tests for the client covering missing/unreadable/passphrase-protected keys, dial failure, and disconnect-on-failure behavior
- [x] 6.6 Add a fake `Client` implementation for handler tests

## 7. API handlers and routes

- [x] 7.1 Create `internal/api` package with a libvirt connection handler struct (repository, client, timeout, known_hosts path, logger) and request/response types
- [x] 7.2 Implement `GET /api/v1/remotes/libvirt/connections` (list)
- [x] 7.3 Implement `POST /api/v1/remotes/libvirt/connections` (create; per-type required-field validation, 400 on invalid type or missing/unreadable/passphrase-protected SSH key, 409 on duplicate name, 201 on success)
- [x] 7.4 Implement `GET /api/v1/remotes/libvirt/connections/:id` (200/404, 400 on non-numeric ID)
- [x] 7.5 Implement `PUT /api/v1/remotes/libvirt/connections/:id` (200/404/409, same validation as create)
- [x] 7.6 Implement `DELETE /api/v1/remotes/libvirt/connections/:id` (204/404)
- [x] 7.7 Implement `POST /api/v1/remotes/libvirt/connections/:id/test` (dials per type: 200 with version/hypervisor/domain counts, 502 on failure/timeout/host-key error, 404 on missing row; persist `last_status` and `last_checked_at`)
- [x] 7.8 Register the routes in `internal/server/server.go` with dependency injection
- [x] 7.9 Add handler unit tests using mock repository and fake libvirt client covering all endpoints and error cases (incl. per-type validation and key-validation failures)

## 8. Application wiring

- [x] 8.1 Wire `LibvirtConnectionRepo`, libvirt client, timeout, and known_hosts path into the API handlers in `cmd/server/main.go`
- [x] 8.2 At startup, resolve the known_hosts file path (default `~/.ssh/known_hosts`) and create the file (mode 0600, parent dir 0700) if it does not exist
- [x] 8.3 Remove the now-wired repository stubs (`_ = ...`) for libvirt connections from `cmd/server/main.go`

## 9. Verification

- [x] 9.1 Run `go vet ./...` and `make test` to verify no regressions
- [x] 9.2 Run `make integrate` to verify migrations and repository behavior on PostgreSQL and MariaDB
- [x] 9.3 Verify all six endpoints manually against a running instance (list/create/get/update/delete/test, including 404/400/409/502 paths and per-type validation)
- [x] 9.4 Verify an SSH connection test against a live libvirt host (key auth, strict known_hosts, and accept-unknown behavior) and a local socket connection — local socket verified against live libvirt 12.6.0 (QEMU, version + domain counts, password-prompt auth); SSH portion now covered by the CI `libvirt-integration` job added in change `libvirt-pipeline-test` (archived 2026-08-18), verified via GitHub Actions run 2026-08-18
