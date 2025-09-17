package infrastructure

import (
	"errors"
	"fmt"
	"oriongo/internal/common/constants"
	"reflect"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	DbContext struct {
		_config           ConnectionConfig
		_models           []ModelContext
		_db               *gorm.DB
		_isConnected      bool
		_currentModel     *ModelContext
		ConnectionStarted bool
	}

	ModelContext struct {
		Name      string
		TableName string
		Value     interface{}
	}
	ConnectionConfig struct {
		AutoConnect bool
		Host        string
		Port        string
		Database    string
		Username    string
		Password    string
	}
)

func NewDbContext(
	config ConnectionConfig,
) *DbContext {
	var db *gorm.DB
	var isConnected bool = false
	var connectionStarted bool = false
	if config.AutoConnect {
		database, connected := _init(config.Host, config.Username, config.Port, config.Database, config.Password)
		if connected {
			db = database
			isConnected = connected
			connectionStarted = true
		}
	}
	return &DbContext{
		_config:           config,
		_db:               db,
		_isConnected:      isConnected,
		ConnectionStarted: connectionStarted,
	}
}

func _init(host, user, port, dbName, pword string) (*gorm.DB, bool) {
	dsn := fmt.Sprintf("%s:%p@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pword, host, port, dbName)
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

	return db, true
}

func (ctx *DbContext) ConnectionStatus() constants.DbConnectionStatus {
	if ctx._isConnected {
		return constants.CONNECTED
	}

	return constants.CONNECTING
}

func (ctx *DbContext) OpenConnection() *DbContext {
	config := ctx._config
	ctx.ConnectionStarted = true
	db, connected := _init(config.Host, config.Username, config.Port, config.Database, config.Password)
	if connected {
		fmt.Println("Successfully connected to database")
		ctx._isConnected = true
		ctx._db = db
	}

	return ctx
}

func (ctx *DbContext) Add(model interface{}) {

}

func (ctx *DbContext) Model(model interface{}) *DbContext {
	modelType := reflect.TypeOf(model)
	ctx._currentModel = &ModelContext{
		Name:      modelType.Name(),
		TableName: modelType.Name(),
		Value:     model,
	}
	return ctx
}

/*
func (ctx *DbContext) Query(condition interface{}) interface{} {
	if ctx._currentModel == nil {
		return nil
	}

	ctx._db.Table(ctx._currentModel.TableName).Find(condition)
}
*/

func (ctx *DbContext) Where(query interface{}, args ...interface{}) interface{} {
	if ctx._currentModel == nil {
		return nil
	}

	result := ctx._db.Where(query, args...).Find(ctx._currentModel.Value)
	return result
}

func (ctx *DbContext) FirstOrDefault() interface{} {
	if ctx._currentModel == nil {
		return nil
	}

	// context := context.Background()

	result := ctx._db.First(&ctx._currentModel)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	// ctx._currentModel = nil
	return result
}
