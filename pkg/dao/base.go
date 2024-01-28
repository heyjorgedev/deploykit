package dao

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/heyjorgedev/deploykit/pkg/model"
)

type Dao struct {
	db *dbx.DB
}

func New(db *dbx.DB) *Dao {
	return &Dao{db: db}
}

func (dao *Dao) DB() dbx.Builder {
	return dao.db
}

func (dao *Dao) ModelQuery(m model.Model) *dbx.SelectQuery {
	tableName := m.TableName()

	return dao.DB().
		Select("{{" + tableName + "}}.*").
		From(tableName)
}

func (dao *Dao) Close() error {
	return dao.db.Close()
}
