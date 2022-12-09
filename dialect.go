package cc

import "reflect"

var dialectsMap = map[string]dialect{}

type dialect interface {
	DataTypeOf(typ reflect.Value) string
	TableExistSQL(tableName string) (string, []any)
}

func registerDialect(name string, dialect dialect) {
	dialectsMap[name] = dialect
}

func getDialect(name string) (dialect dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
