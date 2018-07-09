package server

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tsuru/cst/api"
)

var (
	webserver api.WebServer

	signalChan = make(chan os.Signal, 1)
)

// New creates an instance of server command.
func New() *cobra.Command {

	serverCmd := &cobra.Command{
		Use:    "server",
		Short:  "Run a web server and listen for requests",
		PreRun: serverCommandPreRun,
		Run:    serverCommandRun,
	}

	serverCmd.Flags().
		String("cert-file", "", "certificate file (required)")

	serverCmd.Flags().
		String("key-file", "", "certificate's private key file (required)")

	serverCmd.Flags().
		IntP("port", "p", 8443, "port to listen")

	serverCmd.MarkFlagRequired("cert-file")
	serverCmd.MarkFlagRequired("key-file")

	viper.BindPFlag("server.cert-file", serverCmd.Flags().Lookup("cert-file"))
	viper.BindPFlag("server.key-file", serverCmd.Flags().Lookup("key-file"))
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))

	return serverCmd
}

func serverCommandPreRun(cmd *cobra.Command, args []string) {

	webserver = &api.SecureWebServer{
		CertFile: viper.GetString("server.cert-file"),
		KeyFile:  viper.GetString("server.key-file"),
		Port:     viper.GetInt("server.port"),
	}
}

func serverCommandRun(cmd *cobra.Command, args []string) {

	// initializes a web server in another thread to be able to handle signals
	go func() {
		if err := webserver.Start(); err != nil {
			logrus.
				WithError(err).
				Info("shutting down the web server")
		}

		signalChan <- os.Interrupt
	}()

	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	signal.Stop(signalChan)

	webserver.Shutdown()
}
