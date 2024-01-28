package core

import (
	"database/sql"
	docker "github.com/docker/docker/client"
	"github.com/heyjorgedev/deploykit/pkg/tools/bus"
	"log/slog"
	"os"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	// Internals
	db         *sql.DB
	logger     *slog.Logger
	dataDir    string
	dockerHost *docker.Client

	// Events
	terminateEvent *bus.EventBag[*TerminateEvent]
}

type BaseAppConfig struct {
	IsDev            bool
	DataDir          string
	EncryptionEnv    string
	DataMaxOpenConns int // default to 500
	DataMaxIdleConns int // default 20
	LogsMaxOpenConns int // default to 100
	LogsMaxIdleConns int // default to 5
}

func NewBaseApp(config BaseAppConfig) *BaseApp {
	return &BaseApp{
		dataDir: config.DataDir,

		// Events
		terminateEvent: &bus.EventBag[*TerminateEvent]{},
	}
}

func (app *BaseApp) DB() *sql.DB {
	return app.db
}

func (app *BaseApp) Logger() *slog.Logger {
	return app.logger
}

func (app *BaseApp) DataDir() string {
	return app.dataDir
}

func (app *BaseApp) HostDocker() *docker.Client {
	return app.dockerHost
}

func (app *BaseApp) Bootstrap() error {
	// ensure that data dir exist
	if err := os.MkdirAll(app.DataDir(), os.ModePerm); err != nil {
		return err
	}

	// connect to the host server docker socket
	if err := app.initDockerConnection(); err != nil {
		return err
	}

	// todo: connect to the database
	if err := app.initDatabaseConnection(); err != nil {
		return err
	}

	if err := app.initLogger(); err != nil {
		return err
	}

	return nil
}

func (app *BaseApp) Shutdown() error {
	return nil
}

func (app *BaseApp) initDockerConnection() (err error) {
	app.dockerHost, err = docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	return err
}

func (app *BaseApp) initDatabaseConnection() error {
	return nil
}

func (app *BaseApp) initLogger() error {
	app.logger = slog.Default()
	return nil
}

func (app *BaseApp) OnTerminate() *bus.EventBag[*TerminateEvent] {
	return app.terminateEvent
}
