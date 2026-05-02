package admin

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "admin",
	Short: "Manage users with administrative access to the ProjDocs cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.AddCommand(
		create,
		list,
		rm,
	)
}
