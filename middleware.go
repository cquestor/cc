package cc

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type CorsConfig struct {
	AllowOrigin      string
	AllowMethods     string
	AllowHeaders     string
	ExposeHeaders    string
	MaxAge           int
	AllowCredentials bool
}

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		ip, err := c.GetIP()
		if err != nil {
			ip = "unknow address"
		}
		c.Next()
		lost := time.Since(t)
		var timeColor int
		if lost <= time.Second {
			timeColor = 32
		} else {
			timeColor = 33
		}
		var statusColor int
		if 200 <= c.StatusCode && c.StatusCode < 300 {
			statusColor = 32
		} else if 300 <= c.StatusCode && c.StatusCode < 400 {
			statusColor = 33
		} else if 400 <= c.StatusCode && c.StatusCode < 500 {
			statusColor = 35
		} else if 500 <= c.StatusCode && c.StatusCode < 600 {
			statusColor = 31
		} else {
			statusColor = 37
		}
		Infof("\033[2m[\033[0m\033[1;%dm%d\033[0m\033[2m]\033[0m \033[1m%s\033[0m \033[2m-->\033[0m \033[1m%s\033[0m \033[2mfrom\033[0m \033[1m%s\033[0m \033[2min\033[0m \033[1;%dm%v\033[0m", statusColor, c.StatusCode, c.Method, c.Path, ip, timeColor, lost)
	}
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := trace(fmt.Sprintf("%s", err))
				Errorf("%s\n\n", message)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func Cors(config *CorsConfig) HandlerFunc {
	if config == nil {
		config = new(CorsConfig)
		config.AllowOrigin = "*"
		config.AllowHeaders = "*"
		config.AllowMethods = "*"
		config.MaxAge = 60 * 60 * 24
		config.AllowCredentials = true
	}
	return func(c *Context) {
		if config.AllowOrigin != "" {
			c.SetHeader("Access-Control-Allow-Origin", config.AllowOrigin)
		}
		if config.AllowMethods != "" {
			c.SetHeader("Access-Control-Allow-Methods", config.AllowMethods)
		}
		if config.AllowHeaders != "" {
			c.SetHeader("Access-Control-Allow-Headers", config.AllowHeaders)
		}
		if config.ExposeHeaders != "" {
			c.SetHeader("Access-Control-Expose-Headers", config.ExposeHeaders)
		}
		if config.MaxAge > 0 {
			c.SetHeader("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
		}
		if config.AllowCredentials {
			c.SetHeader("Access-Control-Allow-Credentials", "true")
		}
		if c.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
		}
	}
}

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
