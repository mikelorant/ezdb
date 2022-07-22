package cmd

import (
	"log"
	"strings"

	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	var context string

	cmd := &cobra.Command{
		Use:   "run",
		Short: "A brief description of your command",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.RunOptions{
				Context: context,
				Command: strings.Join(args, " "),
			}
			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.Run(opts); err != nil {
				log.Fatalf("unable to query: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")

	return cmd
}
