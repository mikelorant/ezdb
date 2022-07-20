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

	fromDBSize, err := fromClient.GetDatabaseSize(fromName)
	if err != nil {
		return fmt.Errorf("unable to get source database size: %w", err)
	}

	toClient, err := a.GetDBClient(toContext)

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

	storer, err := GetStorer(&Store{
		Type: "pipe",
	})
	if err != nil {
		fmt.Errorf("unable to get storer: %w", err)
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		if _, err := fromClient.Backup(fromName, fromDBSize, storer, true); err != nil {
			return fmt.Errorf("unable to backup source database: %v: %w", fromName, err)
		}

		return nil
	})

	g.Go(func() error {
		if err := toClient.Restore(toName, "", storer, false); err != nil {
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
