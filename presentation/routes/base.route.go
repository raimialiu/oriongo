package routes

import (
	"github.com/labstack/echo/v4"
)

type (
	BaseRouter struct {
		BasePath string
		Path     RoutePath
	}

	RoutePath struct {
		Prefix  string
		Path    string
		Method  string
		Handler echo.HandlerFunc
	}

	OrionRoute struct {
		Routes []RoutePath
	}

	RouteRegistry interface{}
)

func NewBaseRouter(basePath string, path RoutePath) *BaseRouter {
	return &BaseRouter{
		BasePath: basePath,
		Path:     path,
	}
}
