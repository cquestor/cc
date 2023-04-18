package orm

import (
	"database/sql"
	"reflect"
	"strings"
)

const myTag = "data"

// parseInsertObject 解析插入对象
func parseInsertObject(v any) ([]string, []any, reflect.Type) {
	var fields []string
	var values []any
	t := reflect.TypeOf(v)
	e := reflect.ValueOf(v)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		if interested, tag := isInsertInterested(field.Tag.Get(myTag)); !interested {
			continue
		} else if tag != "" {
			fields = append(fields, tag)
			values = append(values, e.Field(i).Interface())
		} else {
			fields = append(fields, strings.ToLower(field.Name))
			values = append(values, e.Field(i).Interface())
		}
	}
	return fields, values, t
}

// parseSelectObject 解析查询对象
func parseSelectObject(t reflect.Type) []string {
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		if interested, tag := isSelectInterested(field.Tag.Get(myTag)); !interested {
			continue
		} else if tag != "" {
			fields = append(fields, tag)
		} else {
			fields = append(fields, strings.ToLower(field.Name))
		}
	}
	return fields
}

// isInsertInterested 是否是要处理的插入字段
func isInsertInterested(tag string) (bool, string) {
	parts := strings.Split(tag, ",")
	if len(parts) < 2 {
		return true, parts[0]
	}
	mark := parts[1]
	switch mark {
	case "auto_increment", "use_default":
		return false, ""
	default:
		return true, parts[0]
	}
}

// isSelectInterested 是否是要处理的查询字段
func isSelectInterested(tag string) (bool, string) {
	parts := strings.Split(tag, ",")
	if len(parts) < 2 {
		return true, parts[0]
	}
	mark := parts[1]
	switch mark {
	case "-":
		return false, ""
	default:
		return true, parts[0]
	}
}

// scanFetchObject 多结构体赋值
func scanFetchObject(t reflect.Type, rows *sql.Rows) (reflect.Value, error) {
	newElem := reflect.New(t)
	addrs := make([]any, 0)
	for i := 0; i < newElem.Elem().NumField(); i++ {
		field := newElem.Elem().Field(i)
		if !field.CanInterface() {
			continue
		}
		if interested, _ := isSelectInterested(t.Field(i).Tag.Get(myTag)); !interested {
			continue
		}
		addrs = append(addrs, field.Addr().Interface())
	}
	if err := rows.Scan(addrs...); err != nil {
		return reflect.Value{}, err
	}
	return newElem, nil
}

// scanOneObject 单结构体赋值
func scanOneObject(v any, rows *sql.Rows) error {
	_value := reflect.ValueOf(v)
	addrs := make([]any, 0)
	for i := 0; i < _value.Elem().NumField(); i++ {
		field := _value.Elem().Field(i)
		if !field.CanInterface() {
			continue
		}
		if interested, _ := isSelectInterested(_value.Elem().Type().Field(i).Tag.Get(myTag)); !interested {
			continue
		}
		addrs = append(addrs, field.Addr().Interface())
	}
	if err := rows.Scan(addrs...); err != nil {
		return err
	}
	return nil
}
