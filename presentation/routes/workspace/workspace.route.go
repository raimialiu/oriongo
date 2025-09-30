package workspace

import (
	"oriongo/internal/infrastructure"
	"oriongo/presentation/handlers"
	"oriongo/presentation/routes"

	"github.com/labstack/echo/v4"
)

type WorkspaceRoute struct {
	dbContext *infrastructure.DbContext
	basePath  string
	routes    []routes.BaseRouter
}

func WorkspaceRoutes(prefix string, db infrastructure.DbContext) []routes.BaseRouter {
	a := &WorkspaceRoute{
		dbContext: &db,
		basePath:  prefix,
		routes:    []routes.BaseRouter{},
	}

	a.CreateNewWorkspace()
	return a.routes
}

func (a *WorkspaceRoute) CreateNewWorkspace() {
	workspaceHandlers := handlers.NewWorkspaceHandler(*a.dbContext)
	a.routes = append(a.routes, *routes.NewBaseRouter(a.basePath,
		routes.RoutePath{
			Path:    "",
			Handler: workspaceHandlers.CreateWorkspace,
			Method:  echo.POST,
		}))
}
