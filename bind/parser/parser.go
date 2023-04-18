package parser

import "reflect"

// IParser 解析器接口
type IParser interface {
	GetData(v any, target reflect.Kind) (any, error)
}

// GetParser 获取指定解析器
func GetParser(name string) IParser {
	switch name {
	case "number":
		return NewNumberParser()
	case "string":
		return NewStringParser()
	case "bool":
		return NewBoolParser()
	}
	panic("no parser for: " + name)
}
