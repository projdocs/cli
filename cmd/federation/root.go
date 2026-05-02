package federation

import (
	"github.com/projdocs/cli/cmd/federation/admin"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:     "federation",
	Aliases: []string{"fed"},
	Short:   "Utilities for managing a ProjDocs cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.AddCommand(serve)
	Command.AddCommand(admin.Command)
}
