package deploykit

import (
	"github.com/heyjorgedev/deploykit/pkg/core"
	"github.com/heyjorgedev/deploykit/pkg/http"
	"github.com/spf13/cobra"
)

func newServeCommand(app core.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return http.Serve(app, http.ServeConfig{HttpAddr: "127.0.0.1:8090"})
		},
	}

	return cmd
}
