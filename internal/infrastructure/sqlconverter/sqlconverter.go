package sqlconverter

import "reflect"

type (
	SqlConverter struct {
		Dialect SqlDialect
		Options ConvertOptions
	}
	ColumnDefinition struct {
		Name          string
		Type          string
		Nullable      bool
		PrimaryKey    bool
		AutoIncrement bool
		Unique        bool
		Default       string
		Size          int
		Precision     int
		Scale         int
		Comment       string
		References    string
	}

	ConvertOptions struct {
		TablePrefix     string
		UseBackticks    bool
		IncludeIndexes  bool
		IncludeComments bool
		CamelToSnake    bool
	}

	IndexDefinition struct {
		Name    string
		Type    string // INDEX, UNIQUE, PRIMARY
		Columns []string
	}

	// TableDefinition represents a complete table definition
	TableDefinition struct {
		Name    string
		Columns []ColumnDefinition
		Indexes []IndexDefinition
		Comment string
	}
)

func NewSqlConverter(dialect SqlDialect) *SqlConverter {
	return &SqlConverter{
		Dialect: dialect,
	}
}

func (s *SqlConverter) parseStruct(valueType reflect.Type) (*TableDefinition, error) {

	tableDefinition := &TableDefinition{}

	return tableDefinition, nil

}
