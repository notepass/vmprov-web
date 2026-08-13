package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdaptor_EmptyConnString(t *testing.T) {
	_, err := NewAdaptor(Config{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection string is required")
}

func TestNewAdaptor_InvalidConnection(t *testing.T) {
	_, err := NewAdaptor(Config{
		ConnString:   "postgres://localhost:9999/nonexistent",
		MaxOpenConns: 25,
		MaxIdleConns: 5,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed")
}

func TestResolveDriver_Postgres(t *testing.T) {
	drv, connStr := resolveDriver("postgres://localhost:5432/db")
	assert.Equal(t, "pgx", drv)
	assert.Equal(t, "postgres://localhost:5432/db", connStr)
}

func TestResolveDriver_MySQL(t *testing.T) {
	drv, connStr := resolveDriver("user:pass@tcp(localhost:3306)/db")
	assert.Equal(t, "mysql", drv)
	assert.Equal(t, "user:pass@tcp(localhost:3306)/db", connStr)
}
