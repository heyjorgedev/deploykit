package sqlite

import (
	"context"
	"testing"
)

func NewTestingDB(t *testing.T) (*DB, func() error) {
	db := NewDB(context.Background(), ":memory:")
	err := db.Open()
	if err != nil {
		t.Fatalf("error opening db: %v", err)
	}

	return db, func() error {
		return db.Close()
	}
}

func NewMigratedTestingDB(t *testing.T) (*DB, func() error) {
	db, cleanup := NewTestingDB(t)

	err := db.RunMigrations()
	if err != nil {
		t.Fatalf("error running migrations: %v", err)
	}

	return db, func() error {
		return cleanup()
	}
}

func TestDB(t *testing.T) {
	t.Run("RunMigrations", func(t *testing.T) {
		db, cleanup := NewTestingDB(t)
		defer cleanup()

		err := db.RunMigrations()
		if err != nil {
			t.Fatalf("error running migrations: %v", err)
		}
	})
}

func TestFormatMigrationName(t *testing.T) {
	t.Run("FormatWithPath", func(t *testing.T) {
		result := formatMigrationName("migrations/000_initial.sql")

		if result != "000_initial" {
			t.Fatalf("expected 000_initial, got %s", result)
		}
	})

	t.Run("FormatWithoutPath", func(t *testing.T) {
		result := formatMigrationName("000_initial.sql")

		if result != "000_initial" {
			t.Fatalf("expected 000_initial, got %s", result)
		}
	})
}
