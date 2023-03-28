package binder

import (
	"reflect"
)

// IBinder 绑定器接口
type IBinder interface {
	Bind(value any, v any) error
}

// GetBinder 获取绑定器
func GetBinder(v reflect.Kind) IBinder {
	switch v {
	case reflect.String:
		return NewStringBinder()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return NewNumberBinder()
	case reflect.Bool:
		return NewBoolBinder()
	}
	panic("no parser for: " + v.String())
}
