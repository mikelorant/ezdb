package database

import (
	"context"
	"database/sql"
	"net"

	"github.com/mikelorant/ezdb2/internal/database/mysqlshim"
	"github.com/mikelorant/ezdb2/internal/database/postgresshim"
)

type Shim interface {
	BackupCommand(verbose bool) string
	RestoreCommand(verbose bool) string
	CreateDatabase(name string) error
	CreateUserGrant(name, password, database string) error
	DropDatabase(name string) error
	GetDatabaseSize(name string) (int64, error)
	IsDatabase(name string) bool
	ListDatabases() ([]string, error)
	Query(query string) ([][]string, error)
	Format(rows [][]string) string
	ShowSession() ([][]string, error)
	ShowVariable(variable string) ([][]string, error)
}

type Config struct {
	Engine   string `yaml:"engine"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type Database struct {
	Config Config
	DB     *sql.DB
	shim   Shim
}

func New(cfg Config, dialFunc func(ctx context.Context, address string) (net.Conn, error)) (*Database, error) {
	switch cfg.Engine {
	case "mysql":
		shim, _ := mysqlshim.New(NewMySQLConfig(cfg), dialFunc)
		return &Database{
			Config: cfg,
			DB:     shim.DB,
			shim:   shim,
		}, nil
	case "postgres":
		shim, _ := postgresshim.New(NewPostgresConfig(cfg), dialFunc)
		return &Database{
			Config: cfg,
			DB:     shim.DB,
			shim:   shim,
		}, nil
	}

	return &Database{}, nil
}

func (d *Database) BackupCommand(verbose bool) string {
	return d.shim.BackupCommand(verbose)
}

func (d *Database) RestoreCommand(verbose bool) string {
	return d.shim.RestoreCommand(verbose)
}

func (d *Database) CreateDatabase(name string) error {
	return d.shim.CreateDatabase(name)
}

func (d *Database) CreateUserGrant(name, password, database string) error {
	return d.shim.CreateUserGrant(name, password, database)
}

func (d *Database) DropDatabase(name string) error {
	return d.shim.DropDatabase(name)
}

func (d *Database) GetDatabaseSize(name string) (int64, error) {
	return d.shim.GetDatabaseSize(name)
}

func (d *Database) IsDatabase(name string) bool {
	return d.shim.IsDatabase(name)
}

func (d *Database) ListDatabases() ([]string, error) {
	return d.shim.ListDatabases()
}

func (d *Database) Query(query string) ([][]string, error) {
	return d.shim.Query(query)
}

func (d *Database) Format(rows [][]string) string {
	return d.shim.Format(rows)
}

func (d *Database) ShowSession() ([][]string, error) {
	return d.shim.ShowSession()
}

func (d *Database) ShowVariable(variable string) ([][]string, error) {
	return d.shim.ShowVariable(variable)
}
