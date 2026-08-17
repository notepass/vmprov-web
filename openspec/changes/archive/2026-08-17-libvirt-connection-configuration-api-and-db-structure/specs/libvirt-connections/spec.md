## ADDED Requirements

### Requirement: Libvirt connection persistence
The system SHALL persist libvirt connection configurations in a `libvirt_connections` table with columns `id`, `name` (unique), `type` (CHECK: `ssh` or `socket`), `host` (nullable), `username` (nullable), `ssh_key_path` (nullable), `accept_unknown_host_key` (boolean, default false), `socket_path` (nullable), `description` (nullable), `last_status` (nullable), `last_checked_at` (nullable), `created_at`, and `updated_at`, applied via goose migration `002` for both the PostgreSQL and MySQL dialects.

#### Scenario: Table exists after startup on PostgreSQL
- **WHEN** the application starts against a PostgreSQL database
- **THEN** migration `002` is applied and the `libvirt_connections` table exists with a unique constraint on `name` and a check constraint on `type`

#### Scenario: Table exists after startup on MySQL
- **WHEN** the application starts against a MySQL or MariaDB database
- **THEN** migration `002` is applied and the `libvirt_connections` table exists with a unique constraint on `name` and a check constraint on `type`

#### Scenario: Migration rollback
- **WHEN** a rollback to the migration prior to `002` is requested
- **THEN** the down script removes the `libvirt_connections` table

### Requirement: Libvirt connection domain model
The domain package SHALL define a `LibvirtConnection` struct with `ID`, `Name`, `Type`, `Host`, `Username`, `SSHKeyPath`, `AcceptUnknownHostKey`, `SocketPath`, `Description`, `LastStatus`, `LastCheckedAt`, `CreatedAt`, and `UpdatedAt` fields mapped to the `libvirt_connections` table columns.

#### Scenario: Model maps to table columns
- **WHEN** a row in `libvirt_connections` is read through the repository
- **THEN** it is returned as a `LibvirtConnection` with all columns populated, using nullable types for `Host`, `Username`, `SSHKeyPath`, `SocketPath`, `Description`, `LastStatus`, and `LastCheckedAt`

### Requirement: Libvirt connection repository
The repository package SHALL define a `LibvirtConnectionRepository` interface with `Create`, `GetByID`, `GetByName`, `Update`, `Delete`, and `List` methods, implemented with Squirrel queries compatible with both supported database dialects.

#### Scenario: Repository interface exists
- **WHEN** the `repository` package is imported
- **THEN** a `LibvirtConnectionRepository` interface is available

#### Scenario: Repository is testable with mocks
- **WHEN** a unit test needs to test logic depending on libvirt connection storage
- **THEN** a mock implementation of `LibvirtConnectionRepository` can be substituted

#### Scenario: Cross-database queries
- **WHEN** a repository method executes a query against PostgreSQL or MySQL
- **THEN** Squirrel generates placeholder syntax valid for the target dialect

### Requirement: Route registration
The Echo server SHALL register routes for the libvirt connection endpoints under `/api/v1/remotes/libvirt/connections`.

#### Scenario: Routes are registered
- **WHEN** the application starts
- **THEN** the routes for list, create, get, update, delete, and test libvirt connections are reachable via HTTP

### Requirement: List libvirt connections
`GET /api/v1/remotes/libvirt/connections` SHALL return all stored libvirt connections as a JSON array.

#### Scenario: List returns stored connections
- **WHEN** connections exist in the database and the client sends `GET /api/v1/remotes/libvirt/connections`
- **THEN** the response is 200 with a JSON array containing each connection's `id`, `name`, `type`, `host`, `username`, `ssh_key_path`, `accept_unknown_host_key`, `socket_path`, `description`, `last_status`, `last_checked_at`, `created_at`, and `updated_at`

#### Scenario: List returns empty array
- **WHEN** no connections are stored and the client sends `GET /api/v1/remotes/libvirt/connections`
- **THEN** the response is 200 with an empty JSON array

### Requirement: Create libvirt connection
`POST /api/v1/remotes/libvirt/connections` SHALL create a connection from a JSON body containing `name`, `type`, the type-specific fields, and an optional `description`. SSH connections require `host`, `username`, and `ssh_key_path`; socket connections require `socket_path`. The `type` must be `ssh` or `socket`, `name` must be unique, and for SSH connections the key file must be valid at save time.

#### Scenario: Successful SSH connection creation
- **WHEN** the client sends a create request with `type` `ssh`, a unique `name`, `host`, `username`, and a valid `ssh_key_path`
- **THEN** the response is 201 with the created connection including its assigned `id`

#### Scenario: Successful socket connection creation
- **WHEN** the client sends a create request with `type` `socket`, a unique `name`, and a `socket_path`
- **THEN** the response is 201 with the created connection including its assigned `id`

#### Scenario: Missing required field
- **WHEN** the client sends a create request missing `name`, `type`, or a field required by the given type
- **THEN** the response is 400 with an error describing the missing field

#### Scenario: Invalid connection type
- **WHEN** the client sends a create request with a `type` other than `ssh` or `socket`
- **THEN** the response is 400 with an error describing the invalid type

