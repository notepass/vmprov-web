## ADDED Requirements

### Requirement: Libvirt integration job in the CI workflow
The CI workflow SHALL include a `libvirt-integration` job that runs on the same ubuntu runner as the existing jobs and is executed only via the existing `workflow_dispatch` trigger. The job SHALL install `libvirt-daemon-system`, `libvirt-clients`, `libvirt-daemon-driver-qemu`, `qemu-kvm`, and `openssh-server` (none of which are preinstalled on the runner image), enable and start `libvirtd` and `ssh`, wait until `libvirtd` is active, and then run the gated libvirt integration tests. The job SHALL NOT require containerized databases.

#### Scenario: Job installs and starts libvirt and sshd
- **WHEN** the `libvirt-integration` job starts
- **THEN** the libvirt packages and `openssh-server` are installed, `libvirtd` and `ssh` are started, and the job waits for `libvirtd` to be active before running tests

#### Scenario: Gated tests execute in the job
- **WHEN** the `libvirt-integration` job reaches the test step
- **THEN** it runs the `internal/libvirt` integration tests with `LIBVIRT_INTEGRATION=1` and the job succeeds only if all tests pass

#### Scenario: Test failure fails the job
- **WHEN** any live integration test fails
- **THEN** the `libvirt-integration` job fails and reports the test error

#### Scenario: No automatic triggers
- **WHEN** a push or pull request event occurs
- **THEN** the `libvirt-integration` job does not run

### Requirement: Libvirt socket access for the runner user
The `libvirt-integration` job SHALL grant the runner user access to the libvirt socket by adding the user to the `libvirt` group and running the test command under that group, so that both direct socket dials and SSH-forwarded dials (which connect as the same user) can reach `/var/run/libvirt/libvirt-sock`.

#### Scenario: Tests run with libvirt group access
- **WHEN** the job runs the gated tests
- **THEN** the test command executes with the `libvirt` group effective, and no libvirt permission errors occur

### Requirement: SSH fixtures for the libvirt integration job
The `libvirt-integration` job SHALL generate an ed25519 SSH keypair for the runner user, authorize the public key in the runner user's `authorized_keys` (with `~/.ssh` at mode 0700 and `authorized_keys` at mode 0600), and build a known_hosts fixture for `127.0.0.1` using `ssh-keyscan` before running the tests.

#### Scenario: SSH key fixture is established
- **WHEN** the job sets up SSH fixtures
- **THEN** key-based authentication from the runner user to `127.0.0.1` succeeds using the generated key

#### Scenario: Known hosts fixture is populated
- **WHEN** the job builds the known_hosts fixture
- **THEN** the fixture file contains a host key entry for `127.0.0.1` usable by the strict-mode SSH test
