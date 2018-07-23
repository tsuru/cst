package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
	schd "github.com/tsuru/cst/scan/scheduler"
)

var scheduler schd.Scheduler

func init() {
	scheduler = &schd.DefaultScheduler{}
}

type scanRequest struct {
	Image string `json:"image"`
}

type tsuruEvent struct {
	scanRequest
	EndCustomData string `json:"endcustomdata"`
}

func createScan(ctx echo.Context) error {
	scanRequest, err := loadScanRequestFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
	case schd.ErrImageHasAlreadyBeenScheduled:
		return ctx.NoContent(http.StatusNoContent)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
}

func loadScanRequestFromContext(ctx echo.Context) (*scanRequest, error) {
	evt := tsuruEvent{}
	err := ctx.Bind(&evt)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest)
	}
	scanRequest := &scanRequest{}
	if evt.EndCustomData != "" {
		img, err := unmarshalImage(evt)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		scanRequest.Image = img
	} else {
		scanRequest.Image = evt.scanRequest.Image
	}
	return scanRequest, nil
}

func unmarshalImage(evt tsuruEvent) (string, error) {
	b64, err := base64.StdEncoding.DecodeString(evt.EndCustomData)
	if err != nil {
		return "", err
	}
	var v interface{}
	err = bson.Unmarshal(b64, &v)
	if err != nil {
		return "", err
	}
	m, ok := v.(bson.M)
	if !ok {
		return "", fmt.Errorf("invalid metadata")
	}
	img, ok := m["image"].(string)
	if !ok {
		return "", fmt.Errorf("invalid metadata")
	}
	return img, nil
}
