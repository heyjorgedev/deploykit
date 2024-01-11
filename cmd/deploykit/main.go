package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/jorgemurta/deploykit"
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

	r.AddCommand(p.appsCmd())

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
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := p.HTTPClient.AppsList(context.Background())
			if err != nil {
				return err
			}

			fmt.Printf("Apps (%d):\n", len(resp.Data))

			sort.Slice(resp.Data, func(i, j int) bool {
				return resp.Data[i].Name < resp.Data[j].Name
			})

			for _, app := range resp.Data {
				fmt.Printf("- %s\n", app.Name)
			}

			return nil
		},
	}
}

func (p *Program) appsCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Say hello",
		Long:  "Say hello to the world",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			resp, err := p.HTTPClient.AppsCreate(context.Background(), deploykit.App{
				Name: appName,
			})

			if err != nil {
				return err
			}

			if !resp.Success {
				return errors.New(resp.Message)
			}

			fmt.Println(resp.Message)
			return nil
		},
	}
}
