package cc

// Handler 处理器
type Handler func(*Context) IResponse

// Invoke 实现 IHandler 接口
func (handler Handler) Invoke(ctx *Context) IResponse {
	return handler(ctx)
}
