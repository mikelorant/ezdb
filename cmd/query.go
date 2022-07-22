package cmd

import (
	"log"
	"strings"

	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewQueryCmd() *cobra.Command {
	var (
		context string
		name    string
	)

	cmd := &cobra.Command{
		Use:   "query",
		Short: "A brief description of your command",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.QueryOptions{
				Context: context,
				Name:    name,
				Query:   strings.Join(args, " "),
			}
			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.Query(opts); err != nil {
				log.Fatalf("unable to query: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")
	cmd.Flags().StringVar(&name, "name", "", "Database name")

	return cmd
}
