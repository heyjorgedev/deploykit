package core

import (
	"database/sql"
	docker "github.com/docker/docker/client"
	"github.com/heyjorgedev/deploykit/pkg/core/hook"
	"github.com/heyjorgedev/deploykit/pkg/dao"
	"log/slog"
)

type App interface {
	// DB returns the database connection.
	DB() *sql.DB

	Dao() *dao.Dao

	// HostDocker returns the docker connection of the same host.
	HostDocker() *docker.Client

	// Logger returns the logger.
	Logger() *slog.Logger

	// DataDir returns the data directory path.
	DataDir() string

	// Bootstrap initializes the app, e.g. by opening a database connection.
	Bootstrap() error

	// Shutdown shuts down the app, e.g. by closing the database connection.
	Shutdown() error

	// -----------------
	// Events
	// -----------------

	OnTerminate() *hook.Hook[*TerminateEvent]
}
