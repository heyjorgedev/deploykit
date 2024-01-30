package deploykit

import (
	"github.com/heyjorgedev/deploykit/pkg/caddy"
	"github.com/heyjorgedev/deploykit/pkg/core"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

var Version = "dev"

var _ core.App = (*DeployKit)(nil)

type appWrapper struct {
	core.App
}

type DeployKit struct {
	*appWrapper

	dataDirFlag string

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
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}

	baseDir, isDev := inspectRuntime()
	if config.DataDir == "" {
		config.DataDir = filepath.Join(baseDir, "dk-data")
	}

	dk.appWrapper = &appWrapper{
		App: core.NewBaseApp(core.BaseAppConfig{
			DataDir: config.DataDir,
			IsDev:   isDev,
		}),
	}

	dk.RootCmd.PersistentFlags().StringVar(
		&dk.dataDirFlag,
		"dir",
		config.DataDir,
		"directory where DeployKit will store all data.",
	)

	dk.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	_ = dk.RootCmd.ParseFlags(os.Args[1:])

	return dk
}

func (dk *DeployKit) Start() error {
	dk.RootCmd.AddCommand(newServeCommand(dk))
	return dk.Execute()
}

func (dk *DeployKit) Bootstrap() error {
	// Register Dependencies
	caddy.NewManager(dk)

	return dk.appWrapper.Bootstrap()
}

func (dk *DeployKit) Execute() error {
	if !dk.skipBootstrap() {
		if err := dk.Bootstrap(); err != nil {
			return err
		}
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

func (dk *DeployKit) skipBootstrap() bool {
	flags := []string{
		"-h",
		"--help",
		"-v",
		"--version",
	}

	if dk.IsBootstrapped() {
		return true // already bootstrapped
	}

	cmd, _, err := dk.RootCmd.Find(os.Args[1:])
	if err != nil {
		return true // unknown command
	}
	if cmd == dk.RootCmd {
		return true
	}

	for _, arg := range os.Args {
		if !ExistInSlice(arg, flags) {
			continue
		}

		// ensure that there is no user defined flag with the same name/shorthand
		trimmed := strings.TrimLeft(arg, "-")
		if len(trimmed) > 1 && cmd.Flags().Lookup(trimmed) == nil {
			return true
		}
		if len(trimmed) == 1 && cmd.Flags().ShorthandLookup(trimmed) == nil {
			return true
		}
	}

	return false
}

// ExistInSlice checks whether a comparable element exists in a slice of the same type.
func ExistInSlice[T comparable](item T, list []T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}

	return false
}

func inspectRuntime() (baseDir string, withGoRun bool) {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// probably ran with go run
		withGoRun = true
		baseDir, _ = os.Getwd()
	} else {
		// probably ran with go build
		withGoRun = false
		baseDir = filepath.Dir(os.Args[0])
	}
	return
}
