-- +goose Up
CREATE TABLE libvirt_connections (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('ssh', 'socket')),
    host VARCHAR(255),
    username VARCHAR(255),
    ssh_key_path VARCHAR(1024),
    accept_unknown_host_key BOOLEAN NOT NULL DEFAULT FALSE,
    socket_path VARCHAR(1024),
    description TEXT,
    last_status VARCHAR(20),
    last_checked_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS libvirt_connections;
