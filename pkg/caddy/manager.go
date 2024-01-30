package caddy

import (
	"github.com/heyjorgedev/deploykit/pkg/core"
	"sync"
)

type Manager struct {
	app core.App
	mux sync.RWMutex
}

func NewManager(app core.App) *Manager {
	m := &Manager{app: app}

	app.OnAfterBootstrap().Add(func(event *core.BootstrapEvent) error {
		app.Logger().Info("Initializing Caddy Manager")
		return nil
	})

	app.OnTerminate().Add(func(event *core.TerminateEvent) error {
		app.Logger().Info("Terminating Caddy Manager")
		return nil
	})

	return m
}
