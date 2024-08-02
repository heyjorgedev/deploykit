package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	db     *sql.DB
	ctx    context.Context // background context
	cancel func()          // cancel background context

	DSN string
	Now func() time.Time
}

func NewDB(ctx context.Context, dsn string) *DB {
	db := &DB{
		DSN: dsn,
		Now: time.Now,
	}

	db.ctx, db.cancel = context.WithCancel(ctx)

	return db
}

func (db *DB) Open() (err error) {
	if db.DSN == "" {
		return fmt.Errorf("dsn required")
	}

	if db.DSN != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(db.DSN), 0700); err != nil {
			return err
		}
	}

	if db.db, err = sql.Open("sqlite3", db.DSN); err != nil {
		return err
	}

	// Enable WAL. SQLite performs better with the WAL  because it allows
	// multiple readers to operate while data is being written.
	if _, err := db.db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return fmt.Errorf("enable wal: %w", err)
	}

	// Enable foreign key checks. For historical reasons, SQLite does not check
	// foreign key constraints by default... which is kinda insane. There's some
	// overhead on inserts to verify foreign key integrity but it's definitely
	// worth it.
	if _, err := db.db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.db != nil {
		return db.db.Close()
	}

	return nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

func (db *DB) RunMigrations() error {
	// Create migrations table
	_, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Get the list of migration files
	files, err := fs.Glob(migrationsFS, "migrations/*.sql")

	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	sort.Strings(files)

	// Loop through the files and run each one
	for _, file := range files {
		if err := runMigration(db, file); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", file, err)
		}
	}

	return nil
}

func runMigration(db *DB, file string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	friendlyMigrationName := formatMigrationName(file)

	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, friendlyMigrationName).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil // already run migration, skip
	}

	// Run the migration
	if buf, err := fs.ReadFile(migrationsFS, file); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	// Record the migration
	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES (?)`, friendlyMigrationName); err != nil {
		return err
	}

	return tx.Commit()
}

type Tx struct {
	*sql.Tx
	db  *DB
	now time.Time
}

func formatMigrationName(file string) string {
	parts := strings.Split(file, "/")
	fileName := parts[len(parts)-1]

	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
