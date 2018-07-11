package worker

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tsuru/cst/db"
	"github.com/tsuru/cst/db/mongodb"
	"github.com/tsuru/cst/queue"
)

var (
	signalChan = make(chan os.Signal, 1)

	newQueue   = queue.NewQueue
	newStorage = mongodb.NewMongoDB
)

// New creates an instance of worker command.
func New() *cobra.Command {

	workerCmd := &cobra.Command{
		Use:    "worker",
		Short:  "Run worker to fire on scheduled scans",
		PreRun: workerCommandPreRun,
		Run:    workerCommandRun,
	}

	workerCmd.Flags().
		String("database", "", "database URL connection (required)")

	workerCmd.MarkFlagRequired("database")

	viper.BindPFlag("worker.database", workerCmd.Flags().Lookup("database"))

	return workerCmd
}

func workerCommandPreRun(cmd *cobra.Command, args []string) {

	databaseURL := viper.GetString("worker.database")

	storage, err := newStorage(databaseURL)

	if err != nil {
		logrus.WithError(err).Fatal("problem to connect on storage service")
	}

	db.SetStorage(storage)

	q, err := newQueue(databaseURL)

	if err != nil {
		logrus.WithError(err).Fatal("problem to connect on queue service")
	}

	queue.SetQueue(q)
}

func workerCommandRun(cmd *cobra.Command, args []string) {

	q := queue.GetQueue()

	// process the jobs in another thread to be able to handle signals
	go q.ProcessLoop()

	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	signal.Stop(signalChan)

	q.Stop()
	db.GetStorage().Close()
}
