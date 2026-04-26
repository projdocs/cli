package cmd

import (
	"fmt"
	"os"

	"github.com/projdocs/cli/config"
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
		return nil
	},
}

func init() {
	ProjDocs.AddCommand(initCmd)
}
