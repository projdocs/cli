package cmd

import (
	"fmt"

	"github.com/moby/moby/client"
	config2 "github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/pkg/docker"
	"github.com/projdocs/cli/pkg/services"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the ProjDocs server",
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			dkr *docker.Client
			cfg *config2.Config
		)

		// load config
		if cfgFile, err := config2.LoadFile(); err != nil {
			return fmt.Errorf("could not load config file: %w", err)
		} else if cfg, err = config2.FromFile(cfgFile); err != nil {
			return fmt.Errorf("could not build config: %w", err)
		}

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

		// create the runner
		runner := services.
			NewRunner(dkr, services.GetAll()...).
			Build(*cfg)

		// start server
		if err := runner.Start(cmd.Context()); err != nil {
			// TODO: cleanup

			return fmt.Errorf("could not start runner: %w", err)
		}
		return nil
	},
}

func init() {
	ProjDocs.AddCommand(serveCmd)
}
