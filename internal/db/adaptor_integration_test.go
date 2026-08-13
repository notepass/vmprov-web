package db

import (
	"testing"

	"github.com/notepass/vmprov-web/internal/config"
	"github.com/stretchr/testify/require"
)

func TestAdaptorLifecycle_Postgres(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	adaptor, err := NewAdaptor(Config{
		ConnString:      "postgres://postgres:password@localhost:5432/postgres?sslmode=disable",
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 30,
	})
	require.NoError(t, err, "should connect to postgres")
	require.NotNil(t, adaptor)

	err = adaptor.HealthCheck()
	require.NoError(t, err, "health check should pass")

	err = adaptor.Close()
	require.NoError(t, err, "should close without error")
}

func TestAdaptorLifecycle_MySQL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	adaptor, err := NewAdaptor(Config{
		ConnString:   "root:password@tcp(localhost:3306)/testdb",
		MaxOpenConns: 5,
		MaxIdleConns: 2,
	})
	require.NoError(t, err, "should connect to mysql")
	require.NotNil(t, adaptor)

	err = adaptor.HealthCheck()
	require.NoError(t, err, "health check should pass")

	err = adaptor.Close()
	require.NoError(t, err, "should close without error")
}

func TestNewAdaptor_PoolConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg := &config.Config{
		DBConnString:      "postgres://postgres:password@localhost:5432/postgres?sslmode=disable",
		DBMaxOpenConns:    10,
		DBMaxIdleConns:    3,
		DBConnMaxLifetime: 60,
	}

	adaptor, err := NewAdaptor(Config{
		ConnString:      cfg.DBConnString,
		MaxOpenConns:    cfg.DBMaxOpenConns,
		MaxIdleConns:    cfg.DBMaxIdleConns,
		ConnMaxLifetime: cfg.DBConnMaxLifetime,
	})
	require.NoError(t, err)
	defer adaptor.Close()

	db := adaptor.DB()
	stats := db.Stats()
	require.Equal(t, 10, stats.MaxOpenConnections)
	require.True(t, stats.OpenConnections <= 10)
}
