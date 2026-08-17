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

	_, err = db.Exec("DROP TABLE IF EXISTS audit_logs, templates, users, libvirt_connections, goose_db_version")
	require.NoError(t, err)

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
	require.Equal(t, int64(2), version, "should be at version 2 after up")

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

func TestMigration002_MySQL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/testdb")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS audit_logs, templates, users, libvirt_connections, goose_db_version")
	require.NoError(t, err)

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
	require.Equal(t, int64(2), version, "should be at version 2 after up")

	var count int
	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'testdb' AND table_name = 'libvirt_connections'").Scan(&count)
	require.Equal(t, 1, count, "libvirt_connections table should exist")

	err = goose.DownTo(db, "migrations", 1)
	require.NoError(t, err)

	version, err = goose.GetDBVersion(db)
	require.NoError(t, err)
	require.Equal(t, int64(1), version, "should be at version 1 after down")

	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'testdb' AND table_name = 'libvirt_connections'").Scan(&count)
	require.Equal(t, 0, count, "libvirt_connections table should be dropped")

	err = goose.DownTo(db, "migrations", 0)
	require.NoError(t, err)
}

func TestMigrationRollback_Postgres(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := sql.Open("pgx", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS audit_logs, templates, users, libvirt_connections, goose_db_version")
	require.NoError(t, err)

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
	require.Equal(t, int64(2), version, "should be at version 2 after up")

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

func TestMigration002_Postgres(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := sql.Open("pgx", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS audit_logs, templates, users, libvirt_connections, goose_db_version")
	require.NoError(t, err)

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
	require.Equal(t, int64(2), version, "should be at version 2 after up")

	var count int
	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'libvirt_connections'").Scan(&count)
	require.Equal(t, 1, count, "libvirt_connections table should exist")

	err = goose.DownTo(db, "migrations", 1)
	require.NoError(t, err)

	version, err = goose.GetDBVersion(db)
	require.NoError(t, err)
	require.Equal(t, int64(1), version, "should be at version 1 after down")

	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'libvirt_connections'").Scan(&count)
	require.Equal(t, 0, count, "libvirt_connections table should be dropped")

	err = goose.DownTo(db, "migrations", 0)
	require.NoError(t, err)
}
