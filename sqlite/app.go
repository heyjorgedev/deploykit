package sqlite

import (
	"context"

	"github.com/jorgemurta/deploykit"
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

	rows, err := tx.QueryContext(ctx, "SELECT id, name FROM apps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := make([]*deploykit.App, 0)
	for rows.Next() {
		var app deploykit.App
		if err := rows.Scan(&app.ID, &app.Name); err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}

	return apps, nil
}

func (s *AppService) Create(ctx context.Context, app *deploykit.App) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, "INSERT INTO apps (name) VALUES (?) RETURNING id, name", app.Name)
	if row.Err() != nil {
		return err
	}

	err = row.Scan(&app.ID, &app.Name)
	if err != nil {
		return err
	}

	return tx.Commit()
}
