package repositories

import (
	"database/sql"
	"fmt"
	"oriongo/internal/infrastructure"
	"os"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	BaseRepository struct {
		_context       *infrastructure.DbContext
		_tableName     string
		_selectBuilder sq.SelectBuilder
		_config        infrastructure.ConnectionConfig
	}
	DB struct {
	}

	QBuilder struct{}
)

func dsn(config infrastructure.ConnectionConfig) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database)
	return connectionString
}

func OpenDB() *gorm.DB {
	var config infrastructure.ConnectionConfig = infrastructure.ConnectionConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
	}
	dsn := dsn(config)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func openSql(config infrastructure.ConnectionConfig) *sql.DB {
	db, err := sql.Open("mysql", dsn(config))
	if err != nil {
		panic(err)
	}
	return db
}

func NewBaseRepository(ctx infrastructure.DbContext) *BaseRepository {
	/*
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Database)

		db, error := sql.Open("mysql", connectionString)
		if error != nil {
			panic(error)
		}

		if pingError := db.Ping(); pingError != nil {
			panic(pingError)
		}

	*/
	//db := OpenDB()
	return &BaseRepository{
		_context: &ctx,
	}
}

func FirstOrDefault[T any](rows *sql.Rows) interface{} {
	var zeroOfTargetType T
	if response, error := QueryToStruct[T](rows); error != nil {
		return zeroOfTargetType
	} else {
		if len(response) > 0 {
			return response[0]
		} else {
			return zeroOfTargetType
		}
	}

	return zeroOfTargetType
}

func (b BaseRepository) selectColumnsAndValuesFromInterface(value interface{}) ([]string, []interface{}) {
	target := reflect.ValueOf(value)
	targetType := target.Type()

	var columns = make([]string, targetType.NumField())
	var values = make([]interface{}, len(columns))
	for i := 0; i < targetType.NumField(); i++ {
		columns[i] = targetType.Field(i).Name

		field := target.Field(i)
		if field.Kind() == reflect.Ptr {
			values[i] = target.Field(i).Elem().Interface()
		} else {
			values[i] = target.Field(i).Interface()
		}
	}

	return columns, values
}

func (b BaseRepository) tableName(value interface{}) string {
	if value == nil {
		return ""
	}

	if value, ok := value.(string); ok {
		return value
	}
	target := reflect.ValueOf(value)
	targetType := target.Type()

	return targetType.Name()
}

func (b *BaseRepository) Insert(value interface{}) (int, error) {
	target := reflect.ValueOf(value)
	targetType := target.Type()

	var columns = make([]string, targetType.NumField())
	var values = make([]interface{}, len(columns))
	for i := 0; i < targetType.NumField(); i++ {
		columns[i] = targetType.Field(i).Name
		field := target.Field(i)
		if field.Kind() == reflect.Ptr {
			values[i] = target.Field(i).Elem().Interface()
		} else {
			values[i] = target.Field(i).Interface()
		}
	}

	if b._tableName == "" {
		b._tableName = targetType.Name()
	}

	result, insertError := sq.Insert(b._tableName).Columns(columns...).Values(values).Exec()
	if insertError != nil {
		return -223, insertError
	}

	rowsAffected, rowsAffectedError := result.RowsAffected()
	if rowsAffectedError != nil {
		return -223, rowsAffectedError
	}

	return int(rowsAffected), nil
}

func (b *BaseRepository) Find(tableName string, query string, values interface{}, columns ...string) (*sql.Rows, error) {
	//columns, values := b.selectColumnsAndValuesFromInterface(args)
	//tableName := b.tableName(args)

	var selectBuilder sq.SelectBuilder
	if len(columns) == 0 {
		selectBuilder = sq.Select("*").From(tableName).Where(query, values)
	} else {
		selectBuilder = sq.Select(columns...).From(tableName).Where(query, values)
	}

	rows, rowError := selectBuilder.Query()
	if rowError != nil {
		return nil, rowError
	}

	return rows, nil
}

func (b *BaseRepository) Table(tableName string) *BaseRepository {
	b._tableName = tableName
	return b
}

