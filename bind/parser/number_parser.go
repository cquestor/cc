package parser

import (
	"fmt"
	"reflect"
	"strconv"
)

// NumberParser 数字解析器
type NumberParser struct{}

// NewNumberParser 构造数字解析器
func NewNumberParser() *NumberParser {
	return &NumberParser{}
}

// GetData 实现 IParser 接口
func (parser *NumberParser) GetData(v any, target reflect.Kind) (any, error) {
	switch target {
	case reflect.Bool:
		return parser.bool(v), nil
	case reflect.String:
		return parser.string(v), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return v, nil
	default:
		return nil, fmt.Errorf("invalid target: %s", target)
	}
}

// bool 将数字转换成布尔值
func (parser *NumberParser) bool(v any) bool {
	return v != 0
}

// string 将数字转换成字符串
func (parser *NumberParser) string(v any) string {
	switch v := v.(type) {
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	}
	return ""
}
