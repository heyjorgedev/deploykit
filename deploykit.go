package deploykit

import (
	"github.com/heyjorgedev/deploykit/pkg/caddy"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var Version = "dev"

var _ core.App = (*DeployKit)(nil)

type appWrapper struct {
	core.App
}

type DeployKit struct {
	*appWrapper

	RootCmd *cobra.Command
}

type Config struct {
	IsDev   bool
	DataDir string
}

func NewWithConfig(config Config) *DeployKit {
	dk := &DeployKit{
		RootCmd: &cobra.Command{
			Use:     filepath.Base(os.Args[0]),
			Short:   "DeployKit",
			Version: Version,
		},
	}

	dk.appWrapper = &appWrapper{
		App: core.NewBaseApp(core.BaseAppConfig{
			DataDir: config.DataDir,
			IsDev:   config.IsDev,
		}),
	}

	dk.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	return dk
}

func (dk *DeployKit) Start() error {
	dk.RootCmd.AddCommand(newServeCommand(dk))
	return dk.Execute()
}

func (dk *DeployKit) Execute() error {
	caddy.NewManager(dk)

	if err := dk.Bootstrap(); err != nil {
		return err
	}

	exit := make(chan bool, 1)

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigch := make(chan os.Signal, 1)
		signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
		<-sigch

		exit <- true
	}()

	// execute the root command
	go func() {
		_ = dk.RootCmd.Execute()

		exit <- true
	}()

	<-exit

	return dk.OnTerminate().Trigger(&core.TerminateEvent{App: dk}, func(e *core.TerminateEvent) error {
		return dk.Shutdown()
	})
}
