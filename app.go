package deploykit

import "context"

type App struct {
	ID      uint16 `json:"id"`
	Name    string `json:"name"`
	Network string `json:"network"`
}

func (a *App) Validate() error {
	return nil
}

type AppService interface {
	FindAll(ctx context.Context) ([]*App, error)
}
