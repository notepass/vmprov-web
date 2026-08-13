package db

import (
	"fmt"
	"io/fs"
	"testing"
)

func TestEmbeddedMigrations(t *testing.T) {
	entries, err := fs.ReadDir(embeddedMigrations, "migrations")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Embedded migrations:\n")
	for _, e := range entries {
		fmt.Printf("  %s\n", e.Name())
	}
}

func TestDialectFS(t *testing.T) {
	dfs := &dialectFS{
		embedded: embeddedMigrations,
		dialect:  "postgres",
		prefix:   "migrations/",
	}
	
	entries, err := dfs.ReadDir("migrations")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("PostgreSQL migrations:\n")
	for _, e := range entries {
		fmt.Printf("  %s\n", e.Name())
	}
}
