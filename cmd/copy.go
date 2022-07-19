package cmd

import (
	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewCopyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "copy",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.Copy()

			return nil
		},
	}

	cmd.Flags().String("from-context", "", "Database context")
	cmd.Flags().String("from-name", "", "Database name")
	cmd.Flags().String("to-context", "", "Database context")
	cmd.Flags().String("to-name", "", "Database name")

	return cmd
}
