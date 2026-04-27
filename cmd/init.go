package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	config2 "github.com/projdocs/cli/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"setup"},
	Short:   "Setup a ProjDocs server",
	RunE: func(cmd *cobra.Command, args []string) error {

		if cfg, _ := config2.LoadFile(); cfg != nil {
			return fmt.Errorf("projdocs already initialized")
		}

		if _, err := os.Stat(config2.GetConfigFilePath()); err != nil {
			if os.IsNotExist(err) {

				if err := os.MkdirAll(config2.GetConfigDirPath(), 0o755); err != nil {
					return fmt.Errorf("could not create parent directories for config file: %w", err)
				}
				if err := os.WriteFile(config2.GetConfigFilePath(), config2.TemplateConfigFile, 0o644); err != nil {
					return fmt.Errorf("could not write config file: %w", err)
				}
			} else {
				return fmt.Errorf("unable to check for config file: %w", err)
			}
		}

		postgresDir := filepath.Join(config2.GetConfigDirPath(), "postgres")
		if err := os.MkdirAll(postgresDir, 0o755); err != nil {
			return fmt.Errorf("could not create postgres directory: %w", err)
		}

		return nil
	},
}

func init() {
	ProjDocs.AddCommand(initCmd)
}
