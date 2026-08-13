package db

import "database/sql"

// Adaptor defines the interface for database operations.
type Adaptor interface {
	DB() *sql.DB
	HealthCheck() error
	Close() error
}
