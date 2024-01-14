package sqlite

import (
	"context"

	"github.com/heyjorgedev/deploykit"
)

type AppService struct {
	db *DB
}

func NewAppService(db *DB) *AppService {
	return &AppService{db}
}

func (s *AppService) FindAll(ctx context.Context) ([]*deploykit.App, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "SELECT name FROM apps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := make([]*deploykit.App, 0)
	for rows.Next() {
		var app deploykit.App
		if err := rows.Scan(&app.Name); err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}

	return apps, nil
}

func (s *AppService) Create(ctx context.Context, app *deploykit.App) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return FormatError(err)
	}
	defer tx.Rollback()

	if err := createApp(ctx, tx, app); err != nil {
		return err
	}

	return tx.Commit()
}
