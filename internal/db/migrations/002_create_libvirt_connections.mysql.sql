-- +goose Up
CREATE TABLE libvirt_connections (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL UNIQUE,
    `type` VARCHAR(20) NOT NULL,
    `host` VARCHAR(255),
    `username` VARCHAR(255),
    ssh_key_path VARCHAR(1024),
    accept_unknown_host_key BOOLEAN NOT NULL DEFAULT FALSE,
    socket_path VARCHAR(1024),
    description TEXT,
    last_status VARCHAR(20),
    last_checked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_libvirt_connections_type CHECK (`type` IN ('ssh', 'socket'))
);

-- +goose Down
DROP TABLE IF EXISTS libvirt_connections;
