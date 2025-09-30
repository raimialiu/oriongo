package origongo

import (
	"fmt"
	"net/http"
	"oriongo/internal/common/constants"
	"oriongo/internal/config"
	"oriongo/internal/domain/entities"
	"oriongo/internal/infrastructure"
	"oriongo/presentation/routes"
	"oriongo/presentation/routes/accounts"
	"oriongo/presentation/routes/workspace"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	OrionGo struct {
		_host         *echo.Echo
		Configuration config.OrionConfig
		_viper        *viper.Viper
		route         routes.OrionRoute
		_logger       echo.Logger
		_dbContext    *infrastructure.DbContext
	}

	Configuration struct{}
)

func Instance() *OrionGo {
	e := echo.New()
	app := &OrionGo{
		route: routes.OrionRoute{
			Routes: make([]routes.RoutePath, 0),
		},
		_host: e,
	}

	app.AddConfiguration()
	return app
}

func (o *OrionGo) Use(handler func(next echo.HandlerFunc) echo.HandlerFunc) *OrionGo {
	if o._host == nil {
		panic("OrionGo host is nil")
	}

	o._host.Use(handler)
	return o
}

func (o *OrionGo) AddPingRoute() {
	o.IncludeRoute(*routes.NewBaseRouter("", routes.RoutePath{
		Path: "ping",
		Handler: func(context echo.Context) error {
			return context.JSON(http.StatusOK, "pong")
		},
		Method: http.MethodGet,
	}))
}

func (o *OrionGo) MapPost(path string, handler func(echo.Context) error) *OrionGo {
	o._host.POST(path, handler)
	return o
}

func (o *OrionGo) UseRouting() {
	if len(o.route.Routes) == 0 {
		o._logger.Debug("OrionGo routes is empty")
	}

	for _, route := range o.route.Routes {
		switch route.Method {
		case http.MethodGet:
			o._host.GET(route.Path, route.Handler)
			break
		case http.MethodPost:
			o._host.POST(route.Path, route.Handler)
			break
		case http.MethodPut:
			o._host.PUT(route.Path, route.Handler)
			break
		case http.MethodDelete:
			o._host.DELETE(route.Path, route.Handler)
			break
		}
	}
}

func (o *OrionGo) AddLogging() *OrionGo {
	o._logger = o._host.Logger
	return o
}

func (o *OrionGo) AddRequestHandlers() *OrionGo {

	return o
}

func (o *OrionGo) Set(key string, value interface{}) *OrionGo {
	o._viper.Set(key, value)
	return o
}

func (o *OrionGo) Run() {
	o.Use(middleware.CORS())
	o.UseRouting()
	port := strings.Join([]string{":", strconv.Itoa(o.Configuration.App.PORT)}, "")
	o._host.Logger.Info(fmt.Sprintf("Starting server on port %s", port))
	o._host.Logger.Fatal(o._host.Start(port))
}

func CreateDefaultApp() *OrionGo {
	app := Instance()

	app._logger.Info(strings.Join([]string{"Welcome to ", constants.APPLICATION_NAME}, ""))

	return app
}

func (o *OrionGo) AddRoutes() *OrionGo {
	o.AddPingRoute()
	o.AddWorkspaceRoutes("workspaces")
	o.AddAccountRoutes("accounts")

	return o
}

func (o *OrionGo) IncludeRoute(route routes.BaseRouter) *OrionGo {
	path := strings.Join([]string{fmt.Sprintf("/%s%s", route.BasePath, route.Path.Path)}, "")
	o.route.Routes = append(o.route.Routes, routes.RoutePath{
		Handler: route.Path.Handler,
		Path:    path,
		Method:  route.Path.Method,
	})

	return o
}

func (o *OrionGo) AddAccountRoutes(prefix string) *OrionGo {
	accountRoutes := accounts.AccountRoutes("accounts")
	for _, accountRoute := range accountRoutes {
		o.IncludeRoute(accountRoute)
	}

	return o
}

func (o *OrionGo) AddWorkspaceRoutes(prefix string) {
	workspaceRoutes := workspace.WorkspaceRoutes(prefix, *o._dbContext)
	for _, route := range workspaceRoutes {
		o.IncludeRoute(route)
	}
}

func (app *OrionGo) AddDbContext() *OrionGo {
	var config = infrastructure.ConnectionConfig{
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		Username:    os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASS"),
		Database:    os.Getenv("DB_NAME"),
		AutoMigrate: os.Getenv("DB_AUTO_MIGRATE") == "1",
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: dsn,
		}),
		&gorm.Config{},
	)

	if err != nil {
		fmt.Printf(`failed to connect to database: %s`, err.Error())
		return app
	}

	fmt.Println("Successfully connected to database")
	app._dbContext = infrastructure.NewDbContext(config, db)
	app._dbContext.Migrate()

	//migrationSteps := app.setupMigrationSteps()
	//app._dbContext.Migrator.AddSteps(migrationSteps...)
	//app._dbContext.Migrate()
	return app
}

func (o *OrionGo) DB() *infrastructure.DbContext {
	return o._dbContext
}

func (o *OrionGo) setupMigrationSteps() []infrastructure.MigrationStep {
	var workspace = entities.Workspace{}
	return []infrastructure.MigrationStep{
		{
			Name: "Workspaces",
			Up: func(db *gorm.DB) error {
				return db.AutoMigrate(&workspace)
			},
		},

		{
			Name: "WorkspaceSettings",
			Up: func(db *gorm.DB) error {
				return db.AutoMigrate(&entities.WorkspaceSettings{})
			},
		},
	}

}

func (o *OrionGo) Host() *echo.Echo {
	return o._host
}

func (o *OrionGo) AddConfiguration() {
	v := viper.New()

	env := os.Getenv(constants.ORION_ENVIRONMENT)
	if env == "" {
		env = constants.Development.String()
	}

	o._host.Logger.Debug("Environment " + env)

	configName := strings.Join([]string{"config.", env}, "")
	o._host.Logger.Debug("Config " + configName)

	v.SetConfigType("yaml")
	v.AddConfigPath("./internal/config/paths")
	v.SetConfigName(configName)

	v.AutomaticEnv()
	v.SetEnvPrefix(constants.ENV_PREFIX)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	} else {
		var config config.OrionConfig
		if err := v.Unmarshal(&config); err != nil {
			panic(err)
		}

		o.Configuration = config
		o._viper = v
	}

}
