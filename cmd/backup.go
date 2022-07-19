package cmd

import (
	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewBackupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			app.Backup()

			return nil
		},
	}

	cmd.Flags().String("context", "", "Database context")
	cmd.Flags().String("name", "", "Database name")

	return cmd
}
