package main

import (
	"fmt"
	"os"

	"github.com/heyjorgedev/deploykit"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

func (p *Program) deployCmd() *cobra.Command {
	var configPath string
	c := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy",
		Long:  "Deploy",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := expand(configPath)
			if err != nil {
				return err
			}

			buf, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}

			var config deploykit.DeploymentConfig
			toml.Unmarshal(buf, &config)

			resp, err := p.HTTPClient.AppsCreate(cmd.Context(), deploykit.App{
				Name: config.App.Name,
			})
			if err != nil {
				return err
			}

			if !resp.Success {
				return fmt.Errorf("failed to create app: %s", resp.Message)
			}

			fmt.Println(resp.Message)
			return nil
		},
	}

	c.Flags().StringVarP(&configPath, "file", "f", "./deploykit.toml", "Deployment config file")

	return c
}
