package binder

import (
	"errors"
	"reflect"

	"github.com/cquestor/cc/bind/parser"
)

// BoolBinder 布尔绑定器
type BoolBinder struct {
	parser parser.IParser
}

// NewBoolBinder 构造布尔绑定器
func NewBoolBinder() *BoolBinder {
	return &BoolBinder{
		parser: parser.GetParser("bool"),
	}
}

// Bind 绑定布尔值
func (binder *BoolBinder) Bind(value any, v any) error {
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
