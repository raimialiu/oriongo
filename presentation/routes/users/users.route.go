package users

import (
	"oriongo/internal/origongo"

	"github.com/labstack/echo/v4"
)

func UserRoutes(prefix string, app *origongo.OrionGo) {
	e := app.Host()
	groupRoute := e.Group(prefix)

	groupRoute.GET("/", func(c echo.Context) error {

	})
}
