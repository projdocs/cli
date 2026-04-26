package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/projdocs/cli/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cmd.ProjDocs.SilenceErrors = true
	cmd.ProjDocs.SilenceUsage = true

	if err := cmd.ProjDocs.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error: %s", err.Error()))
		os.Exit(1)
	}
	os.Exit(0)
}
