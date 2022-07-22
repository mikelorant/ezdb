package cmd

import (
	"log"

	"github.com/mikelorant/ezdb2/internal/app"
	"github.com/spf13/cobra"
)

func NewCopyCmd() *cobra.Command {
	var (
		fromContext string
		fromName    string
		toContext   string
		toName      string
	)

	cmd := &cobra.Command{
		Use:   "copy",
		Short: "A brief description of your command",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := app.CopyOptions{
				FromContext: fromContext,
				FromName:    fromName,
				ToContext:   toContext,
				ToName:      toName,
			}

			a, err := app.New()
			if err != nil {
				log.Fatalf("unable to start app: %v", err)
			}
			if err := a.Copy(opts); err != nil {
				log.Fatalf("unable to restore database: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().SortFlags = false
	cmd.Flags().StringVar(&fromContext, "from-context", "", "Database context")
	cmd.Flags().StringVar(&fromName, "from-name", "", "Database name")
	cmd.Flags().StringVar(&toContext, "to-context", "", "Database context")
	cmd.Flags().StringVar(&toName, "to-name", "", "Database name")
	cmd.MarkFlagRequired("to-name")

	return cmd
}
