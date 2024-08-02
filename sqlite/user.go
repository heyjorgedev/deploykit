package sqlite

import (
	"context"

	"github.com/heyjorgedev/deploykit"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, user *deploykit.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createUser(ctx, tx, user); err != nil {
		return err
	}

	return tx.Commit()
}

func createUser(ctx context.Context, tx *Tx, user *deploykit.User) error {
	result, err := tx.ExecContext(ctx, "INSERT INTO users (name, email) values (?, ?);", user.Name, user.Email)
	if err != nil {
		return formatError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return formatError(err)
	}

	user.ID = int(id)

	return nil
}
