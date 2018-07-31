package server

import (
	"errors"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tsuru/cst/api"
	"github.com/tsuru/cst/db"
	"github.com/tsuru/cst/db/mongodb"
	"github.com/tsuru/cst/queue"
)

var (
	webserver api.WebServer

	signalChan = make(chan os.Signal, 1)

	newQueue   = queue.NewQueue
	newStorage = mongodb.NewMongoDB
)

// New creates an instance of server command.
func New() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:    "server",
		Short:  "Run a web server and listen for requests",
		PreRun: serverCommandPreRun,
		Run:    serverCommandRun,
		Args: func(cmd *cobra.Command, args []string) error {
			if !viper.GetBool("server.insecure") {
				if len(viper.GetString("server.cert-file")) == 0 {
					return errors.New("cert-file is required")
				}
				if len(viper.GetString("server.key-file")) == 0 {
					return errors.New("key-file is required")
				}
			}
			return nil
		},
	}

	serverCmd.Flags().
		String("cert-file", "", "certificate file")

	serverCmd.Flags().
		String("key-file", "", "certificate's private key file")

	serverCmd.Flags().
		IntP("port", "p", 8443, "port to listen")

	serverCmd.Flags().
		String("database", "", "database URL connection (required)")

	serverCmd.Flags().
		Bool("insecure", false, "start server without TLS")

	serverCmd.MarkFlagRequired("database")

	viper.BindPFlag("server.cert-file", serverCmd.Flags().Lookup("cert-file"))
	viper.BindPFlag("server.key-file", serverCmd.Flags().Lookup("key-file"))
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.database", serverCmd.Flags().Lookup("database"))
	viper.BindPFlag("server.insecure", serverCmd.Flags().Lookup("insecure"))

	return serverCmd
}

func serverCommandPreRun(cmd *cobra.Command, args []string) {
	databaseURL := viper.GetString("server.database")

	database, err := newStorage(databaseURL)

	if err != nil {
		logrus.WithError(err).Fatal("problem to connect on storage service")
	}

	db.SetStorage(database)

	q, err := newQueue(databaseURL)

	if err != nil {
		logrus.WithError(err).Fatal("problem to connect on queue service")
	}

	queue.SetQueue(q)

	webserver = &api.SecureWebServer{
		CertFile: viper.GetString("server.cert-file"),
		KeyFile:  viper.GetString("server.key-file"),
		Port:     viper.GetInt("server.port"),
		UseTLS:   !viper.GetBool("server.insecure"),
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
	db.GetStorage().Close()
}
