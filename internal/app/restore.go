package app

import (
	"fmt"
	"log"
)

type RestoreOptions struct {
	Context  string
	Name     string
	Store    string
	Filename string
}

func (a *App) Restore(opts RestoreOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err := a.GetDBClient(context)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	if err := cl.CreateDatabase(opts.Name); err != nil {
		return fmt.Errorf("unable to create database: %v: %w", opts.Name, err)
	}

	store, err := Select(opts.Store, a.Config.getStores(), "Choose a store:")
	if err != nil {
		return fmt.Errorf("unable to select a store: %w", err)
	}

	storeCfg := a.Config.getStore(store)

	cl, err = a.GetDBClient(context,
		WithDBName(opts.Name),
	)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	shell, err := a.GetShell(context)
	if err != nil {
		return fmt.Errorf("unable to get a shell: %w", err)
	}

	storer, err := GetStorer(storeCfg)
	if err != nil {
		return fmt.Errorf("unable to get storer: %w", err)
	}

	filenames, err := storer.List()
	if err != nil {
		return fmt.Errorf("unable to list store: %w", err)
	}
	filename, err := Select(opts.Filename, filenames, "Choose a file:")
	if err != nil {
		return fmt.Errorf("unable to select a file: %w", err)
	}

	_, err = cl.RestoreCompat(opts.Name, filename, storer, shell, true)
	if err != nil {
		return fmt.Errorf("unable to restore database: %v: %w", opts.Name, err)
	}

	log.Println("Database successfully restored.")

	return nil
}
