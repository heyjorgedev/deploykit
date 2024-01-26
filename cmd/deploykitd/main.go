package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/heyjorgedev/deploykit"
	"github.com/heyjorgedev/deploykit/mock"
	"github.com/heyjorgedev/deploykit/security"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/heyjorgedev/deploykit/docker"
	"github.com/heyjorgedev/deploykit/http"
	"github.com/heyjorgedev/deploykit/sqlite"
	"github.com/pelletier/go-toml"
)

func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	p := NewProgram()

	// Parse command line flags & load configuration.
	if err := p.ParseFlags(ctx, os.Args[1:]); err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Run the program
	if err := p.Run(ctx); err != nil {
		p.Close()
		fmt.Fprintln(os.Stderr, err)
		// wtf.ReportError(ctx, err)
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := p.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Program struct {
	Config     Config
	ConfigPath string

	DB         *sqlite.DB
	Docker     *docker.Connection
	HTTPServer *http.Server
}

func NewProgram() *Program {
	return &Program{
		Config:     DefaultConfig(),
		ConfigPath: DefaultConfigPath,

		Docker:     docker.NewConnection(),
		DB:         sqlite.NewDB(""),
		HTTPServer: http.NewServer(),
	}
}

func (p *Program) Run(ctx context.Context) (err error) {

	if p.DB.DSN, err = expandDSN(p.Config.DB.DSN); err != nil {
		return fmt.Errorf("cannot expand dsn: %w", err)
	}
	if err := p.DB.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}
	if err := p.Docker.Open(); err != nil {
		return fmt.Errorf("cannot connect to docker: %w", err)
	}

	// Setup Services
	projectService := &deploykit.ProjectService{}
	userService := &mock.UserService{
		FindByUsernameFunc: func(username string) (*deploykit.User, error) {
			return &deploykit.User{
				ID:       100,
				Name:     "User",
				Username: "example",
				// password
				Password: "$2a$10$P58UnXr5fBlvLDtUQ9NQNe3GODFzC/Mf/APTM1kUykpPbM43xNDfG",
			}, nil
		},
	}
	authService := security.NewAuthService(userService)

	// Setup HTTP Server Configurations
	p.HTTPServer.Addr = p.Config.HTTP.ListenAddr

	// Setup HTTP Server Dependencies
	p.HTTPServer.ProjectService = projectService
	p.HTTPServer.AuthService = authService
	p.HTTPServer.UserService = userService

	// Setup HTTP Sessions
	sessionManager := scs.New()
	sessionManager.Cookie.Name = "deploykit_session"
	sqliteSessionStore := sqlite3store.NewWithCleanupInterval(p.DB.DB, 30*time.Minute)
	defer sqliteSessionStore.StopCleanup()
	sessionManager.Store = sqliteSessionStore
	p.HTTPServer.SessionManager = sessionManager

	// TODO: Implement CSRF

	// Start the HTTP server.
	if err := p.HTTPServer.Open(); err != nil {
		return err
	}

	return nil
}

func (p *Program) Close() error {
	if p.HTTPServer != nil {
		if err := p.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if p.DB != nil {
		if err := p.DB.Close(); err != nil {
			return err
		}
	}
	if p.Docker != nil {
		if err := p.Docker.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Program) ParseFlags(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("deploykitd", flag.ContinueOnError)
	fs.StringVar(&p.ConfigPath, "config", DefaultConfigPath, "config path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	configPath, err := expand(p.ConfigPath)
	if err != nil {
		return err
	}

	// Read our TOML formatted configuration file.
	config, err := ReadConfigFile(configPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", p.ConfigPath)
	} else if err != nil {
		return err
	}
	p.Config = config

	return nil
}

const (
	DefaultConfigPath = "~/deploykitd.conf"
	DefaultDSN        = "~/.deploykitd/db"
	DefaultDataDir    = "~/.deploykitd/"
)

type Config struct {
	HTTP struct {
		ListenAddr string `toml:"listen_addr"`
	} `toml:"http"`
	DB struct {
		DSN string `toml:"dsn"`
	} `toml:"db"`
	DeployKit struct {
		DataDir string `toml:"data_dir"`
	} `toml:"deploykit"`
	SSL struct {
		IssuerEmail string `toml:"issuer_email"`
	} `toml:"ssl"`
}

func DefaultConfig() Config {
	var config Config
	config.DB.DSN = DefaultDSN
	config.DeployKit.DataDir = DefaultDataDir
	return config
}

func ReadConfigFile(filename string) (Config, error) {
	config := DefaultConfig()
	if buf, err := os.ReadFile(filename); err != nil {
		return config, err
	} else if err := toml.Unmarshal(buf, &config); err != nil {
		return config, err
	}
	return config, nil
}

func expand(path string) (string, error) {
	// Ignore if path has no leading tilde.
	if path != "~" && !strings.HasPrefix(path, "~"+string(os.PathSeparator)) {
		return path, nil
	}

	// Fetch the current user to determine the home path.
	u, err := user.Current()
	if err != nil {
		return path, err
	} else if u.HomeDir == "" {
		return path, fmt.Errorf("home directory unset")
	}

	if path == "~" {
		return u.HomeDir, nil
	}
	return filepath.Join(u.HomeDir, strings.TrimPrefix(path, "~"+string(os.PathSeparator))), nil
}

func expandDSN(dsn string) (string, error) {
	if dsn == ":memory:" {
		return dsn, nil
	}
	return expand(dsn)
}
