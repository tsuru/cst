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
)

func TestServerCommandPreRun(t *testing.T) {

	defer func() {
		webserver = nil
	}()

	t.Run(``, func(t *testing.T) {

		assert.Nil(t, webserver)

		viper.Set("server.cert-file", "/path/to/cert.pem")
		viper.Set("server.key-file", "/path/to/key.pem")
		viper.Set("server.port", 443)

		serverCommandPreRun(nil, []string{})

		expected := &api.SecureWebServer{
			CertFile: "/path/to/cert.pem",
			KeyFile:  "/path/to/key.pem",
			Port:     443,
		}

		assert.Equal(t, expected, webserver)
	})
}

func TestServerCommandRun(t *testing.T) {

	t.Run(`When webserver.Start correctly, receives a SIGINT, should stops the webserver gracefully`, func(t *testing.T) {

		webserverIsStarted := false
		webserverIsStopped := false

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

		go serverCommandRun(nil, []string{})

		time.Sleep(time.Second)

		assert.True(t, webserverIsStarted)
		assert.False(t, webserverIsStopped)

		signalChan <- os.Interrupt

		time.Sleep(time.Second)

		assert.True(t, webserverIsStopped)
	})

	t.Run(`When webserver.Start returns an error, should calls webserver.Shutdown internally`, func(t *testing.T) {

		webserverIsStopped := false

		webserver = &api.MockWebServer{
			MockStart: func() error {
				return fmt.Errorf("just another error on web server")
			},

			MockShutdown: func() error {
				webserverIsStopped = true
				return nil
			},
		}

		go serverCommandRun(nil, []string{})

		time.Sleep(time.Second)

		assert.True(t, webserverIsStopped)
	})

	t.Run(`When webserver.Start doesn't hold the execution, should calls webserver.Shutdown internally`, func(t *testing.T) {

		webserverIsStopped := false

		webserver = &api.MockWebServer{
			MockStart: func() error {
				return nil
			},

			MockShutdown: func() error {

				webserverIsStopped = true

				return nil
			},
		}

		go serverCommandRun(nil, []string{})

		time.Sleep(time.Second)

		assert.True(t, webserverIsStopped)
	})
}

func TestNew(t *testing.T) {

	t.Run(`When required args are not assigned, should retuns a error`, func(t *testing.T) {

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
		}

		for _, args := range errorArgs {
			serverCmd := New()

			serverCmd.PreRun = nil
			serverCmd.Run = func(cmd *cobra.Command, args []string) {}

			serverCmd.SetOutput(bytes.NewBufferString(""))
			serverCmd.SetArgs(args)

			assert.Error(t, serverCmd.Execute())
		}
	})

	t.Run(`When all required parameters are defined, should returns no errors`, func(t *testing.T) {

		successfulArgs := [][]string{
			[]string{
				"--cert-file", "/path/to/cert.pem",
				"--key-file", "/path/to/key.pem",
			},
			[]string{
				"--cert-file", "/path/to/cert.pem",
				"--key-file", "/path/to/key.pem",
				"--port", "443",
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
