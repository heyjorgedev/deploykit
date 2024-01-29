package dao

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/heyjorgedev/deploykit/pkg/model"
)

type Dao struct {
	db dbx.Builder
}

func New(db dbx.Builder) *Dao {
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

func (dao *Dao) InTransaction(fn func(dao *Dao) error) error {
	switch txOrDb := dao.DB().(type) {

	case *dbx.Tx:
		txDao := New(txOrDb)
		return fn(txDao)
	case *dbx.DB:
		return txOrDb.Transactional(func(tx *dbx.Tx) error {
			txDao := New(tx)
			return fn(txDao)
		})
	}

	return nil
}
