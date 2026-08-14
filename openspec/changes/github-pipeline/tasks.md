## 1. Dockerfile

- [ ] 1.1 Create multi-stage Dockerfile with golang:1.26 builder stage
- [ ] 1.2 Configure Alpine runtime stage with CGO disabled and minimal dependencies
- [ ] 1.3 Verify Dockerfile builds successfully with `docker build -t vmprov-web .`

## 2. Docker Compose

- [ ] 2.1 Create docker-compose.yml with PostgreSQL 18 service (POSTGRES_PASSWORD=password, healthcheck on port 5432)
- [ ] 2.2 Add MariaDB 12 service (MYSQL_ROOT_PASSWORD=password, MYSQL_DATABASE=testdb, healthcheck on port 3306)
- [ ] 2.3 Verify services start correctly with `docker-compose up -d` and pass healthchecks

## 3. Makefile integration target

- [ ] 3.1 Add `integrate` target that starts docker-compose services, runs `go test ./...`, then tears down services
- [ ] 3.2 Add `integration-down` target for manual cleanup
- [ ] 3.3 Verify `make integrate` completes successfully locally

## 4. GitHub Actions workflow

- [ ] 4.1 Create `.github/workflows/ci.yml` with `workflow_dispatch` trigger
- [ ] 4.2 Add unit test job: checkout, setup Go, run `go test ./... -short`
- [ ] 4.3 Add integration test job: checkout, setup Go, start docker-compose services, wait for health, run `go test ./...`
- [ ] 4.4 Add build job: checkout, setup Go, cross-compile linux/amd64 and linux/arm64, upload artifacts
- [ ] 4.5 Verify workflow triggers correctly via manual dispatch
