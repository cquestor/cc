package cc

import (
	"encoding/json"
	"fmt"
)

// typeResponse 响应类型
type typeResponse int

const (
	respString typeResponse = iota
	respHtml
	respJson
	respBytes
)

// Response 响应
type Response struct {
	Code int
	Data any
	Type typeResponse
}

// Invoke 响应执行
func (response *Response) Invoke(ctx *Context) {
	switch response.Type {
	case respString:
		ctx.SetHeader("Content-Type", "text/plain; charset=utf-8")
		ctx.Writer.WriteHeader(response.Code)
		ctx.Writer.Write([]byte(response.Data.(string)))
	case respHtml:
		ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
		ctx.Writer.WriteHeader(response.Code)
		ctx.Writer.Write(response.Data.([]byte))
	case respJson:
		ctx.SetHeader("Content-Type", "application/json; charset=utf-8")
		ctx.Writer.WriteHeader(response.Code)
		encoder := json.NewEncoder(ctx.Writer)
		if err := encoder.Encode(response.Data); err != nil {
			panic(err)
		}
	case respBytes:
		ctx.Writer.WriteHeader(response.Code)
		ctx.Writer.Write(response.Data.([]byte))
	}
}

// String 构造字符串响应
func String(code int, format string, v ...any) *Response {
	return &Response{
		Code: code,
		Data: fmt.Sprintf(format, v...),
		Type: respString,
	}
}

// Html 构造网页响应
func Html(code int, v []byte) *Response {
	return &Response{
		Code: code,
		Data: v,
		Type: respHtml,
	}
}

// Json 构造 json 响应
func Json(code int, v any) *Response {
	return &Response{
		Code: code,
		Data: v,
		Type: respJson,
	}
}

// Data 构造二进制响应
func Data(code int, v []byte) *Response {
	return &Response{
		Code: code,
		Data: v,
		Type: respBytes,
	}
}
