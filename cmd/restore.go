package cmd

import (
	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.Restore()

			return nil
		},
	}

	cmd.Flags().String("context", "", "Database context")
	cmd.Flags().String("name", "", "Database name")
	cmd.Flags().String("file", "", "SQL filename")

	return cmd
}
