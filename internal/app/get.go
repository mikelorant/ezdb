package app

import (
	"fmt"
)

type GetSessionOptions struct {
	Context string
}

type GetVariableOptions struct {
	Context 	string
	Variable	string
}

func (a *App) GetSession(opts GetSessionOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err := a.GetDBClient(context)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	out, err := cl.ShowSession()
	if err != nil {
		return fmt.Errorf("unable to get session: %w", err)
	}

	fmt.Print(cl.Format(out))

	return nil
}

func (a *App) GetVariable(opts GetVariableOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err := a.GetDBClient(context)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	out, err := cl.ShowVariable(opts.Variable)
	if err != nil {
		return fmt.Errorf("unable to get variable: %w", err)
	}

	fmt.Print(cl.Format(out))

	return nil
}
