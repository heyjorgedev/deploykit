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

	rows, err := tx.QueryContext(ctx, "SELECT id, name, network FROM apps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := make([]*deploykit.App, 0)
	for rows.Next() {
		var app deploykit.App
		if err := rows.Scan(&app.ID, &app.Name, &app.Network); err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}

	return apps, nil
}
