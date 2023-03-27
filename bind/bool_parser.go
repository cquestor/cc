package bind

// BoolParser 布尔解析器
type BoolParser struct{}

// NewBoolParser 构造布尔解析器
func NewBoolParser() *BoolParser {
	return &BoolParser{}
}

// string 将布尔值转换成字符串
func (parser *BoolParser) string(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
