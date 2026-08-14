FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server ./cmd/server

FROM alpine:3.22

RUN adduser -D -u 10001 vmprov
WORKDIR /app
COPY --from=builder /app/server /app/server
USER vmprov

EXPOSE 8080
ENTRYPOINT ["/app/server"]
