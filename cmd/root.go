package cmd

import "github.com/spf13/cobra"

const (
	_headerClientID     = "CF-Access-Client-Id"
	_headerClientSecret = "CF-Access-Client-Secret"
)

var (
	databasePath string
	serverAddr   string
	localOp      bool
	clientID     string
	clientSecret string
	confPath     string
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

	cobra.OnInitialize(initConfig)

	c.PersistentFlags().StringVar(&confPath, "config", "", "Path to the config file")

	c.AddCommand(serveCommand(), createCommand(), deleteCommand(), listCommand(), openCommand())

	return c
}
