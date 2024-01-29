package caddy

import "github.com/heyjorgedev/deploykit/pkg/core"

type Manager struct {
	app core.App
}

func NewManager(app core.App) *Manager {
	m := &Manager{app: app}

	app.OnAfterBootstrap().Add(func(event *core.BootstrapEvent) error {
		app.Logger().Info("Initializing Caddy Manager")
		return nil
	})

	app.OnTerminate().Add(func(event *core.TerminateEvent) error {
		app.Logger().Info("Terminating Caddy Manager")
		return m.Close()
	})

	return m
}

func (m *Manager) Close() error {
	return nil
}
