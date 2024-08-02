package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/heyjorgedev/deploykit/http"
	"github.com/heyjorgedev/deploykit/sqlite"
)

func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	program := NewProgram()

	// Execute program
	if err := program.Run(ctx); err != nil {
		program.Close()
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	if err := program.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Program struct {
	DB         *sqlite.DB
	HTTPServer *http.Server
}

func NewProgram() *Program {
	return &Program{
		HTTPServer: http.NewServer(),
	}
}

func (p *Program) Run(ctx context.Context) error {
	p.DB = sqlite.NewDB(ctx, "./db.sqlite")
	if err := p.DB.Open(); err != nil {
		return err
	}

	if err := p.DB.RunMigrations(); err != nil {
		return err
	}

	p.HTTPServer.Addr = "127.0.0.1:8080"

	teamService := sqlite.NewTeamService(p.DB)
	p.HTTPServer.TeamService = teamService

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
