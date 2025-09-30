package sqlconverter

type SqlDialect string

const (
	SQLSERVER SqlDialect = "sqlserver"
	MYSQL     SqlDialect = "mysql"
	POSTGRES  SqlDialect = "postgres"
)
