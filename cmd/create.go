package cmd

import (
	"log"

	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "A brief description of your command",
	}

	cmd.AddCommand(NewCreateUserCmd())
	cmd.AddCommand(NewCreateDatabaseCmd())

	return cmd
}

func NewCreateUserCmd() *cobra.Command {
	var (
		context  string
		name     string
		password string
		database string
	)

	cmd := &cobra.Command{
		Use:   "user",
		Short: "A brief description of your command",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.CreateUserOptions{
				Context:  context,
				Name:     name,
				Password: password,
				Database: database,
			}
			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.CreateUser(opts); err != nil {
				log.Fatalf("unable to create user: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")
	cmd.Flags().StringVar(&name, "name", "", "User name")
	cmd.Flags().StringVar(&password, "password", "", "User password")
	cmd.Flags().StringVar(&database, "database", "", "Database grant for user")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("database")

	return cmd
}

func NewCreateDatabaseCmd() *cobra.Command {
	var context string

	cmd := &cobra.Command{
		Use:   "database",
		Short: "A brief description of your command",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.CreateDatabaseOptions{
				Context: context,
				Name:    args[0],
			}
			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.CreateDatabase(opts); err != nil {
				log.Fatalf("unable to create database: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")

	return cmd
}
