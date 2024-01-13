package deploykit

import (
	"context"
	"strings"
)

type App struct {
	Name string `json:"name" validate:"required"`
}

func (a *App) Validate() error {
	err := a.ValidateName()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ValidateName() error {
	// name cannot be empty
	if a.Name == "" {
		return Errorf(EINVALID, "name cannot be empty")
	}

	// name cannot have spaces
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
