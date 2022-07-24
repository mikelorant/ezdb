package app

import (
	"fmt"

	"github.com/mikelorant/ezdb2/internal/shell"
)

func (a *App) GetShell(context string) (*shell.Shell, error) {
	_, tun := a.Config.getContext(context)

	if tun == nil {
		return shell.New(shell.Config{})
	}

	sess, err := getTunnelSession(tun)
	if err != nil {
		return nil, fmt.Errorf("unable to get tunnel client: %w", err)
	}

	return shell.New(shell.Config{
		Session: sess,
	})
}
