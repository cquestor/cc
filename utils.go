package cc

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// handleMiddlewares 处理中间件
func handleMiddlewares(ctx *Context, middlewares []IHandler) Response {
	for _, handler := range middlewares {
		if response := handler.Invoke(ctx); response != nil {
			return response
		}
	}
	return nil
}

// handleErr 处理错误
func handleErr(ctx *Context) {
	if err := recover(); err != nil {
		message := trace(fmt.Sprintf("%s", err))
		LogErrf("%s\n\n", message)
		Code(http.StatusInternalServerError).Invoke(ctx)
	}
}

// trace 堆栈信息
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func banner() {
	fmt.Println(" \033[1;32m   ______   \033[1;36m____     \033[1;33m_   __    \033[1;31m______\033[0m")
	fmt.Println(" \033[1;32m  / ____/  \033[1;36m/ __ \\   \033[1;33m/ | / /   \033[1;31m/ ____/\033[0m")
	fmt.Println(" \033[1;32m / / __   \033[1;36m/ / / /  \033[1;33m/  |/ /   \033[1;31m/ __/   \033[0m")
	fmt.Println(" \033[1;32m/ /_/ /  \033[1;36m/ /_/ /  \033[1;33m/ /|  /   \033[1;31m/ /___   \033[0m")
	fmt.Println(" \033[1;32m\\____/   \033[1;36m\\____/  \033[1;33m/_/ |_/   \033[1;31m/_____/   \033[0m")
	fmt.Println()
}
