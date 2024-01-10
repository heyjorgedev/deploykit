package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jorgemurta/deploykit"
	"github.com/jorgemurta/deploykit/http"
	"github.com/spf13/cobra"
)

func main() {
	p := &Program{}
	if err := p.Run(context.Background()); err != nil {
		p.Close()
		fmt.Fprintln(os.Stderr, err)
		// wtf.ReportError(ctx, err)
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

	r.AddCommand(p.appsCmd())

	return r.Execute()
}

func (p *Program) Close() error {
	return nil
}

func (p *Program) rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deploykit",
		Short: "Hugo is a very fast static site generator",
		Long:  `Deploy`,
	}
}

func (p *Program) appsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "apps",
		Short: "Say hello",
		Long:  "Say hello to the world",
	}

	c.AddCommand(p.appsCreateCmd())
	c.AddCommand(p.appsListCmd())

	return c
}

func (p *Program) appsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Say hello",
		Long:  "Say hello to the world",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, World!")
		},
	}
}

func (p *Program) appsCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Say hello",
		Long:  "Say hello to the world",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			appName := args[0]
			app, err := p.HTTPClient.AppsCreate(context.Background(), deploykit.App{
				Name: appName,
			})

			if err != nil {
				fmt.Println(err)
				return
			}

			if app.ID == 0 {
				fmt.Println("App not created")
				return
			}

			fmt.Printf("App created with ID: %d\n", app.ID)

		},
	}
}