func (b *BaseRepository) FindOne(resultType, query, values interface{}, columns ...string) (interface{}, error) {
	var selectBuilder sq.SelectBuilder
	if len(columns) == 0 {
		selectBuilder = sq.Select("*").From(b._tableName).Where(query, values)
	} else {
		selectBuilder = sq.Select(columns...).From(b._tableName).Where(query, values)
	}

	context := openSql(b._config)
	selectBuilder = selectBuilder.Limit(1)
	sql, _, _ := selectBuilder.ToSql()
	rows, rowError := context.Query(sql)
	if rowError != nil {
		return nil, rowError
	}

	err := queryRow(rows, resultType, true)
	if err != nil {
		return nil, err
	}

	return resultType, nil
}

func (b *BaseRepository) FirstOrDefault(rows *sql.Rows, destination interface{}) error {
	if error := queryRow(rows, destination, true); error != nil {
		return error
	}

	return nil
}

func (b *BaseRepository) DB() *gorm.DB {
	//return b._context
	//return OpenDB()
	return (b._context.DB())
}

func findTag(value reflect.Value, tagName string) reflect.Value {
	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag.Get(tagName)
		if tag != "-" || tag != "" {
			continue
		}

		return value.Field(i)
	}

	return reflect.Value{}
}

func QueryToStruct[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	var results = make([]T, 0)
	if _, error := ScanToStruct(rows, &results); error != nil {
		return nil, error
	}

	return results, nil
}

func queryRow(rows *sql.Rows, target interface{}, firstItem bool) error {
	defer rows.Close()
	destination := reflect.ValueOf(target).Elem()
	destinationType := destination.Type()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	var values = make([]interface{}, len(columns))
	for i, _ := range columns {
		values[i] = &values[i]
	}

	for rows.Next() {
		newElement := reflect.New(destinationType).Elem()
		err := rows.Scan(values...)
		if err != nil {
			return err
		}

		for i, column := range columns {
			field := newElement.FieldByName(column)
			if field.IsValid() && field.CanSet() {
				setValue(field, values[i])
			}
		}

		destination.Set(reflect.Append(destination, newElement.Addr()))
		if firstItem {
			break
		}
	}

	return nil
}

func ScanToStruct[T any](rows *sql.Rows, dest T) ([]T, error) {
	target := reflect.ValueOf(dest).Elem()
	targetType := target.Type()

	columns, err := rows.Columns()
	if err != nil {
		return make([]T, 0), err
	}

	var values = make([]interface{}, len(columns))
	for i, _ := range values {
		values[i] = &columns[i]
	}

	for rows.Next() {
		newElement := reflect.New(targetType).Elem()
		if err := rows.Scan(values...); err != nil {
			return make([]T, 0), err
		}

		for i, _ := range columns {
			fieldTag := findTag(newElement, "db")
			if fieldTag.IsValid() && fieldTag.CanSet() {
				setValue(fieldTag, values[i])
			}
		}

		target.Set(reflect.Append(target, newElement.Addr()))
	}

	return make([]T, 0), rows.Err()
}

func setValue(field reflect.Value, value interface{}) {
	if value == nil {
		return
	}

	fieldType := field.Type()
	valueType := reflect.TypeOf(value)

	if fieldType.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(fieldType.Elem()))
		}

		field = field.Elem()
		fieldType = fieldType.Elem()
	}

	if valueType.ConvertibleTo(fieldType) {
		field.Set(reflect.ValueOf(value).Convert(fieldType))
	}
}

func (r *BaseRepository) First(rows *sql.Rows, destination interface{}) interface{} {
	var response = make([]interface{}, 0)
	destinationValue := reflect.ValueOf(destination).Elem()
	for rows.Next() {
		var current = reflect.New(reflect.TypeOf(destinationValue))
		var currentValue = reflect.ValueOf(current.Interface())
		args := make([]interface{}, 0)

		for i := 0; i < currentValue.NumField(); i++ {
			args = append(args, currentValue.Field(i).Addr().Interface())
		}

		rows.Scan(args...)

		response = append(response, current)
		break
	}

	if len(response) == 0 {
		return nil
	}
	return response[0]
}

func (r *BaseRepository) build() string {
	query, _, error := r._selectBuilder.ToSql()
	if error != nil {
		panic(error)
	}

	return query
}
