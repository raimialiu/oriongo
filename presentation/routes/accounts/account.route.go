package accounts

import (
	"net/http"
	"oriongo/presentation/routes"

	"github.com/labstack/echo/v4"
)

type AccountRouter struct {
	basePath string
	routes   []routes.BaseRouter
}

func AccountRoutes(prefix string) []routes.BaseRouter {
	a := &AccountRouter{
		basePath: prefix,
		routes:   []routes.BaseRouter{},
	}

	a.CreateNewAccount(func(c echo.Context) error {
		return c.JSON(http.StatusCreated, map[string]string{})
	})

	a.GetAccountById(func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(http.StatusOK, id)
	})
	return a.routes
}

func (a AccountRouter) Routes() []routes.BaseRouter {
	return a.routes
}

func (a AccountRouter) GetAccountById(handler echo.HandlerFunc) {
	a.routes = append(a.routes, *routes.NewBaseRouter(
		a.basePath,
		routes.RoutePath{
			Path:    "/:encodedKey",
			Method:  echo.GET,
			Handler: handler,
		},
	))
}

func (a AccountRouter) CreateNewAccount(handler echo.HandlerFunc) {
	a.routes = append(a.routes, *routes.NewBaseRouter(
		a.basePath,
		routes.RoutePath{
			Path:    "",
			Method:  echo.POST,
			Handler: handler,
		},
	))
}
