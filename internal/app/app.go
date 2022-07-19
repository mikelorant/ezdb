package app

import (
	"fmt"
	"log"

	"github.com/mikelorant/ezdb2/internal/database"
	"github.com/mikelorant/ezdb2/internal/password"
	"github.com/mikelorant/ezdb2/internal/selector"
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
	context := opts.Context
	if context == "" {
		c, err := selector.Select(a.Config.getContexts(),
			selector.WithMessage("Choose a context:"),
		)
		if err != nil {
			return fmt.Errorf("unable to select context: %w", err)
		}
		context = c
	}

	cl, err := a.GetDBClient(context)

	db, err := cl.ListDatabases()
	if err != nil {
		return fmt.Errorf("unable to list databases: %w", err)
	}

	// TODO: move to func
	name := opts.Name
	if name == "" {
		name, err = selector.Select(db,
			selector.WithExclude(IgnoreDatabases),
			selector.WithMessage("Choose a database:"),
		)
		if err != nil {
			return fmt.Errorf("unable to select database: %w", err)
		}
	}

	store := opts.Store
	if store == "" {
		s, err := selector.Select(a.Config.getStores(),
			selector.WithMessage("Choose a store:"),
		)
		if err != nil {
			return fmt.Errorf("unable to select stores: %w", err)
		}
		store = s
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

	location, err := cl.Backup(name, dbSize, storer)
	if err != nil {
		return fmt.Errorf("unable to backup database: %v: %w", name, err)
	}

	log.Println("Database successfully dumped.")
	log.Println("Location:", location)

	return nil
}

func Copy()    {}
func Restore() {}
