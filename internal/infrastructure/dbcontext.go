package infrastructure

import (
	"fmt"
	"oriongo/internal/common/constants"
	"oriongo/internal/domain/entities"

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
		RunMigration      bool
		Migrator          *GormMigrator
	}

	ModelContext struct {
		Name      string
		TableName string
		Value     interface{}
	}
	ConnectionConfig struct {
		AutoConnect bool
		AutoMigrate bool
		Host        string
		Port        string
		Database    string
		Username    string
		Password    string
		Dialect     string
	}
)

func NewDbContext(
	config ConnectionConfig,
	gormDB *gorm.DB,
) *DbContext {
	var db *gorm.DB
	var isConnected bool = false
	var connectionStarted bool = false

	if gormDB != nil {
		return &DbContext{
			_config:           config,
			_db:               gormDB,
			_isConnected:      isConnected,
			ConnectionStarted: connectionStarted,
			Migrator:          NewMigrator(),
			RunMigration:      config.AutoMigrate,
		}
	}

	if config.AutoConnect {
		database, connected := _init(config.Host, config.Username, config.Port, config.Database, config.Password)
		if connected {
			db = database
			isConnected = connected
			connectionStarted = true
		}
	}

	migrator := NewMigrator()
	return &DbContext{
		_config:           config,
		_db:               db,
		_isConnected:      isConnected,
		ConnectionStarted: connectionStarted,
		Migrator:          migrator,
		RunMigration:      config.AutoMigrate,
	}
}

func (db *DbContext) Migrate() {
	if !db.RunMigration {
		return
	}

	//db.Migrator.Up(db._db)
	migrateError := db._db.AutoMigrate(
		&entities.Workspace{},
		&entities.WorkspaceSettings{},
		&entities.Organization{},
		&entities.OrganizationSettings{},
		&entities.OrganizationUser{},
	)
	if migrateError != nil {
		fmt.Println(migrateError)
	}
}

func _init(host, user, port, dbName, pword string) (*gorm.DB, bool) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pword, host, port, dbName)
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: dsn,
		}),
		&gorm.Config{},
	)

	if err != nil {
		fmt.Println()
		fmt.Printf("Failed to connect to database: %v", err)
		return nil, false
	}

	fmt.Println("Successfully connected to database")
	return db, true
}

func (d *DbContext) DB() *gorm.DB {
	return (d._db)
}

func (ctx *DbContext) ConnectionStatus() constants.DbConnectionStatus {
	if ctx._isConnected {
		return constants.CONNECTED
	}

	return constants.CONNECTING
}
