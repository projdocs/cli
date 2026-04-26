package cmd

import (
	"github.com/projdocs/cli/pkg"
	"github.com/spf13/cobra"
)

var ProjDocs = &cobra.Command{
	Use:     "projdocs",
	Short:   "A CLI for managing a ProjDocs instance",
	Version: pkg.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
