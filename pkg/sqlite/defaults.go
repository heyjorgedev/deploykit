package sqlite

import (
	"database/sql"
	"fmt"
)

func EnableWAL(db *sql.DB) error {
	if _, err := db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("journal mode pragma: %w", err)
	}

	return nil
}

func EnableForeignKeys(db *sql.DB) error {
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	return nil
}
