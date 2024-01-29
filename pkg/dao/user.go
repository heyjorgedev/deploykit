package dao

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/heyjorgedev/deploykit/pkg/model"
)

func (dao *Dao) UserQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&model.User{})
}

func (dao *Dao) FindUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	if err := dao.UserQuery().Where(dbx.HashExp{"username": username}).One(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (dao *Dao) FindUserById(id int) (*model.User, error) {
	user := &model.User{}
	if err := dao.UserQuery().Where(dbx.HashExp{"id": id}).One(user); err != nil {
		return nil, err
	}

	return user, nil
}
