package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/projdocs/cli/cmd"
	"golang.org/x/sys/unix"
)

func suppressInterruptEcho() (restore func()) {
	const (
		ioctlGetTermios = unix.TIOCGETA
		ioctlSetTermios = unix.TIOCSETA
	)
	fd := int(os.Stdin.Fd())

	termios, err := unix.IoctlGetTermios(fd, ioctlGetTermios)
	if err != nil {
		return func() {}
	}

	original := *termios
	termios.Lflag &^= unix.ECHOCTL

	if err := unix.IoctlSetTermios(fd, ioctlSetTermios, termios); err != nil {
		return func() {}
	}

	return func() {
		_ = unix.IoctlSetTermios(fd, ioctlSetTermios, &original)
	}
}

func main() {
	restore := suppressInterruptEcho()
	defer restore()

	cli.ProjDocs.SilenceErrors = true
	cli.ProjDocs.SilenceUsage = true

	if err := cli.ProjDocs.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("Error: %s", err.Error()))
		os.Exit(1)
	}
	os.Exit(0)
}
