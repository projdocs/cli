package cli

import (
	"github.com/projdocs/cli/cmd/federation"
	initCmd "github.com/projdocs/cli/cmd/init"
	"github.com/projdocs/cli/cmd/serve"
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

func init() {
	ProjDocs.AddCommand(initCmd.Command)
	ProjDocs.AddCommand(serve.Command)
	ProjDocs.AddCommand(federation.Command)
}
