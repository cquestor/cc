package parser

import (
	"fmt"
	"reflect"
)

// BoolParser 布尔解析器
type BoolParser struct{}

// NewBoolParser 构造布尔解析器
func NewBoolParser() *BoolParser {
	return &BoolParser{}
}

// GetData 实现 IParser 接口
func (parser *BoolParser) GetData(v any, target reflect.Kind) (any, error) {
	switch target {
	case reflect.String:
		return parser.string(v.(bool)), nil
	case reflect.Bool:
		return v, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return parser.number(v.(bool), target)
	default:
		return nil, fmt.Errorf("invalid target: %s", target)
	}
}

// string 将布尔值转换成字符串
func (parser *BoolParser) string(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// number 将布尔值转换成数字
func (parser *BoolParser) number(v bool, target reflect.Kind) (any, error) {
	switch target {
	case reflect.Uint:
		if v {
			return uint(1), nil
		}
		return uint(0), nil
	case reflect.Uint8:
		if v {
			return uint8(1), nil
		}
		return uint8(0), nil
	case reflect.Uint16:
		if v {
			return uint16(1), nil
		}
		return uint16(0), nil
	case reflect.Uint32:
		if v {
			return uint32(1), nil
		}
		return uint32(0), nil
	case reflect.Uint64:
		if v {
			return uint64(1), nil
		}
		return uint64(0), nil
	case reflect.Int:
		if v {
			return 1, nil
		}
		return 0, nil
	case reflect.Int8:
		if v {
			return int8(1), nil
		}
		return int8(0), nil
	case reflect.Int16:
		if v {
			return int16(1), nil
		}
		return int16(0), nil
	case reflect.Int32:
		if v {
			return int32(1), nil
		}
		return int32(0), nil
	case reflect.Int64:
		if v {
			return int64(1), nil
		}
		return int64(0), nil
	case reflect.Float32:
		if v {
			return float32(1), nil
		}
		return float32(0), nil
	case reflect.Float64:
		if v {
			return float64(1), nil
		}
		return float64(0), nil
	}
	return nil, fmt.Errorf("invalid target: %s", target)
}
