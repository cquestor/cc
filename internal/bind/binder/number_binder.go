package binder

import (
	"errors"
	"reflect"

	"github.com/cquestor/cc/internal/bind/parser"
)

// NumberBinder 数字绑定器
type NumberBinder struct {
	parser parser.IParser
}

// NewNumberBinder 构造数字绑定器
func NewNumberBinder() *NumberBinder {
	return &NumberBinder{
		parser: parser.GetParser("number"),
	}
}

// Bind 绑定数字
func (binder *NumberBinder) Bind(value any, v any) error {
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
