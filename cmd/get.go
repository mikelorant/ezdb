package cmd

import (
	"log"
	"strings"

	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "A brief description of your command",
	}

	cmd.AddCommand(NewGetSessionCmd())
	cmd.AddCommand(NewGetVariableCmd())

	return cmd
}

func NewGetSessionCmd() *cobra.Command {
	var context string

	cmd := &cobra.Command{
		Use:     "session",
		Aliases: []string{"sessions"},
		Short:   "A brief description of your command",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.GetSessionOptions{
				Context: context,
			}
			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.GetSession(opts); err != nil {
				log.Fatalf("unable to get session: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")

	return cmd
}

func NewGetVariableCmd() *cobra.Command {
	var context string

	cmd := &cobra.Command{
		Use:     "variable",
		Aliases: []string{"variables"},
		Short:   "A brief description of your command",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.GetVariableOptions{
				Context:  context,
				Variable: strings.Join(args, ""),
			}
			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.GetVariable(opts); err != nil {
				log.Fatalf("unable to get variable: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&context, "context", "", "Database context")

	return cmd
}
