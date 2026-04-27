package internal

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type Spinner struct {
	s *spinner.Spinner
}

func NewSpinner(message string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 80*time.Millisecond)
	s.Suffix = " " + message
	s.Start()
	return &Spinner{s: s}
}

func (s *Spinner) Update(message string) {
	s.s.Suffix = " " + message
}

func (s *Spinner) Stop() {
	s.s.Stop()
}

func (s *Spinner) Success(message string) {
	s.s.FinalMSG = color.GreenString("✓ " + message + "\n")
	s.s.Stop()
}

func (s *Spinner) Fail(message string) {
	s.s.FinalMSG = color.RedString("✗ " + message + "\n")
	s.s.Stop()
}
