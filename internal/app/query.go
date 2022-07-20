package app

import (
	"fmt"

	"github.com/mikelorant/ezdb2/internal/database"
)

type QueryOptions struct {
	Context string
	Name    string
	Query   string
}

func (a *App) Query(opts QueryOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err := a.GetDBClient(context)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	db, err := cl.ListDatabases()
	if err != nil {
		return fmt.Errorf("unable to list databases: %w", err)
	}

	name, err := SelectWithExclude(opts.Name, db, "Choose a database:", IgnoreDatabases)
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err = a.GetDBClient(context,
		WithDBName(name),
	)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	out, err := cl.Query(opts.Query)
	if err != nil {
		return fmt.Errorf("unable to query: %w", err)
	}

	fmt.Print(database.Format(out))

	return nil
}
