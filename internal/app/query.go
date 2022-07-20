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
	cl, err := a.GetDBClient(opts.Context)
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
