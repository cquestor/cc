package binder

import (
	"errors"
	"reflect"

	"github.com/cquestor/cc/internal/bind/parser"
)

// StringBinder 字符串绑定器
type StringBinder struct {
	parser parser.IParser
}

// NewStringBinder 构造字符串绑定器
func NewStringBinder() *StringBinder {
	return &StringBinder{
		parser: parser.GetParser("string"),
	}
}

// Bind 绑定字符串
func (binder *StringBinder) Bind(value any, v any) error {
	_v := reflect.ValueOf(v)
	if _v.Kind() != reflect.Ptr {
		return errors.New("bind target must be a ptr")
	}
	_t := _v.Elem().Type()
	dst, err := binder.parser.GetData(value, _t.Kind())
	if err != nil {
		return err
	}
	_v.Elem().Set(reflect.ValueOf(dst))
	return nil
}
