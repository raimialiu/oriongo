package constants

type DbDialect string

const (
	POSTGRES  DbDialect = "postgres"
	MYSQL     DbDialect = "mysql"
	SQLITE    DbDialect = "sqlite3"
	SQLSERVER DbDialect = "sqlserver"
	ORACLE    DbDialect = "oracle"
)
