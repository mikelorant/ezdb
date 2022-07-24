package app

import (
	"fmt"
	"log"

	"github.com/sethvargo/go-password/password"
)

type CreateUserOptions struct {
	Context  string
	Name     string
	Password string
	Database string
}

func (a *App) CreateUser(opts CreateUserOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err := a.GetDBClient(context)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	pass := opts.Password
	if pass == "" {
		pass, err = password.Generate(32, 10, 0, true, false)
		if err != nil {
			return fmt.Errorf("unable to generate password: %w", err)
		}
	}

	if err := cl.CreateUserGrant(opts.Name, pass, opts.Database); err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}

	log.Printf("Created user: %v with grants for database: %v\n", opts.Name, opts.Database)
	if opts.Password == "" {
		log.Printf("No password provided. Generated password is: %v\n", pass)
	}

	return nil
}
