## Purpose

Provide a manually-triggered GitHub Actions CI pipeline that runs unit and integration tests (against containerized databases), builds target-platform binaries, and produces a working Docker image.

## Requirements

### Requirement: CI workflow triggers manually
The project SHALL provide a GitHub Actions workflow that is triggered only by `workflow_dispatch`.

#### Scenario: Workflow dispatch available
- **WHEN** a user triggers the workflow via GitHub Actions UI or `gh workflow run`
- **THEN** the workflow starts and executes all configured jobs

#### Scenario: No automatic triggers
- **WHEN** a push or pull request event occurs
- **THEN** the CI workflow does not run

### Requirement: Unit tests execute successfully
The workflow SHALL run `go test ./... -short` to execute all unit tests.

#### Scenario: Unit tests pass
- **WHEN** unit tests are executed in the workflow
- **THEN** all tests in packages without database dependencies pass

#### Scenario: Unit tests fail
- **WHEN** a unit test fails
- **THEN** the workflow job fails and reports the error

### Requirement: Integration tests run with containerized databases
The workflow SHALL start PostgreSQL 18 and MariaDB 12 containers before running integration tests via `go test ./...`.

#### Scenario: PostgreSQL integration tests pass
- **WHEN** the PostgreSQL 18 container is running with `POSTGRES_PASSWORD=password`
- **THEN** integration tests connecting to `postgres://postgres:password@localhost:5432/postgres?sslmode=disable` succeed

#### Scenario: MariaDB integration tests pass
- **WHEN** the MariaDB 12 container is running with `MYSQL_ROOT_PASSWORD=password` and `MYSQL_DATABASE=testdb`
- **THEN** integration tests connecting to `root:password@tcp(localhost:3306)/testdb` succeed

#### Scenario: Database startup failure
- **WHEN** a database container fails to start
- **THEN** the workflow job fails with a clear error message

### Requirement: Binaries build for target platforms
The workflow SHALL build the Go binary for linux/amd64 and linux/arm64.

#### Scenario: linux/amd64 binary builds
- **WHEN** the build job runs with `GOOS=linux GOARCH=amd64`
- **THEN** a functional binary is produced at `bin/server`

#### Scenario: linux/arm64 binary builds
- **WHEN** the build job runs with `GOOS=linux GOARCH=arm64`
- **THEN** a functional binary is produced at `bin/server`

#### Scenario: Build artifacts uploaded
- **WHEN** binaries are built successfully
- **THEN** they are uploaded as workflow artifacts for download

### Requirement: Dockerfile produces working image
The project SHALL include a Dockerfile that builds the Go application into a minimal container image.

#### Scenario: Dockerfile builds
- **WHEN** `docker build -t vmprov-web .` is run
- **THEN** a Docker image is produced containing the Go binary

#### Scenario: Multi-stage build
- **WHEN** the Dockerfile is inspected
- **THEN** it uses a multi-stage build with a Go builder stage and a minimal runtime stage

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
