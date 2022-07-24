package app

import (
	"fmt"
	"io"
	"os"
)

type Runner interface {
	Run(out io.Writer, in io.Reader, cmd string, combinedOutput bool) error
}

type RunOptions struct {
	Context string
	Command string
}

func (a *App) Run(opts RunOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	sh, err := a.GetShell(context)
	if err != nil {
		return fmt.Errorf("unable to get a shell: %w", err)
	}

	w := os.Stdout
	err = sh.Run(w, nil, opts.Command, true)
	if err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}

	return nil
}
