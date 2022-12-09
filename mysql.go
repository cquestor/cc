package cc

import (
	"fmt"
	"reflect"
	"time"
)

type mysql struct{}

var _ dialect = (*mysql)(nil)

func init() {
	registerDialect("mysql", &mysql{})
}

func (m *mysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Int, reflect.Int32:
		return "integer"
	case reflect.Int8:
		return "tinyint"
	case reflect.Int16:
		return "smallint"
	case reflect.Int64:
		return "bigint"
	case reflect.Bool:
		return "bool"
	case reflect.String:
		return "varchar(255)"
	case reflect.Float32, reflect.Float64:
		return "double precision"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (m *mysql) TableExistSQL(tableName string) (string, []any) {
	args := []any{tableName}
	return "SHOW TABLES LIKE ?", args
}
