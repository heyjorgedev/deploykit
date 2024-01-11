package deploykit

import (
	"context"
	"strings"
)

type App struct {
	ID   uint16 `json:"id"`
	Name string `json:"name" validate:"required"`
}

func (a *App) Validate() error {
	// name cannot be empty
	if a.Name == "" {
		return Errorf(EINVALID, "name cannot be empty")
	}

	if strings.Contains(a.Name, " ") {
		return Errorf(EINVALID, "name cannot have spaces")
	}

	// name cannot be longer than 32 characters
	if len(a.Name) > 32 {
		return Errorf(EINVALID, "name cannot be longer than 32 characters")
	}

	return nil
}

type AppService interface {
	FindAll(ctx context.Context) ([]*App, error)
	Create(ctx context.Context, app *App) error
}
