package core

import (
	"database/sql"
	docker "github.com/docker/docker/client"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/heyjorgedev/deploykit/pkg/core/hook"
	"github.com/heyjorgedev/deploykit/pkg/dao"
	"github.com/heyjorgedev/deploykit/pkg/sqlite"
	"github.com/lmittmann/tint"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
	"path"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	// Internals
	db         *sql.DB
	dbx        *dbx.DB
	dao        *dao.Dao
	logger     *slog.Logger
	dataDir    string
	dockerHost *docker.Client
	isDev      bool

	// Events
	terminateEvent *hook.Hook[*TerminateEvent]
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
		isDev:   config.IsDev,

		// Events
		terminateEvent: &hook.Hook[*TerminateEvent]{},
	}
}

func (app *BaseApp) DB() *sql.DB {
	return app.db
}

func (app *BaseApp) Dao() *dao.Dao {
	return app.dao
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

	if err := app.initDockerConnection(); err != nil {
		return err
	}
	if err := app.initDatabaseConnection(); err != nil {
		return err
	}
	if err := app.initLogger(); err != nil {
		return err
	}

	return nil
}

func (app *BaseApp) Shutdown() error {
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			return err
		}
	}

	if app.dbx != nil {
		if err := app.dbx.Close(); err != nil {
			return err
		}
	}

	if app.dockerHost != nil {
		if err := app.dockerHost.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (app *BaseApp) initDockerConnection() (err error) {
	app.dockerHost, err = docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	return err
}

func (app *BaseApp) initDatabaseConnection() (err error) {
	app.db, err = sql.Open("sqlite3", path.Join(app.DataDir(), "db.sqlite"))
	if err != nil {
		return err
	}
	app.dbx = dbx.NewFromDB(app.db, "sqlite3")
	app.dao = dao.New(app.dbx)

	if err := sqlite.EnableWAL(app.db); err != nil {
		return err
	}

	if err := sqlite.EnableForeignKeys(app.db); err != nil {
		return err
	}

	if err := sqlite.Migrate(app.db); err != nil {
		return err
	}

	return nil
}

func (app *BaseApp) initLogger() error {
	app.logger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		AddSource: app.isDev,
	}))
	return nil
}

func (app *BaseApp) OnTerminate() *hook.Hook[*TerminateEvent] {
	return app.terminateEvent
}
