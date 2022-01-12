package cmd

import "github.com/spf13/cobra"

var (
	databasePath string
	listenAddr   string
)

func Execute() error {
	return NewCommand().Execute()
}

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "shrt",
		Short: "A golinks implementation",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	c.PersistentFlags().StringVar(&databasePath, "database", "routes.db", "Path for SQLite Database")

	c.AddCommand(serveCommand())

	return c
}
