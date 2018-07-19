package worker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tsuru/cst/db"
	"github.com/tsuru/cst/queue"
	"github.com/tsuru/cst/scan"
	"github.com/tsuru/monsterqueue"
)

func TestScanTask_Name(t *testing.T) {

	t.Run(`Ensure name returned by scanTaks.Name method`, func(t *testing.T) {

		st := &ScanTask{}
		assert.Equal(t, queue.ScanTaskName, st.Name())
	})
}

func TestScanTask_Run(t *testing.T) {

	t.Run(`Ensure expected methods ared correctly called`, func(t *testing.T) {

		gotImageOnScanner := ""
		gotResult := scan.Result{}
		gotStatus := scan.Status("")

		wasSuccessful := false

		st := &ScanTask{
			Scanners: []scan.Scanner{
				&scan.MockScanner{
					MockScan: func(image string) scan.Result {

						gotImageOnScanner = image

						return scan.Result{
							Scanner: "mocked-scanner",
						}
					},
				},
			},
		}

		job := queue.MockJob{
			MockParameters: func() monsterqueue.JobParams {
				return monsterqueue.JobParams{
					"id":    "d29b39eb-a5e5-4237-acb4-e7203cd6e2cf",
					"image": "tsuru/cst:latest",
				}
			},

			MockSucess: func(result monsterqueue.JobResult) (bool, error) {

				wasSuccessful = true

				return false, nil
			},
		}

		storage := &db.MockStorage{
			MockUpdateScanStatusByID: func(id string, status scan.Status) error {

				gotStatus = status

				return nil
			},
			MockAppendResultToScanByID: func(id string, result scan.Result) error {

				gotResult = result

				return nil
			},
		}

		db.SetStorage(storage)

		st.Run(job)

		assert.Equal(t, "tsuru/cst:latest", gotImageOnScanner)
		assert.Equal(t, "mocked-scanner", gotResult.Scanner)
		assert.Equal(t, scan.StatusFinished, gotStatus)
		assert.True(t, wasSuccessful)
	})

	t.Run(`When storage returns any error on UpdateScanStatusByID method with scan.StatusRunning param, should abort execution and call the job.Error method`, func(t *testing.T) {

		gotJobError := false

		storage := &db.MockStorage{
			MockUpdateScanStatusByID: func(id string, status scan.Status) error {

				if status == scan.StatusRunning {
					return errors.New("just another error on storage")
				}

				return nil
			},
		}

		db.SetStorage(storage)

		job := queue.MockJob{
			MockParameters: func() monsterqueue.JobParams {
				return monsterqueue.JobParams{
					"id":    "d29b39eb-a5e5-4237-acb4-e7203cd6e2cf",
					"image": "tsuru/cst:latest",
				}
			},

			MockError: func(err error) (bool, error) {
				gotJobError = true

				return false, err
			},
		}

		st := &ScanTask{}

		st.Run(job)

		assert.True(t, gotJobError)
	})

	t.Run(`When storage returns any error on UpdateScanStatusByID method with scan.StatusFinished param, should abort execution and call the job.Error method`, func(t *testing.T) {
		gotJobError := false

		storage := &db.MockStorage{
			MockUpdateScanStatusByID: func(id string, status scan.Status) error {

				if status == scan.StatusFinished {
					return errors.New("just another error on storage")
				}

				return nil
			},
		}

		db.SetStorage(storage)

		job := queue.MockJob{
			MockParameters: func() monsterqueue.JobParams {
				return monsterqueue.JobParams{
					"id":    "d29b39eb-a5e5-4237-acb4-e7203cd6e2cf",
					"image": "tsuru/cst:latest",
				}
			},

			MockError: func(err error) (bool, error) {
				gotJobError = true

				return false, err
			},
		}

		st := &ScanTask{}

		st.Run(job)

		assert.True(t, gotJobError)
	})
}
