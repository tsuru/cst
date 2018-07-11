package api

import (
	"net/http"
	"strings"

	"github.com/tsuru/cst/scan/schd"

	"github.com/labstack/echo"
)

var scheduler schd.Scheduler

func init() {
	scheduler = &schd.DefaultScheduler{}
}

type scanRequest struct {
	Image string `json:"image"`
}

func createScan(ctx echo.Context) error {

	scanRequest := new(scanRequest)

	if err := ctx.Bind(scanRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	imageWithoutSpaces := strings.Replace(scanRequest.Image, " ", "", -1)

	if imageWithoutSpaces == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "the image key is required")
	}

	scanRequest.Image = imageWithoutSpaces

	scan, err := scheduler.Schedule(scanRequest.Image)

	switch err {
	case nil:
		return ctx.JSON(http.StatusCreated, scan)

	case schd.ErrImageHasAlreadyBeenSchedule:
		return ctx.NoContent(http.StatusNoContent)

	default:
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
}
