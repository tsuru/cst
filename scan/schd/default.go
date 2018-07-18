package schd

import (
	"time"

	"github.com/tsuru/monsterqueue"

	uuid "github.com/satori/go.uuid"
	"github.com/tsuru/cst/db"
	"github.com/tsuru/cst/queue"
	"github.com/tsuru/cst/scan"
)

// DefaultScheduler implements a Scheduler interface.
type DefaultScheduler struct{}

// Schedule registers a new analysis of a given image. It returns the complete
// entry of scan if successful else retuns an error instance to indicate the
// wrong state.
func (ds *DefaultScheduler) Schedule(image string) (scan.Scan, error) {

	storage := db.GetStorage()

	if storage.HasScheduledScanByImage(image) {
		return scan.Scan{}, ErrImageHasAlreadyBeenScheduled
	}

	newScan := scan.Scan{
		ID:        uuid.NewV4().String(),
		Status:    scan.StatusScheduled,
		Image:     image,
		CreatedAt: time.Now(),
		Result:    []scan.Result{},
	}

	if err := storage.Save(newScan); err != nil {
		return scan.Scan{}, err
	}

	enqueueScan(newScan)

	return newScan, nil
}

func enqueueScan(scan scan.Scan) {

	q := queue.GetQueue()

	params := monsterqueue.JobParams{
		"id":    scan.ID,
		"image": scan.Image,
	}

	q.Enqueue(queue.ScanTaskName, params)
}
