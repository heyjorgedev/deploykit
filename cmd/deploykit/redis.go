package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (p *Program) databaseRedisListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "redis:list",
		Short: "List Redis databases",
		Long:  "List Redis databases",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Redis databases:")
			return nil
		},
	}
}
