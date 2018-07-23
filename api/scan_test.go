package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsuru/cst/scan"
	schd "github.com/tsuru/cst/scan/scheduler"
)

func TestCreateScan(t *testing.T) {
	t.Run(`When payload is empty, should return a bad request response`, func(t *testing.T) {
		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(``))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)
		err := createScan(context)

		require.Error(t, err)
		e.HTTPErrorHandler(err, context)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run(`When image value contains only spaces, should return bad request`, func(t *testing.T) {
		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{ "image" : "    " }`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)
		err := createScan(context)

		require.Error(t, err)
		e.HTTPErrorHandler(err, context)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run(`When payload is OK, should return created status code`, func(t *testing.T) {
		requestBody := `{ "image": "tsuru/cst:latest" }`

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)
		scheduler = &schd.MockScheduler{
			MockSchedule: func(image string) (scan.Scan, error) {
				return scan.Scan{
					ID:     `2c5a5f48-9801-40ef-8d81-cd4a0f9c0ee2`,
					Status: scan.StatusScheduled,
					Image:  `tsuru/cst:latest`,
				}, nil
			},
		}

		require.NoError(t, createScan(context))
		assert.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run(`When scheduler returns an ErrImageHasAlreadyBeenScheduled error, should return NoContent status code`, func(t *testing.T) {
		requestBody := `{ "image": "tsuru/cst:latest" }`

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)
		scheduler = &schd.MockScheduler{
			MockSchedule: func(image string) (scan.Scan, error) {
				return scan.Scan{}, schd.ErrImageHasAlreadyBeenScheduled
			},
		}
		err := createScan(context)

		require.NoError(t, err)
		e.HTTPErrorHandler(err, context)
		assert.Equal(t, http.StatusNoContent, recorder.Code)
	})

	t.Run(`When payload is OK and scheduler return an error, should return Internal Server Error`, func(t *testing.T) {
		requestBody := `{ "image": "tsuru/cst:latest" }`

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)
		scheduler = &schd.MockScheduler{
			MockSchedule: func(image string) (scan.Scan, error) {
				return scan.Scan{}, errors.New("something went wrong")
			},
		}
		err := createScan(context)

		require.Error(t, err)
		e.HTTPErrorHandler(err, context)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run(`When payload contains endcustomdata, decodes it and posts a scan`, func(t *testing.T) {
		requestBody := `{"endcustomdata": "IQAAAAJpbWFnZQARAAAAdHN1cnUvY3N0OmxhdGVzdAAA", "image": ""}`

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)
		scheduler = &schd.MockScheduler{
			MockSchedule: func(image string) (scan.Scan, error) {
				require.Equal(t, "tsuru/cst:latest", image)
				return scan.Scan{}, nil
			},
		}
		err := createScan(context)

		require.Nil(t, err)
		e.HTTPErrorHandler(err, context)
		require.Equal(t, http.StatusCreated, recorder.Code)
	})

	t.Run(`When payload contains bad endcustomdata, returns bad request`, func(t *testing.T) {
		requestBody := `{"endcustomdata": "invalid_data"}`

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)
		scheduler = &schd.MockScheduler{
			MockSchedule: func(image string) (scan.Scan, error) {
				require.Fail(t, "shouldn't schedule for invalid data")
				return scan.Scan{}, errors.New("invalid metadata")
			},
		}
		err := createScan(context)

		require.Error(t, err)
		e.HTTPErrorHandler(err, context)
		require.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run(`When payload contains both image and endcustomdata, ignores endcustomdata`, func(t *testing.T) {
		requestBody := `{"image": "tsuru/cst:latest", "endcustomdata": "invalid_data"}`

		e := echo.New()
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder := httptest.NewRecorder()
		context := e.NewContext(request, recorder)
		scheduler = &schd.MockScheduler{
			MockSchedule: func(image string) (scan.Scan, error) {
				require.Equal(t, "tsuru/cst:latest", image)
				return scan.Scan{}, nil
			},
		}
		err := createScan(context)

		require.Nil(t, err)
		e.HTTPErrorHandler(err, context)
		require.Equal(t, http.StatusCreated, recorder.Code)
	})
}
