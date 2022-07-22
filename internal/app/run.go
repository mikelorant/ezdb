package app

import (
	"fmt"
	"os"
)

// type Shell interface {
// 	Run(out io.WriteCloser, cmd string) error
// }

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
	err = sh.Run(w, opts.Command)
	if err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}

	return nil
}
