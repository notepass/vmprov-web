## ADDED Requirements

### Requirement: Gated live integration tests
The `internal/libvirt` package SHALL provide live connection integration tests that are skipped unless the `LIBVIRT_INTEGRATION` environment variable is set to a non-empty value, so that default and `-short` test runs remain green on machines without a libvirt daemon.

#### Scenario: Tests skip without the gate
- **WHEN** `go test ./internal/libvirt/...` runs without `LIBVIRT_INTEGRATION` set
- **THEN** the live integration tests are skipped and the test run passes

#### Scenario: Tests run with the gate
- **WHEN** `go test ./internal/libvirt/...` runs with `LIBVIRT_INTEGRATION=1` on a machine with a running libvirt daemon
- **THEN** the live integration tests execute and dial the configured libvirt endpoints

### Requirement: Local socket connection test
The live integration tests SHALL include a test that dials the local libvirt socket (path configurable via the `LIBVIRT_SOCKET_PATH` environment variable, defaulting to `/var/run/libvirt/libvirt-sock`) using the concrete go-libvirt-based client and asserts on the returned test result.

#### Scenario: Socket connection succeeds
- **WHEN** the local socket test runs and the libvirt daemon is reachable
- **THEN** the test result contains a dotted-numeric libvirt version (matching `^\d+\.\d+\.\d+$`), the hypervisor type `QEMU`, and total domain counts greater than or equal to active domain counts

### Requirement: SSH connection test via localhost
The live integration tests SHALL include a test that dials the libvirt daemon over SSH to `127.0.0.1` (host, username, key path, and known_hosts file configurable via the `LIBVIRT_SSH_HOST`, `LIBVIRT_SSH_USER`, `LIBVIRT_SSH_KEY`, and `LIBVIRT_KNOWN_HOSTS` environment variables) using a generated SSH private key with strict host-key verification.

#### Scenario: SSH connection with known host key succeeds
- **WHEN** the SSH test runs with a known_hosts file that already contains the host key for `127.0.0.1`
- **THEN** the connection test succeeds and returns a dotted-numeric libvirt version and the hypervisor type `QEMU`

### Requirement: Unknown host key rejected in strict mode
The live integration tests SHALL include a negative test verifying that an SSH connection whose host key is absent from the known_hosts file fails host-key verification when `accept_unknown_host_key` is false.

#### Scenario: Strict mode rejects unknown host key
- **WHEN** the SSH test dials with an empty known_hosts file and `accept_unknown_host_key` set to false
- **THEN** the connection test fails with an error indicating the host key could not be verified

### Requirement: Accept-unknown host key appends to known hosts
The live integration tests SHALL include a test verifying that an SSH connection with `accept_unknown_host_key` set to true succeeds when the host key is unknown, and that the host key is appended to the known_hosts file.

#### Scenario: Accept-unknown mode records host key
- **WHEN** the SSH test dials with a fresh, writable known_hosts file and `accept_unknown_host_key` set to true
- **THEN** the connection test succeeds and the known_hosts file contains a host key entry for the dialed host
