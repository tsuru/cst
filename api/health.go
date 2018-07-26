package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/tsuru/cst/db"
)

func health(ctx echo.Context) error {

	if db.GetStorage().Ping() {
		return ctx.String(http.StatusOK, "WORKING")
	}

	return ctx.String(http.StatusInternalServerError, "DOWN")
}
