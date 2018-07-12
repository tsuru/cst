package worker

import (
	"github.com/sirupsen/logrus"
	"github.com/tsuru/cst/db"
	"github.com/tsuru/cst/queue"
	"github.com/tsuru/cst/scan"
	"github.com/tsuru/monsterqueue"
)

// ScanTask implements a monsterqueue.Task interface.
type ScanTask struct {
	Scanners []scan.Scanner
}

// Run executes a scheduled scan over all scanners available.
func (st *ScanTask) Run(job monsterqueue.Job) {

	log := logrus.WithField("job.id", job.ID())

	log.Info("initializing a new job")

	defer log.Info("finishing job")

	scanID := job.Parameters()["id"].(string)
	image := job.Parameters()["image"].(string)

	storage := db.GetStorage()

	storage.UpdateScanStatusByID(scanID, scan.StatusRunning)

	results := make([]scan.Result, len(st.Scanners))

	for index, scanner := range st.Scanners {

		result := scanner.Scan(image)
		results[index] = result

		storage.AppendResultToScanByID(scanID, result)
	}

	storage.UpdateScanStatusByID(scanID, scan.StatusFinished)

	job.Success(results)
}

// Name returns the name of this task.
func (st *ScanTask) Name() string {
	return queue.ScanTaskName
}
