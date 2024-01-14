package sqlite

import (
	"context"

	"github.com/heyjorgedev/deploykit"
)

func createApp(ctx context.Context, tx *Tx, app *deploykit.App) error {
	row := tx.QueryRowContext(ctx, "INSERT INTO apps (name) VALUES (?) RETURNING name", app.Name)
	if row.Err() != nil {
		return FormatError(row.Err())
	}

	err := row.Scan(&app.Name)
	if err != nil {
		return FormatError(err)
	}

	return nil
}
