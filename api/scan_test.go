package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tsuru/cst/scan"
	"github.com/tsuru/cst/scan/schd"

	"github.com/labstack/echo"
)

func TestCreateScan(t *testing.T) {

	t.Run(`When payload is empty, should return a bad request response`, func(t *testing.T) {

		e := echo.New()

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(``))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()

		context := e.NewContext(request, recorder)

		err := createScan(context)

		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, context)
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run(`When image value contains only spaces, should return bad request`, func(t *testing.T) {
		e := echo.New()

		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{ "image" : "    " }`))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		recorder := httptest.NewRecorder()

		context := e.NewContext(request, recorder)

		err := createScan(context)

		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, context)
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		}
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

		if assert.NoError(t, createScan(context)) {
			assert.Equal(t, http.StatusCreated, recorder.Code)
		}
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
				return scan.Scan{}, schd.ErrImageHasAlreadyBeenSchedule
			},
		}

		err := createScan(context)

		if assert.NoError(t, err) {
			e.HTTPErrorHandler(err, context)
			assert.Equal(t, http.StatusNoContent, recorder.Code)
		}
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

		if assert.Error(t, err) {
			e.HTTPErrorHandler(err, context)
			assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		}
	})
}
