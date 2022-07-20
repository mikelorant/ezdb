package app

import (
	"fmt"
	"log"
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

	done := make(chan bool, 2)

	go func() error {
		if _, err := fromClient.Backup(fromName, fromDBSize, storer, true); err != nil {
			return fmt.Errorf("unable to backup source database: %v: %w", fromName, err)
		}

		done <- true

		return nil
	}()

	go func() error {
		if err := toClient.Restore(toName, "", storer, false); err != nil {
			return fmt.Errorf("unable to restore target database: %v: %w", toName, err)
		}

		done <- true

		return nil
	}()

	<-done

	log.Println("Database successfully copied.")

	return nil
}
