package cmd

import (
	"log"

	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "A brief description of your command",
	}

	cmd.AddCommand(NewDeleteDatabaseCmd())

	return cmd
}

func NewDeleteDatabaseCmd() *cobra.Command {
	var context string

	cmd := &cobra.Command{
		Use:   "database",
		Short: "A brief description of your command",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.DeleteDatabaseOptions{
				Context: context,
				Name:    args[0],
			}
			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.DeleteDatabase(opts); err != nil {
				log.Fatalf("unable to delete database: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")

	return cmd
}
