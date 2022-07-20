package app

import (
	"fmt"
	"log"

	"github.com/mikelorant/ezdb2/internal/database"
	"github.com/mikelorant/ezdb2/internal/password"
)

type App struct {
	Config Config
}

type QueryOptions struct {
	Context string
	Name    string
	Query   string
}

type CreateUserOptions struct {
	Context  string
	Name     string
	Password string
	Database string
}

type BackupOptions struct {
	Context string
	Name    string
	Store   string
}

type RestoreOptions struct {
	Context  string
	Name     string
	Store    string
	Filename string
}

type CopyOptions struct {
	FromContext string
	FromName    string
	ToContext   string
	ToName      string
}

var IgnoreDatabases = []string{
	"sys",
	"mysql",
	"performance_schema",
	"information_schema",
}

func New() (*App, error) {
	cfg, err := newConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load config: %w", err)
	}

	return &App{
		Config: cfg,
	}, nil
}

func (a *App) Query(opts QueryOptions) error {
	cl, err := a.GetDBClient(opts.Context)

	out, err := cl.Query(opts.Query)
	if err != nil {
		return fmt.Errorf("unable to query: %w", err)
	}

	fmt.Print(database.Format(out))

	return nil
}

func (a *App) CreateUser(opts CreateUserOptions) error {
	cl, err := a.GetDBClient(opts.Context)

	pass := opts.Password
	if pass == "" {
		pass = password.Generate()
	}

	if err := cl.CreateUserGrant(opts.Name, pass, opts.Database); err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}

	q := fmt.Sprintf("SHOW GRANTS FOR '%v'", opts.Name)
	out, err := cl.Query(q)
	if err != nil {
		return fmt.Errorf("unable to query: %w", err)
	}

	fmt.Print(database.Format(out))

	log.Printf("Created user: %v with grants for database: %v\n", opts.Name, opts.Database)
	if opts.Password == "" {
		log.Printf("No password provided. Generated password is: %v\n", pass)
	}

	return nil
}

func (a *App) Backup(opts BackupOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err := a.GetDBClient(context)

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

	cl, err = a.GetDBClient(context,
		WithDBName(name),
	)

	dbSize, err := cl.GetDatabaseSize(name)
	if err != nil {
		return fmt.Errorf("unable to get database size: %w", err)
	}

	storer, err := GetStorer(storeCfg)
	if err != nil {
		fmt.Errorf("unable to get storer: %w", err)
	}

	location, err := cl.Backup(name, dbSize, storer, true)
	if err != nil {
		return fmt.Errorf("unable to backup database: %v: %w", name, err)
	}

	log.Println("Database successfully dumped.")
	log.Println("Location:", location)

	return nil
}

func (a *App) Restore(opts RestoreOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err := a.GetDBClient(context)

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

	storer, err := GetStorer(storeCfg)
	if err != nil {
		fmt.Errorf("unable to get storer: %w", err)
	}

	filenames, err := storer.List()
	if err != nil {
		return fmt.Errorf("unable to list store: %w", err)
	}
	filename, err := Select(opts.Filename, filenames, "Choose a file:")
	if err != nil {
		return fmt.Errorf("unable to select a file: %w", err)
	}

	if err := cl.Restore(opts.Name, filename, storer, true); err != nil {
		return fmt.Errorf("unable to backup database: %v: %w", opts.Name, err)
	}

	log.Println("Database successfully restored.")

	return nil
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
