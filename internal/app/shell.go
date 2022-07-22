package app

import (
	"fmt"
	"io"

	"github.com/mikelorant/ezdb2/internal/shell"
)

type Shell interface {
	Run(out io.Writer, in io.Reader, cmd string) error
}

func (a *App) GetShell(context string) (Shell, error) {
	_, tun := a.Config.getContext(context)

	if tun == nil {
		return shell.NewLocalShell(), nil
	}

	sess, err := getTunnelSession(tun)
	if err != nil {
		return nil, fmt.Errorf("unable to get tunnel client: %w", err)
	}

	return shell.NewRemoteShell(sess), nil
}
