package cc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response 响应接口
type Response interface {
	Invoke(ctx *Context)
}

// responseString 字符串响应
type responseString struct {
	Code  int
	Value string
}

// responseHtml 网页响应
type responseHtml struct {
	Code  int
	Value []byte
}

// responseJson json响应
type responseJson struct {
	Code  int
	Value any
}

// responseData 字节流制响应
type responseData struct {
	Code  int
	Value []byte
}

// responseRedirect 重定向响应
type responseRedirect struct {
	Code  int
	Value string
}

// responseCode 状态码响应
type responseCode struct {
	Code int
}

// String 构造字符串响应
func String(code int, format string, v ...any) *responseString {
	return &responseString{
		Code:  code,
		Value: fmt.Sprintf(format, v...),
	}
}

func (response *responseString) Invoke(ctx *Context) {
	ctx.SetHeader("Content-Type", "text/plain; charset=utf-8")
	ctx.setStatusCode(response.Code)
	ctx.Writer.Write([]byte(response.Value))
}

// Html 构造网页响应
func Html(code int, v []byte) *responseHtml {
	return &responseHtml{
		Code:  code,
		Value: v,
	}
}

func (response *responseHtml) Invoke(ctx *Context) {
	ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
	ctx.setStatusCode(response.Code)
	ctx.Writer.Write(response.Value)
}

// Json 构造 json 响应
func Json(code int, v any) *responseJson {
	return &responseJson{
		Code:  code,
		Value: v,
	}
}

func (response *responseJson) Invoke(ctx *Context) {
	ctx.SetHeader("Content-Type", "application/json; charset=utf-8")
	ctx.setStatusCode(response.Code)
	encoder := json.NewEncoder(ctx.Writer)
	if err := encoder.Encode(response.Value); err != nil {
		panic(err)
	}
}

// Data 构造字节流响应
func Data(code int, v []byte) *responseData {
	return &responseData{
		Code:  code,
		Value: v,
	}
}

func (response *responseData) Invoke(ctx *Context) {
	ctx.setStatusCode(response.Code)
	ctx.Writer.Write(response.Value)
}

// Redirect 构造重定向响应
func Redirect(code int, v string) *responseRedirect {
	return &responseRedirect{
		Code:  code,
		Value: v,
	}
}

func (response *responseRedirect) Invoke(ctx *Context) {
	http.Redirect(ctx.Writer, ctx.Req, response.Value, response.Code)
}

// Code 构造状态码响应
func Code(code int) *responseCode {
	return &responseCode{
		Code: code,
	}
}

func (response *responseCode) Invoke(ctx *Context) {
	ctx.setStatusCode(response.Code)
}
