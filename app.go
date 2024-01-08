package deploykit

import (
	"context"
)

type App struct {
	ID   uint16 `json:"id"`
	Name string `json:"name" validate:"required"`
}

func (a *App) Validate() error {
	return nil
}

type AppService interface {
	FindAll(ctx context.Context) ([]*App, error)
	Create(ctx context.Context, app *App) error
}
