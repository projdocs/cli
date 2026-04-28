package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/moby/moby/client"
	config2 "github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/internal/docker"
	"github.com/projdocs/cli/pkg/services"
	"github.com/spf13/cobra"
)

var (
	serveCmdListen *bool = new(false)
	serveCmdForce  *bool = new(false)
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the ProjDocs server",
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			dkr *docker.Client
			cfg *config2.Config
		)

		// setup docker
		if api, err := client.New(); err != nil {
			return fmt.Errorf("could not initialize docker client: %w", err)
		} else {
			dkr = docker.NewClient(api)
		}

		// ping docker
		if err := dkr.Ping(cmd.Context()); err != nil {
			return fmt.Errorf("could not ping docker client: %w", err)
		}

		// cannot be running
		if containers, err := dkr.GetContainers(cmd.Context()); err != nil {
			return fmt.Errorf("could not inspect containers: %w", err)
		} else if len(containers.Items) > 0 && !*serveCmdForce {
			return fmt.Errorf("ProjDocs already running: found %d containers (hint: re-run with the `--force` flag)", len(containers.Items))
		}

		// load config
		if cfgFile, err := config2.LoadFile(); err != nil {
			return fmt.Errorf("could not load config file: %w", err)
		} else if validationErr := cfgFile.Validate(); validationErr != nil {
			return fmt.Errorf("could not validate config file: %w", validationErr)
		} else if cfg, err = config2.FromFile(cfgFile); err != nil {
			return fmt.Errorf("could not build config: %w", err)
		}

		// create the runner
		runner := services.
			NewRunner(dkr, services.GetAll()...).
			Build(*cfg)

		// derive a cancellable context; wire up signal handling if --listen
		// derive context; wire up signal handling if --listen
		ctx := cmd.Context()
		if *serveCmdListen {
			var stop context.CancelFunc
			ctx, stop = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()
		}

		// start server
		if err := runner.Start(ctx); err != nil {
			return fmt.Errorf("could not start runner: %w", err)
		}

		// handle shutdown
		if *serveCmdListen {
			<-ctx.Done()
			fmt.Fprintf(os.Stderr, "\nreceived signal %s, shutting down...\n", context.Cause(ctx))
			runner.Stop()
			color.Blue("Goodbye!")
		}

		return nil
	},
}

func init() {
	ProjDocs.AddCommand(serveCmd)
	serveCmd.Flags().BoolVarP(serveCmdListen, "listen", "l", false, "listen for exit signals to control shutdown")
	serveCmd.Flags().BoolVarP(serveCmdForce, "force", "f", false, "force serve by removing existing containers and re-serving")

}
