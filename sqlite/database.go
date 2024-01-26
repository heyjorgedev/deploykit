package sqlite

import "database/sql"

type Database struct {
	*sql.DB
}

func (db *Database) Open() error {
	return nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}
