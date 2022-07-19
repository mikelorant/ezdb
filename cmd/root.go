package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ezdb2",
		Short: "A brief description of your application",
	}

	cmd.AddCommand(NewCopyCmd())
	cmd.AddCommand(NewQueryCmd())
	cmd.AddCommand(NewBackupCmd())
	cmd.AddCommand(NewRestoreCmd())
	cmd.AddCommand(NewCreateCmd())

	return cmd
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
