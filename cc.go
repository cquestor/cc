// Copyright 2023 cquestor. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cc

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/cquestor/cc/internal/orm"
	"github.com/cquestor/cc/internal/router"
)

// J json结构
type J map[string]any

// Engine Web引擎
type Engine struct {
	*RouteGroup
	config   *AppConfig
	router   router.IRouter
	handlers map[string]map[string]IHandler
	options  map[string]any
	database *orm.Engine
	groups   []*RouteGroup
}

// RouteGroup 分组路由
type RouteGroup struct {
	prefix  string
	befores []IHandler
	afters  []IHandler
	parent  *RouteGroup
	engine  *Engine
}

type (
	CAppConfig   []byte // 项目配置
	CTLSCertFile string // TLS证书
	CTLSKeyFile  string // TLS密钥
)

const (
	optAppConfig = "AppConfig"
	optCertPath  = "CertPath"
	optKeyPath   = "KeyPath"
)

// New 构造Engine
func New() *Engine {
	engine := &Engine{
		config:   NewAppConfig(),
		router:   router.NewRouter(),
		handlers: make(map[string]map[string]IHandler),
		options:  make(map[string]any),
	}
	engine.RouteGroup = &RouteGroup{engine: engine}
	engine.groups = []*RouteGroup{engine.RouteGroup}
	return engine
}

// Group 创建新的路由分组
func (group *RouteGroup) Group(prefix string) *RouteGroup {
	newGroup := &RouteGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: group.engine,
	}
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

// Run 启动 Web Server
func (engine *Engine) Run(options ...any) {
	engine.parseOptions(options...)
	if err := engine.parseConfig(); err != nil {
		LogErrf("Parse config err: %v\n", err)
		os.Exit(1)
	}
	if err := engine.initConfig(); err != nil {
		LogErrf("Init config err: %v\n", err)
		os.Exit(1)
	}
	listenAddr := fmt.Sprintf(":%d", engine.config.Port)
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      engine,
		ReadTimeout:  time.Duration(engine.config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(engine.config.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(engine.config.IdleTimeout) * time.Second,
	}
	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go engine.shutdown(server, quit, done)
	LogInfo("Server is ready to handle requests at", listenAddr)
	engine.start(server)
	<-done
	LogInfo("Server stopped")
}

// Get 添加 GET 请求
func (group *RouteGroup) Get(pattern string, handler func(*Context) Response) {
	group.addRoute(http.MethodGet, pattern, Handler(handler))
}

// Post 添加 POST 请求
func (group *RouteGroup) Post(pattern string, handler func(*Context) Response) {
	group.addRoute(http.MethodPost, pattern, Handler(handler))
}

// Before 添加拦截器
func (group *RouteGroup) Before(v ...func(*Context) Response) {
	for _, handler := range v {
		group.befores = append(group.befores, Handler(handler))
	}
}

// After 添加后置处理拦截器
func (group *RouteGroup) After(v ...func(*Context) Response) {
	for _, handler := range v {
		group.afters = append(group.afters, Handler(handler))
	}
}

// addRoute 添加路由
func (group *RouteGroup) addRoute(method, pattern string, handler IHandler) {
	pattern = path.Join(group.prefix, pattern)
	group.engine.router.AddRoute(method, pattern)
	if group.engine.handlers[method] == nil {
		group.engine.handlers[method] = make(map[string]IHandler)
	}
	group.engine.handlers[method][pattern] = handler
}

// start 启动服务 http/https
func (engine *Engine) start(server *http.Server) {
	if engine.options[optCertPath] != nil && engine.options[optKeyPath] != nil {
		certPath := engine.options[optCertPath].(CTLSCertFile)
		keyPath := engine.options[optKeyPath].(CTLSKeyFile)
		if err := server.ListenAndServeTLS(string(certPath), string(keyPath)); err != nil && err != http.ErrServerClosed {
			LogErrf("Could not listen https on %s: %v \n", server.Addr, err)
			os.Exit(1)
		}
	} else {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			LogErrf("Could not listen http on %s: %v \n", server.Addr, err)
			os.Exit(1)
		}
	}
}

// ParseConfig 读取配置文件
func (engine *Engine) parseConfig() (err error) {
	if content, ok := engine.options[optAppConfig]; ok {
		LogInfo("Loading config from content by provided")
		err = engine.config.ParseContent(content.(CAppConfig))
	} else {
		LogInfof("Loading config from %s\n", DEFAULT_CONFIG_PATH)
		err = engine.config.ParseFile(DEFAULT_CONFIG_PATH)
		if err != nil && os.IsNotExist(err) {
			LogWarn("Local config file not found, using default config")
			return nil
		}
	}
	return err
}

// ParseOptions 解析运行参数
func (engine *Engine) parseOptions(options ...any) {
	for _, option := range options {
		switch option := option.(type) {
		case CAppConfig:
			engine.options[optAppConfig] = option
		case CTLSCertFile:
			engine.options[optCertPath] = option
		case CTLSKeyFile:
			engine.options[optKeyPath] = option
		}
	}
}

// setConfig 依据配置进行初始化
func (engine *Engine) initConfig() error {
	if engine.config.Database.Source != "" {
		LogInfo("Database source found, connecting to database")
		if dataEngine, err := orm.NewEngine(engine.config.Database.Source); err != nil {
			return err
		} else {
			engine.database = dataEngine
		}
	} else {
		LogWarn("Database source not found, you may not be able to use relevant modules")
	}
	return nil
}

// shutdown 服务关闭处理
func (engine *Engine) shutdown(server *http.Server, quit <-chan os.Signal, done chan<- struct{}) {
	<-quit
	LogWarn("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		LogErrf("Cound not gracefully shutdown the server: %v \n", err)
		os.Exit(1)
	}
	if engine.database != nil {
		engine.database.Close()
		LogInfo("Database closed success")
	}
	close(done)
}

// ServeHTTP 实现 http.Handler 接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	defer handleErr(ctx)
	befores, afters := engine.findInterceptor(ctx)
	if response := handleMiddlewares(ctx, befores); response != nil {
		response.Invoke(ctx)
	} else if response := engine.handleHandler(ctx); response != nil {
		response.Invoke(ctx)
	} else if response := handleMiddlewares(ctx, afters); response != nil {
		response.Invoke(ctx)
	}
}

// handleHandler 执行处理器
func (engine *Engine) handleHandler(ctx *Context) Response {
	if handler := engine.findHandler(ctx); handler != nil {
		return handler.Invoke(ctx)
	} else {
		return String(http.StatusNotFound, "404 Not Found: %s", ctx.Path)
	}
}

// findHandler 查找处理器
func (engine *Engine) findHandler(ctx *Context) IHandler {
	if route, params := engine.router.GetRoute(ctx.Method, ctx.Path); route != "" {
		ctx.Params = params
		return engine.handlers[ctx.Method][route]
	}
	return nil
}

// findInterceptor 查找拦截器
func (engine *Engine) findInterceptor(ctx *Context) (befores []IHandler, afters []IHandler) {
	for _, group := range engine.groups {
		if strings.HasPrefix(ctx.Path, group.prefix) {
			befores = append(befores, group.befores...)
			afters = append(afters, group.afters...)
		}
	}
	return befores, afters
}
