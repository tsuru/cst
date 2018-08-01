package server

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/tsuru/cst/api"
	"github.com/tsuru/cst/db"
	"github.com/tsuru/cst/db/mongodb"
	"github.com/tsuru/monsterqueue"
)

func TestServerCommandPreRun(t *testing.T) {
	oldNewQueue := newQueue
	oldNewStorage := newStorage

	defer func() {
		webserver = nil

		newQueue = oldNewQueue
		newStorage = oldNewStorage
		viper.Reset()
	}()

	t.Run(`Ensure WebServer is created with expected params`, func(t *testing.T) {
		newQueue = func(url string) (monsterqueue.Queue, error) {
			return nil, nil
		}

		newStorage = func(url string) (*mongodb.MongoDB, error) {
			return nil, nil
		}

		assert.Nil(t, webserver)

		viper.Set("server.cert-file", "/path/to/cert.pem")
		viper.Set("server.key-file", "/path/to/key.pem")
		viper.Set("server.port", 443)

		serverCommandPreRun(nil, []string{})

		expected := &api.SecureWebServer{
			CertFile: "/path/to/cert.pem",
			KeyFile:  "/path/to/key.pem",
			Port:     443,
			UseTLS:   true,
		}

		assert.Equal(t, expected, webserver)
	})

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

		viper.Set("server.database", "mongodb://localhost/")

		serverCommandPreRun(nil, []string{})

		assert.Equal(t, gotQueueURL, viper.Get("server.database"))
		assert.Equal(t, gotStorageURL, viper.Get("server.database"))
	})
}

func TestServerCommandRun(t *testing.T) {
	t.Run(`When webserver.Start correctly, receives a SIGINT, should stops the webserver and storage gracefully`, func(t *testing.T) {
		webserverIsStarted := false
		webserverIsStopped := false

		storageIsStopped := false

		webserver = &api.MockWebServer{
			MockStart: func() error {
				webserverIsStarted = true

				// used to hold the execution (like a webserver running)
				time.Sleep(time.Minute)

				return nil
			},

			MockShutdown: func() error {
				webserverIsStopped = true
				return nil
			},
		}

		storage := &db.MockStorage{
			MockClose: func() {
				storageIsStopped = true
			},
		}

		db.SetStorage(storage)

		go serverCommandRun(nil, []string{})

		time.Sleep(time.Second)

		assert.True(t, webserverIsStarted)
		assert.False(t, webserverIsStopped)

		signalChan <- os.Interrupt

		time.Sleep(time.Second)

		assert.True(t, webserverIsStopped)
		assert.True(t, storageIsStopped)
	})

	t.Run(`When webserver.Start returns an error, should calls webserver.Shutdown and storage.Close internally`, func(t *testing.T) {
		webserverIsStopped := false

		storageIsStopped := false

		webserver = &api.MockWebServer{
			MockStart: func() error {
				return fmt.Errorf("just another error on web server")
			},

			MockShutdown: func() error {
				webserverIsStopped = true
				return nil
			},
		}

		storage := &db.MockStorage{
			MockClose: func() {
				storageIsStopped = true
			},
		}

		db.SetStorage(storage)

		go serverCommandRun(nil, []string{})

		time.Sleep(time.Second)

		assert.True(t, webserverIsStopped)
		assert.True(t, storageIsStopped)
	})

	t.Run(`When webserver.Start doesn't hold the execution, should calls webserver.Shutdown and storage.Close internally`, func(t *testing.T) {
		webserverIsStopped := false

		storageIsStopped := false

		webserver = &api.MockWebServer{
			MockStart: func() error {
				return nil
			},

			MockShutdown: func() error {
				webserverIsStopped = true

				return nil
			},
		}

		storage := &db.MockStorage{
			MockClose: func() {
				storageIsStopped = true
			},
		}

		db.SetStorage(storage)

		go serverCommandRun(nil, []string{})

		time.Sleep(time.Second)

		assert.True(t, webserverIsStopped)
		assert.True(t, storageIsStopped)
	})
}

func TestNew(t *testing.T) {
	t.Run(`When required args are not assigned, should return an error`, func(t *testing.T) {
		errorArgs := [][]string{
			[]string{},
			[]string{
				"--cert-file", "/path/to/cert.pem",
			},
			[]string{
				"--key-file", "/path/to/key.pem",
			},
			[]string{
				"--port", "8080",
			},
			[]string{
				"--unknown-arg", "unknown-var",
			},
			[]string{
				"--database", "dbhost",
			},
		}

		for _, args := range errorArgs {
			serverCmd := New()

			serverCmd.PreRun = nil
			serverCmd.Run = func(cmd *cobra.Command, args []string) {}

			serverCmd.SetOutput(bytes.NewBufferString(""))
			serverCmd.SetArgs(args)

			assert.Error(t, serverCmd.Execute(), "server should have returned an error with args: ", args)
		}
	})

	t.Run(`When all required parameters are defined, should return no errors`, func(t *testing.T) {
		successfulArgs := [][]string{
			[]string{
				"--cert-file", "/path/to/cert.pem",
				"--key-file", "/path/to/key.pem",
				"--database", "mongodb://127.0.0.1:27017/",
			},
			[]string{
				"--cert-file", "/path/to/cert.pem",
				"--key-file", "/path/to/key.pem",
				"--port", "443",
				"--database", "mongodb://127.0.0.1:27017/",
			},
			[]string{
				"--database", "mongodb://127.0.0.1:27017/",
				"--insecure",
			},
		}

		for _, args := range successfulArgs {
			serverCmd := New()

			serverCmd.PreRun = nil
			serverCmd.Run = func(cmd *cobra.Command, args []string) {}

			serverCmd.SetOutput(bytes.NewBufferString(""))
			serverCmd.SetArgs(args)

			assert.NoError(t, serverCmd.Execute())
		}
	})
}
