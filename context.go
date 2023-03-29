package cc

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

// Context 上下文
type Context struct {
	Req    *http.Request
	Writer http.ResponseWriter
	Method string
	Path   string
	Params map[string]string
}

// NewContext 新建上下文
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Req:    r,
		Writer: w,
		Method: r.Method,
		Path:   r.URL.Path,
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

// Body 获取请求体
func (ctx *Context) Body() []byte {
	b, _ := io.ReadAll(ctx.Req.Body)
	ctx.Req.Body = io.NopCloser(bytes.NewReader(b))
	return b
}

// Header 获取请求头
func (ctx *Context) Header(key string) string {
	return ctx.Req.Header.Get(key)
}

// SetHeader 设置响应头
func (ctx *Context) SetHeader(key, value string) {
	ctx.Writer.Header().Set(key, value)
}

// Cookie 获取请求Cookie
func (ctx *Context) Cookie(key string) *http.Cookie {
	if cookie, err := ctx.Req.Cookie(key); err != nil {
		return nil
	} else {
		return cookie
	}
}

// SetCookie 设置响应Cookie
func (ctx *Context) SetCookie(c *http.Cookie) {
	http.SetCookie(ctx.Writer, c)
}

// File 获取上传文件
// TODO 多文件上传，文件关闭
func (ctx *Context) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return ctx.Req.FormFile(key)
}

// SetMaxFileSize 设置允许最大内存及文件大小，单位字节
func (ctx *Context) SetMaxFileSize(v int64) error {
	return ctx.Req.ParseMultipartForm(v)
}

// setStatusCode 设置响应状态码
func (ctx *Context) setStatusCode(code int) {
	ctx.Writer.WriteHeader(code)
}
