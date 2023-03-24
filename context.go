package cc

import (
	"bytes"
	"io"
	"net/http"
)

// Context 上下文
type Context struct {
	Req        *http.Request
	Writer     http.ResponseWriter
	Method     string
	Path       string
	Params     map[string]string
	StatusCode int
}

// NewContext 构造Context
func NewContext(w http.ResponseWriter, r *http.Request, params map[string]string) *Context {
	return &Context{
		Req:    r,
		Writer: w,
		Method: r.Method,
		Path:   r.URL.Path,
		Params: params,
	}
}

// Param 获取路由参数
func (ctx *Context) Param(key string) string {
	return ctx.Params[key]
}

// Query 获取请求参数
func (ctx *Context) Query(key string) string {
	return ctx.Req.URL.Query().Get(key)
}

// PostForm 获取表单参数
func (ctx *Context) PostForm(key string) string {
	return ctx.Req.PostFormValue(key)
}

// Header 获取请求头
func (ctx *Context) Header(key string) string {
	return ctx.Req.Header.Get(key)
}

// Body 获取请求体
func (ctx *Context) Body() []byte {
	b, _ := io.ReadAll(ctx.Req.Body)
	ctx.Req.Body = io.NopCloser(bytes.NewReader(b))
	return b
}

// SetHeader 设置响应头
func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}
