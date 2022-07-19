package app

import (
	"fmt"
	"log"

	"github.com/mikelorant/ezdb2/internal/dbutil"
	"github.com/mikelorant/ezdb2/internal/passutil"
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

	fmt.Print(dbutil.Format(out))

	return nil
}

func (a *App) CreateUser(opts CreateUserOptions) error {
	cl, err := a.GetDBClient(opts.Context)

	password := opts.Password
	if password == "" {
		password = passutil.Generate()
	}

	if err := cl.CreateUserGrant(opts.Name, password, opts.Database); err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}

	query := fmt.Sprintf("SHOW GRANTS FOR '%v'", opts.Name)
	out, err := cl.Query(query)
	if err != nil {
		return fmt.Errorf("unable to query: %w", err)
	}

	fmt.Print(dbutil.Format(out))

	log.Printf("Created user: %v with grants for database: %v\n", opts.Name, opts.Database)
	if opts.Password == "" {
		log.Printf("No password provided. Generated password is: %v\n", password)
	}

	return nil
}

func Backup()  {}
func Copy()    {}
func Restore() {}
