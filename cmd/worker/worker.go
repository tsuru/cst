package worker

import (
	"os"
	"os/signal"
	"time"

	"github.com/tsuru/monsterqueue"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tsuru/cst/db"
	"github.com/tsuru/cst/db/mongodb"
	"github.com/tsuru/cst/queue"
	"github.com/tsuru/cst/scan"
	"github.com/tsuru/cst/scan/worker"
)

var (
	scanTask monsterqueue.Task

	signalChan = make(chan os.Signal, 1)

	newQueue   = queue.NewQueue
	newStorage = mongodb.NewMongoDB
)

// New creates an instance of worker command.
func New() *cobra.Command {

	workerCmd := &cobra.Command{
		Use:    "worker",
		Short:  "Run worker to analyze scheduled scans",
		PreRun: workerCommandPreRun,
		Run:    workerCommandRun,
	}

	workerCmd.Flags().
		String("database", "", "database URL connection (required)")

	workerCmd.Flags().
		String("clair-address", "", "CoresOS Clair address (required)")

	workerCmd.MarkFlagRequired("database")
	workerCmd.MarkFlagRequired("clair-address")

	viper.BindPFlag("worker.database", workerCmd.Flags().Lookup("database"))
	viper.BindPFlag("worker.clair.address", workerCmd.Flags().Lookup("clair-address"))

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

	clair := &scan.Clair{
		Address: viper.GetString("worker.clair.address"),
		Name:    "clair",
		Timeout: time.Minute,
	}

	scanTask = &worker.ScanTask{
		Scanners: []scan.Scanner{
			clair,
		},
	}
}

func workerCommandRun(cmd *cobra.Command, args []string) {

	q := queue.GetQueue()

	q.RegisterTask(scanTask)

	// process the jobs in another thread to be able to handle signals
	go q.ProcessLoop()

	signal.Notify(signalChan, os.Interrupt)

	<-signalChan
	signal.Stop(signalChan)

	q.Stop()
	db.GetStorage().Close()
}
