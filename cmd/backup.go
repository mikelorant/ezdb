package cmd

import (
	"log"

	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewBackupCmd() *cobra.Command {
	var (
		context string
		name    string
		store   string
	)

	cmd := &cobra.Command{
		Use:   "backup",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.BackupOptions{
				Context: context,
				Name:    name,
				Store:   store,
			}

			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.Backup(opts); err != nil {
				log.Fatalf("unable to backup database: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")
	cmd.Flags().StringVar(&name, "name", "", "Database name")
	cmd.Flags().StringVar(&store, "store", "", "Storage name")

	return cmd
}
