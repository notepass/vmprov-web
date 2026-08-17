## Why

Libvirt connection testing is implemented (CRUD API, `POST /api/v1/remotes/libvirt/connections/:id/test`, and the `internal/libvirt` client), but it has never been exercised against a real libvirt daemon in CI — verification so far was manual only, and the SSH-against-a-live-host check from the previous change (archived task 9.4) is still open. A GitHub Actions ubuntu runner can host a real `libvirtd` and `sshd`, so a client-level connection test (no VMs required) can run on every manual CI dispatch and guard the dialer/RPC layer, SSH key authentication, and known_hosts host-key behavior against regressions.

## What Changes

- Add gated live integration tests to `internal/libvirt` (skipped unless `LIBVIRT_INTEGRATION=1`, so `go test ./... -short` and local dev runs stay green):
  - local socket connection test asserting libvirt version, hypervisor type, and domain counts
  - SSH connection test to `127.0.0.1` with strict host-key verification against a provided known_hosts file
  - negative test: unknown host key with `accept_unknown_host_key=false` is rejected
  - positive test: unknown host key with `accept_unknown_host_key=true` succeeds and the host key is appended to the known_hosts file
- Add a `libvirt-integration` job to the existing `.github/workflows/ci.yml` (trigger remains `workflow_dispatch`-only): install `libvirt-daemon-system`, `libvirt-clients`, `libvirt-daemon-driver-qemu`, `qemu-kvm`, and `openssh-server` (none are preinstalled on the runner image), start `libvirtd` and `sshd`, grant the runner user libvirt-group socket access, set up an ed25519 key plus known_hosts fixture, and run the gated tests via `sg libvirt`. The job needs no containerized databases.
- Add a `make libvirt-integrate` target so the gated tests can also run locally against a machine's own libvirtd/sshd.

## Capabilities

### New Capabilities
- (none)

### Modified Capabilities
- `libvirt-connections`: Add requirements for gated live client-level integration tests (local socket and SSH via localhost, including strict host-key rejection and accept-unknown append behavior)
- `ci-pipeline`: Add requirements for a libvirt integration job (package installation, service startup, libvirt-group socket access, SSH key/known_hosts fixtures, gated test execution) within the manually-triggered workflow

## Impact

- `internal/libvirt/`: new test file (no production code changes)
- `.github/workflows/ci.yml`: new `libvirt-integration` job
- `Makefile`: new `libvirt-integrate` target
- No new Go dependencies, no API or schema changes
