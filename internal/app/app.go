package app

import (
	"fmt"
)

type App struct {
	Config Config
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
