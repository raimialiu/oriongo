package origongo

import (
	"fmt"
	"net/http"
	"oriongo/internal/common/constants"
	"oriongo/internal/config"
	"oriongo/internal/infrastructure"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	OrionGo struct {
		_host         *echo.Echo
		Configuration config.OrionConfig
		_viper        *viper.Viper
		_logger       echo.Logger
		_dbContext    *gorm.DB
	}

	Configuration struct{}
)

func Instance() *OrionGo {
	e := echo.New()
	app := &OrionGo{
		_host: e,
	}

	app.AddConfiguration()
	return app
}

func (app *OrionGo) Info(message string) *OrionGo {
	app._host.Logger.Info(message)
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
	o._host.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "pong")
	})
}

func (o *OrionGo) MapPost(path string, handler func(echo.Context) error) *OrionGo {
	o._host.POST(path, handler)
	return o
}

func (o *OrionGo) AddLogging() *OrionGo {
	o._logger = o._host.Logger

	return o
}

func (o *OrionGo) AddControllers() *OrionGo {

	return o
}

func (o *OrionGo) Set(key string, value interface{}) *OrionGo {
	o._viper.Set(key, value)
	return o
}

func (o *OrionGo) Run() {
	port := strings.Join([]string{":", strconv.Itoa(o.Configuration.App.PORT)}, "")
	o._host.Logger.Info(fmt.Sprintf("Starting server on port %s", port))
	o._host.Logger.Fatal(o._host.Start(port))
}

func CreateDefaultApp() *OrionGo {
	app := Instance()
	app.AddPingRoute()

	print(strings.Join([]string{"Welcome to ", constants.APPLICATION_NAME}, ""))

	return app
}

func (app *OrionGo) AddDbContext(config infrastructure.ConnectionConfig) *OrionGo {
	dsn := fmt.Sprintf("%s:%p@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Database)
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: dsn,
		}),
		&gorm.Config{},
	)

	if err != nil {
		fmt.Errorf("Failed to connect to database: %v", err)
	}

	fmt.Println("Successfully connected to database")
	app._dbContext = db
	return app
}

func (o *OrionGo) DB() *gorm.DB {
	return o._dbContext
}

func (o *OrionGo) Host() *echo.Echo {
	return o._host
}

func (o *OrionGo) AddConfiguration() {
	v := viper.New()

	env := os.Getenv(constants.GetEnvName(constants.ORION_ENVIRONMENT))
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
