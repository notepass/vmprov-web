## Context

The previous change (`libvirt-connection-configuration-api-and-db-structure`, archived 2026-08-17) shipped the libvirt connection stack: `libvirt_connections` persistence, the REST API including `POST /api/v1/remotes/libvirt/connections/:id/test`, and the `internal/libvirt` client wrapping `github.com/digitalocean/go-libvirt` dialers (`dialers.NewSSH`, `dialers.NewLocal`) behind an injectable `Client` interface with a fake for unit tests.

Verification to date: 22 manual endpoint checks against a running instance (all status codes, validation, `last_status` persistence, 5s timeout enforcement) and a local socket conn-test against libvirt 12.6.0 — the 3s `Connect()` delay was an interactive polkit password prompt, which group-based socket access avoids. The archived task 9.4 (SSH verification against a live libvirt host) remains open because no SSH-reachable libvirt host was available.

Runner research (official `actions/runner-images` readmes, August 2026):
- `ubuntu-latest` is Ubuntu 24.04 (26.04 is still a public preview)
- libvirt/qemu are NOT listed in the image readme's installed packages
- `openssh-server` is NOT installed (only `openssh-client`)
- passwordless `sudo` is available
- `/var/run/libvirt/libvirt-sock` is `root:libvirt` mode 0770; users in the `libvirt` group get direct access without a polkit prompt
- no KVM is required for connection-level tests (QEMU falls back to TCG; we never start a domain)

## Goals / Non-Goals

**Goals:**
- Automated live conn-testing of both transports (local socket, SSH) in CI on every manual dispatch
- Cover strict host-key verification (success and rejection) and accept-unknown append behavior — the paths that unit tests with the fake client cannot reach
- Close the verification gap of archived task 9.4 without depending on any external host
- Keep `go test ./... -short` and local dev runs green on machines without libvirt

**Non-Goals:**
- No VM lifecycle testing (domain create/destroy) — connection level only, per the agreed scope
- No end-to-end test through the HTTP API + database — the HTTP/DB path is already covered by the existing containerized-database integration tests; this change tests the conn semantics in `internal/libvirt`
- No changes to production code, configuration, or API
- No self-hosted runners, macOS/Windows runners, or remote SSH hosts
- No libvirt on the CI unit/integration jobs (separate, self-contained job)

## Decisions

### Decision 1: Client-level tests only
The live tests live in `internal/libvirt` and call the concrete client (`New()`) directly with a `Connection` struct.

Rationale: connection-test semantics (dial, version, hypervisor type, domain counts, host-key handling) all live in `internal/libvirt`; the HTTP/DB layer around it is already integration-tested. Client-level tests are faster, need no database, and localize failures to the transport layer.

Alternatives considered:
- E2E via running server + docker DB + `POST .../test` — rejected for this change (heavier, slower, couples transport regressions to DB availability); can be added later as a separate change.

### Decision 2: Environment-variable gate, not `testing.Short()`
Live tests skip unless `LIBVIRT_INTEGRATION` is non-empty.

Rationale: `-short` runs on every dev machine and in the CI `unit` job, none of which have libvirt; a dedicated variable mirrors the opt-in style of `make integrate` and makes intent explicit.

### Decision 3: Explicit package installation in the job
The job apt-installs `libvirt-daemon-system`, `libvirt-clients`, `libvirt-daemon-driver-qemu`, `qemu-kvm`, `openssh-server` instead of relying on the runner image.

Rationale: the current image readme lists neither libvirt nor openssh-server; an explicit install is idempotent and resilient to image changes in both directions.

### Decision 4: Group-based socket access via `usermod` + `sg`
The job runs `sudo usermod -aG libvirt "$(whoami)"` and executes the tests via `sg libvirt -c '...'`.

Rationale: running tests as root would bypass the exact permission semantics the production app relies on; `sg` reads `/etc/group` at invocation, so the `usermod` from a prior step takes effect without re-login. Both direct socket dials and SSH-forwarded dials (same remote user) are covered by one mechanism.

Alternatives considered:
- `sudo go test` — rejected: root access masks socket-permission regressions.
- `chmod 0777` on the socket — rejected: fragile if `libvirtd` restarts, and broader than needed.

### Decision 5: SSH over `127.0.0.1`
The SSH test dials the runner's own `sshd` at `127.0.0.1`.

Rationale: exercises the full go-libvirt SSH dialer path — key authentication, local port-forwarding of the libvirt socket, host-key verification — with zero external dependencies.

### Decision 6: Timeouts and startup waits
Tests use a 15–20s connection timeout (CI is slower than a dev laptop); the job waits for `libvirtd` to be active (systemd active check / `virsh version` loop) before testing.

### Decision 7: known_hosts fixtures
- Strict-mode test: fixture built at job time via `ssh-keyscan 127.0.0.1` (deterministic within the run).
- Accept-unknown test: fresh empty writable file in `t.TempDir()` (no cross-test interference); the test asserts the host key entry appears afterwards.
- Strict-rejection test: fresh empty known_hosts file.

## Risks / Trade-offs

- [libvirtd slow to become active on a fresh runner] → active-wait loop in the job before the test step; generous test timeout.
- [go-libvirt accept-unknown append behavior unverified in our code path] → Decision 7's test asserts the file gains a host key entry; if the library does not append, the test surfaces it and we adjust the client (e.g., pre-seed) in this change.
- [Runner image package churn (libvirt/sshd appear or move)] → explicit idempotent install keeps the job correct either way.
- [KVM absent on the runner] → irrelevant for conn tests; TCG fallback is sufficient.
- [`sshd` host keys missing] → the package generates them on install; the job's `systemctl start ssh` covers it.

## Open Questions

- None blocking. If the runner's `sshd` refuses to start (image policy), fallback: run the SSH test against a dedicated `sshd` instance on a non-standard port configured via `LIBVIRT_SSH_HOST`/port env vars.
