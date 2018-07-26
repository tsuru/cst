package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tsuru/cst/db"
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

func TestShowScans(t *testing.T) {

	t.Run(`When there are no scans for a given image, should return 200 and an empty slice`, func(t *testing.T) {
		storage := &db.MockStorage{
			MockGetScansByImage: func(image string) ([]scan.Scan, error) {
				return []scan.Scan{}, nil
			},
		}

		db.SetStorage(storage)

		e := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(``))
		recorder := httptest.NewRecorder()

		context := e.NewContext(request, recorder)

		context.SetPath("/v1/scan/:image")
		context.SetParamNames("image")
		context.SetParamValues(url.PathEscape("tsuru/cst:latest"))

		err := showScans(context)

		require.NoError(t, err)
		e.HTTPErrorHandler(err, context)
		assert.Equal(t, http.StatusOK, recorder.Code)

		assert.JSONEq(t, `[]`, recorder.Body.String())
	})

	t.Run(`When image param is bad URL encoded, should return 400 status code`, func(t *testing.T) {
		e := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(``))
		recorder := httptest.NewRecorder()

		context := e.NewContext(request, recorder)

		context.SetPath("/v1/scan/:image")
		context.SetParamNames("image")
		context.SetParamValues("badUrlEncoded%ss")

		err := showScans(context)

		require.Error(t, err)
		e.HTTPErrorHandler(err, context)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run(`When storage returns any error, should return 500 status code`, func(t *testing.T) {
		storage := &db.MockStorage{
			MockGetScansByImage: func(image string) ([]scan.Scan, error) {
				return []scan.Scan{}, errors.New("just another error on storage")
			},
		}

		db.SetStorage(storage)

		e := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(``))
		recorder := httptest.NewRecorder()

		context := e.NewContext(request, recorder)

		context.SetPath("/v1/scan/:image")
		context.SetParamNames("image")
		context.SetParamValues(url.PathEscape("tsuru/cst:latest"))

		err := showScans(context)

		require.Error(t, err)
		e.HTTPErrorHandler(err, context)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})

	t.Run(`When storage correctly returns scans, should return 200 status code and scans on body`, func(t *testing.T) {

		expectedScans := []scan.Scan{
			scan.Scan{
				ID:    "1",
				Image: "tsuru/cst:latest",
			},
		}

		storage := &db.MockStorage{
			MockGetScansByImage: func(image string) ([]scan.Scan, error) {
				return expectedScans, nil
			},
		}

		db.SetStorage(storage)

		e := echo.New()

		request := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(``))
		recorder := httptest.NewRecorder()

		context := e.NewContext(request, recorder)

		context.SetPath("/v1/scan/:image")
		context.SetParamNames("image")
		context.SetParamValues(url.PathEscape("tsuru/cst:latest"))

		err := showScans(context)

		require.NoError(t, err)
		e.HTTPErrorHandler(err, context)
		assert.Equal(t, http.StatusOK, recorder.Code)

		expectedScansJSON, _ := json.Marshal(expectedScans)

		assert.JSONEq(t, string(expectedScansJSON), recorder.Body.String())
	})
}
