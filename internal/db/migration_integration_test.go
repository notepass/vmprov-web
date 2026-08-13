package db

import (
	"database/sql"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

func TestMigrationRollback_MySQL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/testdb")
	require.NoError(t, err)
	defer db.Close()

	db.Exec("DROP TABLE IF EXISTS audit_logs, templates, users, goose_db_version")

	goose.SetDialect("mysql")
	goose.SetBaseFS(&dialectFS{
		embedded: embeddedMigrations,
		dialect:  "mysql",
		prefix:   "migrations/",
	})

	err = goose.Up(db, "migrations")
	require.NoError(t, err)

	version, err := goose.GetDBVersion(db)
	require.NoError(t, err)
	require.Equal(t, int64(1), version, "should be at version 1 after up")

	var count int
	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'testdb' AND table_name = 'users'").Scan(&count)
	require.Equal(t, 1, count, "users table should exist")

	err = goose.DownTo(db, "migrations", 0)
	require.NoError(t, err)

	version, err = goose.GetDBVersion(db)
	require.NoError(t, err)
	require.Equal(t, int64(0), version, "should be at version 0 after down")

	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'testdb' AND table_name = 'users'").Scan(&count)
	require.Equal(t, 0, count, "users table should be dropped")
}

func TestMigrationRollback_Postgres(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := sql.Open("pgx", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	db.Exec("DROP TABLE IF EXISTS audit_logs, templates, users")
	db.Exec("DROP TABLE IF EXISTS goose_db_version")

	goose.SetDialect("postgres")
	goose.SetBaseFS(&dialectFS{
		embedded: embeddedMigrations,
		dialect:  "postgres",
		prefix:   "migrations/",
	})

	err = goose.Up(db, "migrations")
	require.NoError(t, err)

	version, err := goose.GetDBVersion(db)
	require.NoError(t, err)
	require.Equal(t, int64(1), version, "should be at version 1 after up")

	var count int
	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users'").Scan(&count)
	require.Equal(t, 1, count, "users table should exist")

	err = goose.DownTo(db, "migrations", 0)
	require.NoError(t, err)

	version, err = goose.GetDBVersion(db)
	require.NoError(t, err)
	require.Equal(t, int64(0), version, "should be at version 0 after down")

	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users'").Scan(&count)
	require.Equal(t, 0, count, "users table should be dropped")
}
