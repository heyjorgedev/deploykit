package main

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/jorgemurta/deploykit"
	"github.com/spf13/cobra"
)

func (p *Program) appsListCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "apps:list",
		Short: "Say hello",
		Long:  "Say hello to the world",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := p.HTTPClient.AppsList(ctx)
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

func (p *Program) appsCreateCmd(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "apps:create",
		Short: "Say hello",
		Long:  "Say hello to the world",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			resp, err := p.HTTPClient.AppsCreate(ctx, deploykit.App{
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
