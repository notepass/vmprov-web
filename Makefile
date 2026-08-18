.PHONY: build run test integrate integration-down libvirt-integrate

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

test:
	go test ./...

integrate:
	docker compose up -d --wait
	go test ./...
	status=$$?
	docker compose down
	exit $$status

integration-down:
	docker compose down

libvirt-integrate:
	LIBVIRT_INTEGRATION=1 go test ./internal/libvirt/... -count=1
