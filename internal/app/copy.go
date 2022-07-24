package app

import (
	"fmt"
	"log"

	"golang.org/x/sync/errgroup"
)

type CopyOptions struct {
	FromContext string
	FromName    string
	ToContext   string
	ToName      string
}

func (a *App) Copy(opts CopyOptions) error {
	fromContext, err := Select(opts.FromContext, a.Config.getContexts(), "Choose a source context:")
	if err != nil {
		return fmt.Errorf("unable to select a source context: %w", err)
	}

	fromClient, err := a.GetDBClient(fromContext)
	if err != nil {
		return fmt.Errorf("unable to get source database client: %w", err)
	}

	fromDB, err := fromClient.ListDatabases()
	if err != nil {
		return fmt.Errorf("unable to list source databases: %w", err)
	}

	fromName, err := SelectWithExclude(opts.FromName, fromDB, "Choose a source database:", IgnoreDatabases)
	if err != nil {
		return fmt.Errorf("unable to select a source database: %w", err)
	}

	toContext, err := Select(opts.ToContext, a.Config.getContexts(), "Choose a target context:")
	if err != nil {
		return fmt.Errorf("unable to select a target context: %w", err)
	}

	toName := opts.ToName

	fromClient, err = a.GetDBClient(fromContext,
		WithDBName(fromName),
	)
	if err != nil {
		return fmt.Errorf("unable to get source database client: %w", err)
	}

	fromDBSize, err := fromClient.GetDatabaseSize(fromName)
	if err != nil {
		return fmt.Errorf("unable to get source database size: %w", err)
	}

	fromShell, err := a.GetShell(fromContext)
	if err != nil {
		return fmt.Errorf("unable to get a source shell: %w", err)
	}

	toClient, err := a.GetDBClient(toContext)
	if err != nil {
		return fmt.Errorf("unable to get source database client: %w", err)
	}

	if toClient.IsDatabase(toName) {
		log.Println("Found existing database:", toName)
		log.Println("Press enter to drop database.")
		fmt.Scanln()

		if err := toClient.DropDatabase(toName); err != nil {
			return fmt.Errorf("unable to drop target database: %v: %w", toName, err)
		}
	}

	if err := toClient.CreateDatabase(toName); err != nil {
		return fmt.Errorf("unable to create target database: %v: %w", toName, err)
	}

	toClient, err = a.GetDBClient(toContext,
		WithDBName(toName),
	)
	if err != nil {
		return fmt.Errorf("unable to get target database client: %w", err)
	}

	toShell, err := a.GetShell(toContext)
	if err != nil {
		return fmt.Errorf("unable to get a target shell: %w", err)
	}

	storer, err := GetStorer(&Store{
		Type: "pipe",
	})
	if err != nil {
		return fmt.Errorf("unable to get storer: %w", err)
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		if _, err := backup(fromClient, fromName, fromDBSize, storer, fromShell, true); err != nil {
			return fmt.Errorf("unable to backup source database: %v: %w", fromName, err)
		}

		return nil
	})

	g.Go(func() error {
		if _, err := restore(toClient, toName, "", storer, toShell, false); err != nil {
			return fmt.Errorf("unable to restore target database: %v: %w", toName, err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("unable to copy databases: %w", err)
	}

	log.Println("Database successfully copied.")

	return nil
}
