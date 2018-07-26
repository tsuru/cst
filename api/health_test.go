package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tsuru/cst/db"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {

	t.Run(`When system is unhealthy, should return 500 error`, func(t *testing.T) {

		db.SetStorage(&db.MockStorage{})

		e := echo.New()

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(``))

		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)

		health(context)

		require.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "DOWN", recorder.Body.String())
	})

	t.Run(`When system is working, should return 200 code`, func(t *testing.T) {

		storage := &db.MockStorage{
			MockPing: func() bool {
				return true
			},
		}

		db.SetStorage(storage)

		e := echo.New()

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(``))

		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)

		health(context)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "WORKING", recorder.Body.String())
	})
}
