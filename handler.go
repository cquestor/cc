package cc

// IHandler 处理器接口
type IHandler interface {
	Invoke(ctx *Context) Response
}

// Handler 处理器
type Handler func(ctx *Context) Response

// Invoke 实现 IHandler 接口
func (handler Handler) Invoke(ctx *Context) Response {
	return handler(ctx)
}
