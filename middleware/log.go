package middleware

import (
	"github.com/cquestor/cc"
	"github.com/cquestor/cc/logger"
)

// CLogMiddleware 日志中间件
type CLogMiddleware struct{}

func (log *CLogMiddleware) Instance() func(*cc.Context) cc.Response {
	return func(ctx *cc.Context) cc.Response {
		ipAddress := ctx.Req.RemoteAddr
		if forwarded := ctx.Header("X-Forwarded-For"); forwarded != "" {
			ipAddress = forwarded
		}
		cc.LogInfof("%s\033[7;32m from \033[1m%s\033[7;35m ==> %s\n", logger.Style(logger.ColorCyan, logger.StyleInverse, " ", ctx.Method, " "), logger.Style(logger.ColorBlue, logger.StyleInverse, " ", ipAddress, " "), logger.Style(logger.ColorRed, logger.StyleBold, " ", ctx.Path, " "))
		return nil
	}
}
