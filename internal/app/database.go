package app

import (
	"fmt"

	"github.com/mikelorant/ezdb2/internal/database"
	"github.com/mikelorant/ezdb2/internal/printer"
)

type Databases []Database

type Database struct {
	Context  string `yaml:"context"`
	Engine   string `yaml:"engine"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
	Tunnel   string `yaml:"tunnel"`
}

type DBOptions struct {
	name string
}

func (d Database) String() string {
	return printer.Struct(d)
}

func WithDBName(name string) func(*DBOptions) {
	return func(d *DBOptions) {
		d.name = name
	}
}

func (a *App) GetDBClient(context string, dbOpts ...func(*DBOptions)) (*database.Database, error) {
	var dbOptions DBOptions
	for _, o := range dbOpts {
		o(&dbOptions)
	}

	db, tun := a.Config.getContext(context)

	if dbOptions.name != "" {
		db.Name = dbOptions.name
	}

	cfg := database.Config{
		Engine:   db.Engine,
		Host:     db.Host,
		Port:     db.Port,
		User:     db.User,
		Name:     db.Name,
		Password: db.Password,
	}

	dialer, err := getDialFunc(tun)
	if err != nil {
		return nil, fmt.Errorf("unable to get dialer: %w", err)
	}

	return database.New(cfg, dialer)
}
