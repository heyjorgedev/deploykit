package main

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

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
	r.AddCommand(p.appsCreateCmd())
	r.AddCommand(p.appsListCmd())

	r.AddCommand(p.deployCmd())

	return r.ExecuteContext(ctx)
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
