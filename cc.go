package cc

// J json数据格式
type J map[string]any

// IResponse 响应接口
type IResponse interface {
	Invoke(*Context)
}

// IHandler 处理器接口
type IHandler interface {
	Invoke(*Context) IResponse
}

// HandlerProvider 处理器存储结构
type HandlerProvider map[string]map[string]IHandler

// AddHandler 添加处理器
func (provider *HandlerProvider) AddHandler(method, pattern string, handler IHandler) {
	if root, ok := (*provider)[method]; ok {
		root[pattern] = handler
	} else {
		temp := map[string]IHandler{pattern: handler}
		(*provider)[method] = temp
	}
}
