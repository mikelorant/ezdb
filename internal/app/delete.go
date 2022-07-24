package app

import (
	"fmt"
	"log"
)

type DeleteDatabaseOptions struct {
	Context string
	Name    string
}

func (a *App) DeleteDatabase(opts DeleteDatabaseOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err := a.GetDBClient(context)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	if !cl.IsDatabase(opts.Name) {
		log.Printf("Database %v not found.", opts.Name)
		return nil
	}

	log.Println("Found existing database:", opts.Name)
	log.Println("Press enter to delete database.")
	fmt.Scanln()

	if err := cl.DropDatabase(opts.Name); err != nil {
		return fmt.Errorf("unable to drop target database: %v: %w", opts.Name, err)
	}

	log.Printf("Deleted database: %v\n", opts.Name)

	return nil
}
