package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
)

func registerStatic(e *core.ServeEvent) {
	e.Router.GET("/*", func(c echo.Context) error {
		return c.String(http.StatusNotImplemented, "Static serving not implemented")
	})
}