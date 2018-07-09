package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tsuru/cst/cmd/server"
)

// Execute calls the CLI entrypoint. If fired subcommand returns any error, so
// logs the error on standard output and exit with non-zero code.
func Execute() {

	err := newRootCommand().Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}

func newRootCommand() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:  "cst",
		Args: cobra.MinimumNArgs(1),
	}

	rootCmd.AddCommand(server.New())

	return rootCmd
}
