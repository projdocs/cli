package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/projdocs/cli/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"setup"},
	Short:   "Setup a ProjDocs server",
	RunE: func(cmd *cobra.Command, args []string) error {

		if cfg, _ := config.LoadFile(); cfg != nil {
			return fmt.Errorf("projdocs already initialized")
		}

		if _, err := os.Stat(config.GetConfigFilePath()); err != nil {
			if os.IsNotExist(err) {

				if err := os.MkdirAll(config.GetConfigDirPath(), 0o755); err != nil {
					return fmt.Errorf("could not create parent directories for config file: %w", err)
				}
				if err := os.WriteFile(config.GetConfigFilePath(), config.TemplateConfigFile, 0o644); err != nil {
					return fmt.Errorf("could not write config file: %w", err)
				}
			} else {
				return fmt.Errorf("unable to check for config file: %w", err)
			}
		}

		postgresDir := filepath.Join(config.GetConfigDirPath(), "postgres")
		if err := os.MkdirAll(postgresDir, 0o755); err != nil {
			return fmt.Errorf("could not create postgres directory: %w", err)
		}

		storageDir := filepath.Join(config.GetConfigDirPath(), "storage")
		if err := os.MkdirAll(storageDir, 0o755); err != nil {
			return fmt.Errorf("could not create storage directory: %w", err)
		}

		color.Green("initialized projdocs!")
		color.Green(" - data directory: %s", config.GetConfigDirPath())
		color.Green(" - config file: %s", config.GetConfigFilePath())

		return nil
	},
}

func init() {
	ProjDocs.AddCommand(initCmd)
}
