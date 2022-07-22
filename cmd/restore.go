package cmd

import (
	"log"

	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewRestoreCmd() *cobra.Command {
	var (
		context  string
		name     string
		store    string
		filename string
	)

	cmd := &cobra.Command{
		Use:   "restore",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.RestoreOptions{
				Context:  context,
				Name:     name,
				Store:    store,
				Filename: filename,
			}

			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.Restore(opts); err != nil {
				log.Fatalf("unable to restore database: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")
	cmd.Flags().StringVar(&name, "name", "", "Database name")
	cmd.Flags().StringVar(&store, "store", "", "Storage name")
	cmd.Flags().StringVar(&filename, "filename", "", "Filename")
	cmd.MarkFlagRequired("name")

	return cmd
}
