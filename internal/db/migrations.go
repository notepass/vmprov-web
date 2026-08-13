package db

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/pressly/goose/v3"
)

// RunMigrations applies all pending migrations to the database.
func RunMigrations(a Adaptor) error {
	dialect := getDialect(a)

	// Collect migration files for this dialect
	migrations, err := collectMigrationsForDialect(dialect)
	if err != nil {
		return fmt.Errorf("failed to collect migrations: %w", err)
	}

	if len(migrations) == 0 {
		return fmt.Errorf("no migrations found for dialect %s", dialect)
	}

	// Create an in-memory FS with only the relevant migrations
	// We'll use goose with the collected migration paths
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Create a temporary embedded FS with only the relevant migrations
	migrationFS := &dialectFS{
		embedded: embeddedMigrations,
		dialect:  dialect,
		prefix:   "migrations/",
	}

	goose.SetBaseFS(migrationFS)

	if err := goose.Up(a.DB(), "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// dialectFS filters embedded migrations by dialect.
type dialectFS struct {
	embedded fs.FS
	dialect  string
	prefix   string
}

func (f *dialectFS) Open(name string) (fs.File, error) {
	return f.embedded.Open(name)
}

func (f *dialectFS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := fs.ReadDir(f.embedded, name)
	if err != nil {
		return nil, err
	}

	var filtered []fs.DirEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.IsDir() {
			filtered = append(filtered, entry)
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, "."+f.dialect+".sql") {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

// collectMigrationsForDialect lists migration files for the given dialect.
func collectMigrationsForDialect(dialect string) ([]string, error) {
	var migrations []string
	err := fs.WalkDir(embeddedMigrations, "migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if strings.HasSuffix(name, "."+dialect+".sql") {
			migrations = append(migrations, path)
		}
		return nil
	})
	return migrations, err
}

func getDialect(a Adaptor) string {
	switch a.(type) {
	case *postgresAdaptor:
		return "postgres"
	case *mysqlAdaptor:
		return "mysql"
	default:
		return "postgres"
	}
}
