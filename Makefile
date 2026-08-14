.PHONY: build run test integrate integration-down

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
