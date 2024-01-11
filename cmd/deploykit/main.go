package main

import (
	"context"
	"os"

	"github.com/jorgemurta/deploykit/http"
	"github.com/spf13/cobra"
)

func main() {
	p := &Program{}
	if err := p.Run(context.Background()); err != nil {
		p.Close()
		os.Exit(1)
	}
	defer p.Close()
}

type Program struct {
	HTTPClient *http.CLIClient
}

func (p *Program) Run(ctx context.Context) error {
	p.HTTPClient = http.NewCliClient("http://127.0.0.1:8080")

	r := p.rootCmd()

	r.AddCommand(p.databaseRedisListCmd())
	r.AddCommand(p.appsCreateCmd(ctx))
	r.AddCommand(p.appsListCmd(ctx))

	return r.Execute()
}

func (p *Program) Close() error {
	return nil
}

func (p *Program) rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "deploykit",
		Short:        "Hugo is a very fast static site generator",
		Long:         `Deploy`,
		SilenceUsage: true,
	}
}