#### Scenario: Missing or unreadable SSH key
- **WHEN** the client sends an SSH create request whose `ssh_key_path` does not exist or is not readable by the application
- **THEN** the response is 400 with an error indicating the key file problem

#### Scenario: Passphrase-protected SSH key
- **WHEN** the client sends an SSH create request whose `ssh_key_path` is a passphrase-protected private key
- **THEN** the response is 400 with an error indicating that passphrase-protected keys are not supported

#### Scenario: Duplicate name
- **WHEN** the client sends a create request with a `name` that already exists
- **THEN** the response is 409 with an error indicating the name is already in use

### Requirement: Get libvirt connection
`GET /api/v1/remotes/libvirt/connections/:id` SHALL return a single connection by ID.

#### Scenario: Connection found
- **WHEN** the client requests an existing connection ID
- **THEN** the response is 200 with the connection JSON

#### Scenario: Connection not found
- **WHEN** the client requests a non-existent connection ID
- **THEN** the response is 404

#### Scenario: Invalid ID
- **WHEN** the client requests a non-numeric `:id`
- **THEN** the response is 400

### Requirement: Update libvirt connection
`PUT /api/v1/remotes/libvirt/connections/:id` SHALL update the `name`, `type`, type-specific fields (`host`, `username`, `ssh_key_path`, `accept_unknown_host_key`, `socket_path`), and `description` of an existing connection, refreshing `updated_at`. The same per-type validation and SSH key validation as creation SHALL apply.

#### Scenario: Successful update
- **WHEN** the client sends a valid update for an existing connection
- **THEN** the response is 200 with the updated connection

#### Scenario: Update of missing connection
- **WHEN** the client sends an update for a non-existent connection ID
- **THEN** the response is 404

#### Scenario: Update to duplicate name
- **WHEN** the client sends an update whose `name` matches another existing connection
- **THEN** the response is 409

#### Scenario: Update with invalid SSH key
- **WHEN** the client sends an update for an SSH connection whose `ssh_key_path` is missing, unreadable, or passphrase-protected
- **THEN** the response is 400 with an error describing the key problem

### Requirement: Delete libvirt connection
`DELETE /api/v1/remotes/libvirt/connections/:id` SHALL remove a connection by ID.

#### Scenario: Successful deletion
- **WHEN** the client deletes an existing connection
- **THEN** the response is 204 and the connection no longer appears in list results

#### Scenario: Deletion of missing connection
- **WHEN** the client deletes a non-existent connection ID
- **THEN** the response is 404

### Requirement: Test libvirt connection
`POST /api/v1/remotes/libvirt/connections/:id/test` SHALL dial the stored connection through the libvirt client using the configured connect timeout — via the SSH dialer (stored host, username, key file, and host-key verification against the configured known_hosts file) for `ssh` connections, or via the local socket dialer (stored socket path) for `socket` connections — and report the libvirt version, hypervisor type, and total and active domain counts. It SHALL persist `last_status` and `last_checked_at` on the connection after each check.

#### Scenario: Reachable libvirt host
- **WHEN** the client tests a connection whose libvirt host is reachable
- **THEN** the response is 200 with `status` set to `ok`, the libvirt version, the hypervisor type, and domain counts, and the stored `last_status` is updated to `ok`

#### Scenario: Unreachable libvirt host
- **WHEN** the client tests a connection whose libvirt host is unreachable
- **THEN** the response is 502 with `status` set to `error` and a message, and the stored `last_status` is updated to `error`

#### Scenario: Dial timeout
- **WHEN** the dial operation exceeds the configured connect timeout
- **THEN** the response is 502 with an error indicating the timeout

#### Scenario: Unknown host key in strict mode
- **WHEN** the client tests an SSH connection whose host key is not in the known_hosts file and `accept_unknown_host_key` is false
- **THEN** the response is 502 with an error indicating the host key could not be verified

#### Scenario: Unknown host key with accept-unknown enabled
- **WHEN** the client tests an SSH connection whose host key is not in the known_hosts file and `accept_unknown_host_key` is true
- **THEN** the host key is appended to the known_hosts file and the response is 200 with `status` set to `ok`

#### Scenario: Test of missing connection
- **WHEN** the client tests a non-existent connection ID
- **THEN** the response is 404

### Requirement: Libvirt client abstraction
The system SHALL define a libvirt client interface in `internal/libvirt` that abstracts connection testing for a typed connection (type plus per-type fields, known_hosts file path, and timeout), with a concrete implementation using `github.com/digitalocean/go-libvirt` dialers (`dialers.NewSSH` and `dialers.NewLocal` via `libvirt.NewWithDialer`). Each test SHALL establish a new connection and disconnect on completion.

#### Scenario: Client interface exists
- **WHEN** the `libvirt` package is imported
- **THEN** a client interface is available that returns libvirt version, hypervisor type, and domain counts for a typed connection

#### Scenario: Client is mockable
- **WHEN** a unit test needs to test handler behavior without a live libvirt host
- **THEN** a fake implementation of the client interface can be substituted

#### Scenario: Connection is closed after use
- **WHEN** a connection test completes or fails
- **THEN** the underlying libvirt connection is disconnected
