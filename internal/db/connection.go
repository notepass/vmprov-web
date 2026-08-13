package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Config holds database connection configuration.
type Config struct {
	ConnString      string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

// postgresAdaptor wraps a postgresql database connection.
type postgresAdaptor struct {
	db *sql.DB
}

// mysqlAdaptor wraps a mysql/mariadb database connection.
type mysqlAdaptor struct {
	db *sql.DB
}

// NewAdaptor creates a new database adaptor.
func NewAdaptor(cfg Config) (Adaptor, error) {
	if cfg.ConnString == "" {
		return nil, fmt.Errorf("connection string is required")
	}

	drv, connStr := resolveDriver(cfg.ConnString)
	db, err := sql.Open(drv, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	switch drv {
	case "pgx":
		return &postgresAdaptor{db: db}, nil
	case "mysql":
		return &mysqlAdaptor{db: db}, nil
	default:
		db.Close()
		return nil, fmt.Errorf("unsupported database driver: %s", drv)
	}
}

// resolveDriver determines the driver and connection string from the provided connection string.
func resolveDriver(connStr string) (string, string) {
	if len(connStr) >= 11 && connStr[:11] == "postgres://" {
		return "pgx", connStr
	}
	return "mysql", connStr
}

func (a *postgresAdaptor) DB() *sql.DB {
	return a.db
}

func (a *postgresAdaptor) HealthCheck() error {
	return a.db.Ping()
}

func (a *postgresAdaptor) Close() error {
	return a.db.Close()
}

func (a *mysqlAdaptor) DB() *sql.DB {
	return a.db
}

func (a *mysqlAdaptor) HealthCheck() error {
	return a.db.Ping()
}

func (a *mysqlAdaptor) Close() error {
	return a.db.Close()
}
