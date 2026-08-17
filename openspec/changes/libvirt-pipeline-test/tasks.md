## 1. Gated live integration tests

- [ ] 1.1 Create `internal/libvirt/libvirt_integration_test.go` with an env-var helper (skip all tests unless `LIBVIRT_INTEGRATION` is non-empty; read `LIBVIRT_SSH_HOST` default `127.0.0.1`, `LIBVIRT_SSH_USER` default current user, `LIBVIRT_SSH_KEY`, `LIBVIRT_KNOWN_HOSTS`, `LIBVIRT_SOCKET_PATH` default `/var/run/libvirt/libvirt-sock`, `LIBVIRT_TIMEOUT` default 15s)
- [ ] 1.2 Add the local socket test: dial `LIBVIRT_SOCKET_PATH` with the concrete client, assert version matches `^\d+\.\d+\.\d+$`, hypervisor type `QEMU`, `TotalDomains >= ActiveDomains >= 0`
- [ ] 1.3 Add the SSH strict test: dial `LIBVIRT_SSH_HOST` as `LIBVIRT_SSH_USER` with `LIBVIRT_SSH_KEY` and `LIBVIRT_KNOWN_HOSTS`, assert success and version/hypervisor type
- [ ] 1.4 Add the SSH unknown-host-key rejection test: fresh empty known_hosts file, `AcceptUnknownHostKey=false`, expect an error indicating the host key could not be verified
- [ ] 1.5 Add the SSH accept-unknown test: fresh empty writable known_hosts file in `t.TempDir()`, `AcceptUnknownHostKey=true`, expect success and assert the file now contains a host key entry for the dialed host
- [ ] 1.6 Verify `gofmt`, `go vet`, and `go test ./... -short` stay green (live tests skip by default)

## 2. Local ergonomics

- [ ] 2.1 Add a `make libvirt-integrate` target that runs the gated tests with `LIBVIRT_INTEGRATION=1` and the default socket path
- [ ] 2.2 Run `make libvirt-integrate` locally to confirm the socket test passes (skip SSH tests locally if sshd/key setup is unavailable, noting what was covered)

## 3. CI job

- [ ] 3.1 Add a `libvirt-integration` job to `.github/workflows/ci.yml` (ubuntu-latest, no `needs`, no database steps): checkout, setup-go (1.26), apt-install `libvirt-daemon-system`, `libvirt-clients`, `libvirt-daemon-driver-qemu`, `qemu-kvm`, `openssh-server`, `systemctl enable --now libvirtd ssh`, wait for `libvirtd` active
- [ ] 3.2 Add libvirt-group access step: `sudo usermod -aG libvirt "$(whoami)"`, and run the test command via `sg libvirt -c '...'`
- [ ] 3.3 Add SSH fixture step: generate an ed25519 keypair, authorize the public key (`~/.ssh` 0700, `authorized_keys` 0600), build the known_hosts fixture with `ssh-keyscan 127.0.0.1`
- [ ] 3.4 Run the gated tests with `LIBVIRT_INTEGRATION=1` plus `LIBVIRT_SSH_KEY`/`LIBVIRT_KNOWN_HOSTS` env vars, `-count=1 -v`

## 4. Verification

- [ ] 4.1 `go build ./...`, `go vet ./...`, `go test ./... -short` all green
- [ ] 4.2 `make libvirt-integrate` green locally (socket test)
- [ ] 4.3 Trigger the workflow (`gh workflow run CI`) and confirm the `libvirt-integration` job passes, including all four live tests
- [ ] 4.4 Note in the archived change's `tasks.md` that task 9.4's SSH-against-a-live-host portion is now covered by this change's CI job
