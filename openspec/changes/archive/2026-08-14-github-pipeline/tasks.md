## 1. Dockerfile

- [x] 1.1 Create multi-stage Dockerfile with golang:1.26 builder stage
- [x] 1.2 Configure Alpine runtime stage with CGO disabled and minimal dependencies
- [x] 1.3 Verify Dockerfile builds successfully with `docker build -t vmprov-web .`

## 2. Docker Compose

- [x] 2.1 Create docker-compose.yml with PostgreSQL 18 service (POSTGRES_PASSWORD=password, healthcheck on port 5432)
- [x] 2.2 Add MariaDB 12 service (MYSQL_ROOT_PASSWORD=password, MYSQL_DATABASE=testdb, healthcheck on port 3306)
- [x] 2.3 Verify services start correctly with `docker-compose up -d` and pass healthchecks

## 3. Makefile integration target

- [x] 3.1 Add `integrate` target that starts docker-compose services, runs `go test ./...`, then tears down services
- [x] 3.2 Add `integration-down` target for manual cleanup
- [x] 3.3 Verify `make integrate` completes successfully locally

## 4. GitHub Actions workflow

- [x] 4.1 Create `.github/workflows/ci.yml` with `workflow_dispatch` trigger
- [x] 4.2 Add unit test job: checkout, setup Go, run `go test ./... -short`
- [x] 4.3 Add integration test job: checkout, setup Go, start docker-compose services, wait for health, run `go test ./...`
- [x] 4.4 Add build job: checkout, setup Go, cross-compile linux/amd64 and linux/arm64, upload artifacts
- [x] 4.5 Verify workflow triggers correctly via manual dispatch
