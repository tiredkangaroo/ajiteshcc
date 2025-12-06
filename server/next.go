package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tiredkangaroo/ajiteshcc/env"
)

func redirectFrontend(c echo.Context, path string) error {
	// note: hard coded
	if env.DefaultEnv.DEBUG {
		path = "http://" + "127.0.0.1" + ":5173" + path
	}
	return c.Redirect(http.StatusFound, path)
}
