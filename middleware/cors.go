package middleware

import (
	"net/http"
	"strings"

	"github.com/cquestor/cc"
)

// CorsMiddleware 跨域中间件
type CorsMiddleware struct {
	Origin  string
	Methods string
	Headers string
}

// SetOrigin设置跨域源
func (cors *CorsMiddleware) SetOrigin(origin string) {
	cors.Origin = origin
}

// SetMethods 设置跨域方法
func (cors *CorsMiddleware) SetMethods(methods ...string) {
	cors.Methods = strings.Join(methods, ", ")
}

// SetHeaders 设置跨域请求头
func (cors *CorsMiddleware) SetHeaders(headers ...string) {
	cors.Headers = strings.Join(headers, ", ")
}

// Cors 跨域设置
func (cors *CorsMiddleware) Instance() func(*cc.Context) cc.Response {
	if cors.Origin == "" {
		cors.Origin = "*"
	}
	if cors.Methods == "" {
		cors.Methods = "*"
	}
	if cors.Headers == "" {
		cors.Methods = "*"
	}
	return func(ctx *cc.Context) cc.Response {
		ctx.SetHeader("Access-Control-Allow-Origin", cors.Origin)
		ctx.SetHeader("Access-Control-Allow-Methods", cors.Methods)
		ctx.SetHeader("Access-Control-Allow-Headers", cors.Headers)
		if ctx.Method == http.MethodOptions {
			return cc.Code(http.StatusNoContent)
		}
		return nil
	}
}
