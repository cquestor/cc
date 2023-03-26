package orm

import (
	"reflect"
	"strings"
)

// parseObject 解析对象
func parseObject(v any) ([]string, []any, reflect.Type) {
	var fields []string
	var values []any
	t := reflect.TypeOf(v)
	e := reflect.ValueOf(v)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		if interested, tag := isInterested(field.Tag.Get("cc")); !interested {
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

// isInterested 是否是要处理的字段
func isInterested(tag string) (bool, string) {
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
