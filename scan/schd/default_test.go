package schd

import (
	"errors"
	"testing"

	"github.com/tsuru/monsterqueue"

	"github.com/stretchr/testify/assert"

	"github.com/tsuru/cst/db"
	"github.com/tsuru/cst/queue"
	"github.com/tsuru/cst/scan"
)

func TestDefaultScheduler_Schedule(t *testing.T) {

	defer func() {
		db.SetStorage(nil)
		queue.SetQueue(nil)
	}()

	t.Run(`When scanning an image is already scheduled, should return ErrImageHasAlreadyBeenScheduled error`, func(t *testing.T) {

		queue.SetQueue(&queue.MockQueue{})

		storage := &db.MockStorage{
			MockHasScheduledScanByImage: func(img string) bool {
				return true
			},
		}

		db.SetStorage(storage)

		ds := &DefaultScheduler{}
		_, err := ds.Schedule("tsuru/cst:latest")

		if assert.Error(t, err) {
			assert.Equal(t, ErrImageHasAlreadyBeenSchedule, err)
		}
	})

	t.Run(`When scanning a new image, should return a new instance of Scan and no errors`, func(t *testing.T) {

		queue.SetQueue(&queue.MockQueue{})

		storage := &db.MockStorage{
			MockHasScheduledScanByImage: func(img string) bool {
				return false
			},
		}

		db.SetStorage(storage)

		ds := &DefaultScheduler{}

		newScan, err := ds.Schedule("tsuru/cst:latest")

		if assert.NoError(t, err) {
			assert.NotEmpty(t, newScan.ID)
			assert.Equal(t, string(scan.StatusScheduled), string(newScan.Status))
			assert.Equal(t, "tsuru/cst:latest", newScan.Image)
		}
	})

	t.Run(`When scanning a new image, when storage returns error on storage.Save method, shoul return an error`, func(t *testing.T) {

		queue.SetQueue(&queue.MockQueue{})

		storage := &db.MockStorage{
			MockHasScheduledScanByImage: func(img string) bool {
				return false
			},

			MockSave: func(s scan.Scan) error {
				return errors.New(`just another error on storage`)
			},
		}

		db.SetStorage(storage)

		ds := &DefaultScheduler{}

		_, err := ds.Schedule("tsuru/cst:latest")

		assert.Error(t, err)
	})

	t.Run(`Ensure queue.Enqueue is called with expected params`, func(t *testing.T) {

		gotTaskName := ""
		gotParams := monsterqueue.JobParams{}

		q := &queue.MockQueue{
			MockEnqueue: func(task string, params monsterqueue.JobParams) (monsterqueue.Job, error) {

				gotTaskName = task
				gotParams = params

				return nil, nil
			},
		}

		queue.SetQueue(q)

		newScan := scan.Scan{
			ID:    "d29b39eb-a5e5-4237-acb4-e7203cd6e2cf",
			Image: "tsuru/cst:latest",
		}

		enqueueScan(newScan)

		assert.Equal(t, queue.ScanTaskName, gotTaskName)
		assert.Equal(t, newScan.ID, gotParams["id"])
		assert.Equal(t, newScan.Image, gotParams["image"])
	})
}
