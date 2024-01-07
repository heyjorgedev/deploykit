package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/jorgemurta/deploykit/http"
	"github.com/jorgemurta/deploykit/sqlite"
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
	HTTPServer *http.Server
}

func NewProgram() *Program {
	return &Program{
		Config:     DefaultConfig(),
		ConfigPath: DefaultConfigPath,

		DB:         sqlite.NewDB(""),
		HTTPServer: http.NewServer(),
	}
}

func (p *Program) Run(ctx context.Context) (err error) {

	// Open the database.
	if p.DB.DSN, err = expandDSN(p.Config.DB.DSN); err != nil {
		return fmt.Errorf("cannot expand dsn: %w", err)
	}
	if err := p.DB.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	// Setup Services
	appService := sqlite.NewAppService(p.DB)

	// Setup HTTP Server Configurations
	p.HTTPServer.Addr = p.Config.HTTP.Addr

	// Setup HTTP Server Dependencies
	p.HTTPServer.AppService = appService

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
)

type Config struct {
	HTTP struct {
		Addr string `toml:"addr"`
	} `toml:"http"`
	DB struct {
		DSN string `toml:"dsn"`
	} `toml:"db"`
}

func DefaultConfig() Config {
	var config Config
	config.DB.DSN = DefaultDSN
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
