package worker

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/tsuru/cst/db"
	"github.com/tsuru/cst/db/mongodb"
	"github.com/tsuru/cst/queue"
	"github.com/tsuru/monsterqueue"
)

func TestWorkerCommandPreRun(t *testing.T) {
	t.Run(`Ensure newQueue and newStorage are called with expected param`, func(t *testing.T) {
		gotStorageURL := ""
		gotQueueURL := ""

		newQueue = func(url string) (monsterqueue.Queue, error) {
			gotQueueURL = url

			return nil, nil
		}

		newStorage = func(url string) (*mongodb.MongoDB, error) {
			gotStorageURL = url

			return nil, nil
		}

		viper.Set("worker.database", "mongodb://localhost/")

		workerCommandPreRun(nil, []string{})

		assert.Equal(t, gotQueueURL, viper.Get("worker.database"))
		assert.Equal(t, gotStorageURL, viper.Get("worker.database"))
	})
}

func TestWorkerCommandRun(t *testing.T) {
	t.Run(`Ensure queue.ProcessLoop is called, when receive a INT singnal should calls storage.Close`, func(t *testing.T) {
		hasCalledProcessLoop := false
		hasCalledStorageClose := false

		q := queue.MockQueue{
			MockProcessLoop: func() {
				hasCalledProcessLoop = true
			},
		}

		queue.SetQueue(q)

		storage := &db.MockStorage{
			MockClose: func() {
				hasCalledStorageClose = true
			},
		}

		db.SetStorage(storage)

		go workerCommandRun(nil, []string{})

		time.Sleep(time.Second)

		assert.True(t, hasCalledProcessLoop)
		assert.False(t, hasCalledStorageClose)

		signalChan <- os.Interrupt

		time.Sleep(time.Second)

		assert.True(t, hasCalledStorageClose)
	})
}

func TestNew(t *testing.T) {
	t.Run(`When required args are not assigned, should retuns a error`, func(t *testing.T) {
		errorArgs := [][]string{
			[]string{},
			[]string{
				"--database", "mongodb://localhost:27018/mydb",
			},
			[]string{
				"--clair-address", "https://clair.tld:6060",
			},
		}

		for _, args := range errorArgs {
			workerCmd := New()

			workerCmd.PreRun = nil
			workerCmd.Run = func(cmd *cobra.Command, args []string) {}

			workerCmd.SetOutput(bytes.NewBufferString(""))
			workerCmd.SetArgs(args)

			assert.Error(t, workerCmd.Execute())
		}
	})

	t.Run(`When all required parameters are defined, should returns no errors`, func(t *testing.T) {
		successfulArgs := [][]string{
			[]string{
				"--database", "mongodb://127.0.0.1:27017/",
				"--clair-address", "https://clair.top.level.domain:8443",
			},
		}

		for _, args := range successfulArgs {
			workerCmd := New()

			workerCmd.PreRun = nil
			workerCmd.Run = func(cmd *cobra.Command, args []string) {}

			workerCmd.SetOutput(bytes.NewBufferString(""))
			workerCmd.SetArgs(args)

			assert.NoError(t, workerCmd.Execute())
		}
	})
}
