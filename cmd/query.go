package cmd

import (
	"fmt"
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
				return fmt.Errorf("unable to start app: %w", err)
			}
			if err := a.Query(opts); err != nil {
				return fmt.Errorf("unable to query: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")
	cmd.Flags().StringVar(&name, "name", "", "Database name")
	cmd.MarkFlagRequired("context")

	return cmd
}
