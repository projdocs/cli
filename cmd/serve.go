package cmd

import (
	"fmt"

	"github.com/moby/moby/client"
	"github.com/projdocs/cli/config"
	"github.com/projdocs/cli/pkg/dkr"
	"github.com/projdocs/cli/pkg/services"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the ProjDocs server",
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			docker *dkr.DockerClient
			cfg    *config.Config
		)

		// load config
		if cfgFile, err := config.LoadFile(); err != nil {
			return fmt.Errorf("could not load config file: %w", err)
		} else if cfg, err = config.FromFile(cfgFile); err != nil {
			return fmt.Errorf("could not build config: %w", err)
		}

		// setup docker
		if api, err := client.New(); err != nil {
			return fmt.Errorf("could not initialize docker client: %w", err)
		} else {
			docker = dkr.NewClient(api)
		}

		// ping docker
		if err := docker.Ping(cmd.Context()); err != nil {
			return fmt.Errorf("could not ping docker client: %w", err)
		}

		// create the runner
		runner := services.
			NewRunner(docker, services.GetAll()...).
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
