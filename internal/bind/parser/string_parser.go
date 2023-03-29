package parser

import (
	"fmt"
	"reflect"
	"strconv"
)

// StringParser 字符串解析器
type StringParser struct{}

// NewStringParser 构造字符串解析器
func NewStringParser() *StringParser {
	return &StringParser{}
}

// GetData 实现 IParser 接口
func (parser *StringParser) GetData(v any, target reflect.Kind) (any, error) {
	switch target {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		return parser.number(v.(string), target)
	case reflect.Bool:
		return parser.bool(v.(string))
	case reflect.String:
		return v, nil
	default:
		return nil, fmt.Errorf("invalid target: %s", target)
	}
}

// number 将字符串转换成数字
func (parser *StringParser) number(v string, target reflect.Kind) (any, error) {
	switch target {
	case reflect.Uint:
		if i, err := strconv.ParseUint(v, 10, 10); err != nil {
			return nil, err
		} else {
			return uint(i), nil
		}
	case reflect.Uint8:
		if i, err := strconv.ParseUint(v, 10, 8); err != nil {
			return nil, err
		} else {
			return uint8(i), nil
		}
	case reflect.Uint16:
		if i, err := strconv.ParseUint(v, 10, 16); err != nil {
			return nil, err
		} else {
			return uint16(i), nil
		}
	case reflect.Uint32:
		if i, err := strconv.ParseUint(v, 10, 32); err != nil {
			return nil, err
		} else {
			return uint32(i), nil
		}
	case reflect.Uint64:
		if i, err := strconv.ParseUint(v, 10, 64); err != nil {
			return nil, err
		} else {
			return i, nil
		}
	case reflect.Int:
		if i, err := strconv.Atoi(v); err != nil {
			return nil, err
		} else {
			return i, nil
		}
	case reflect.Int8:
		if i, err := strconv.ParseInt(v, 10, 8); err != nil {
			return nil, err
		} else {
			return int8(i), nil
		}
	case reflect.Int16:
		if i, err := strconv.ParseInt(v, 10, 16); err != nil {
			return nil, err
		} else {
			return int16(i), nil
		}
	case reflect.Int32:
		if i, err := strconv.ParseInt(v, 10, 32); err != nil {
			return nil, err
		} else {
			return int32(i), nil
		}
	case reflect.Int64:
		if i, err := strconv.ParseInt(v, 10, 64); err != nil {
			return nil, err
		} else {
			return i, nil
		}
	case reflect.Float32:
		if i, err := strconv.ParseFloat(v, 32); err != nil {
			return nil, err
		} else {
			return float32(i), nil
		}
	case reflect.Float64:
		if i, err := strconv.ParseFloat(v, 64); err != nil {
			return nil, err
		} else {
			return i, nil
		}
	}
	return nil, fmt.Errorf("invalid target type, not number: %s", target)
}

// bool 将字符串转换成布尔值
func (parser *StringParser) bool(v string) (bool, error) {
	if v != "0" || v == "true" {
		return true, nil
	}
	return false, nil
}
