package app

import (
	"fmt"
	"io"
	"log"

	"github.com/mikelorant/ezdb2/internal/compress"
	"github.com/mikelorant/ezdb2/internal/database"
	"github.com/mikelorant/ezdb2/internal/progress"
	"golang.org/x/sync/errgroup"
)

type BackupOptions struct {
	Context string
	Name    string
	Store   string
}

func (a *App) Backup(opts BackupOptions) error {
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

	store, err := Select(opts.Store, a.Config.getStores(), "Choose a store:")
	if err != nil {
		return fmt.Errorf("unable to select a store: %w", err)
	}

	storeCfg := a.Config.getStore(store)

	storer, err := GetStorer(storeCfg)
	if err != nil {
		return fmt.Errorf("unable to get storer: %w", err)
	}

	cl, err = a.GetDBClient(context,
		WithDBName(name),
	)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	dbSize, err := cl.GetDatabaseSize(name)
	if err != nil {
		return fmt.Errorf("unable to get database size: %w", err)
	}

	shell, err := a.GetShell(context)
	if err != nil {
		return fmt.Errorf("unable to get a shell: %w", err)
	}

	filename := fmt.Sprintf("%v-%v", context, name)
	location, err := backup(cl, filename, dbSize, storer, shell, true)
	if err != nil {
		return fmt.Errorf("unable to backup database: %v: %w", name, err)
	}

	log.Println("Database successfully dumped.")
	log.Println("Location:", location)

	return nil
}

func backup(cl *database.Database, name string, size int64, storer Storer, shell Shell, verbose bool) (string, error) {
	desc := "Dumping..."
	bar := progress.New(size, desc, verbose)

	// mysqldump (w) -> (w) multiwriter (w) -> (w) progressbar
	//                                      -> (w) pipe (r) -> (r) gzip (r) -> (r) storer

	gzipIn, dumpOut := io.Pipe()
	dumpIn := io.MultiWriter(dumpOut, bar)
	gzipOut := compress.NewGzipCompressor(gzipIn)

	g := new(errgroup.Group)

	result := make(chan string)

	g.Go(func() error {
		location, err := storer.Store(gzipOut, name)
		result <- location
		return err
	})

	if verbose {
		log.Println("Command:", cl.BackupCommand(true))
	}

	if err := shell.Run(dumpIn, nil, cl.BackupCommand(false), false); err != nil {
		return "", fmt.Errorf("unable to run command: %w", err)
	}
	dumpOut.Close()

	location := <-result
	if err := g.Wait(); err != nil {
		return "", fmt.Errorf("store failure: %w", err)
	}
	bar.Finish()

	return location, nil
}
