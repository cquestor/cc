package cc

import (
	"net/http"

	"github.com/cquestor/cc/router"
)

// WebEngine Web引擎
type WebEngine struct {
	router   *router.CRouter
	handlers HandlerProvider
}

// New 构造WebEngine
func New() *WebEngine {
	return &WebEngine{
		router:   router.NewCRouter(),
		handlers: make(HandlerProvider),
	}
}

// addRoute 添加新路由和处理器
func (engine *WebEngine) addRoute(method, pattern string, handler IHandler) {
	engine.router.AddRoute(method, pattern)
	engine.handlers.AddHandler(method, pattern, handler)
}

// Get 添加 GET 方法
func (engine *WebEngine) Get(pattern string, handler func(*Context) IResponse) {
	engine.addRoute(http.MethodGet, pattern, Handler(handler))
}

// Post 添加 POST 方法
func (engine *WebEngine) Post(pattern string, handler func(*Context) IResponse) {
	engine.addRoute(http.MethodPost, pattern, Handler(handler))
}

// Delete 添加 DELETE 方法
func (engine *WebEngine) Delete(pattern string, handler func(*Context) IResponse) {
	engine.addRoute(http.MethodDelete, pattern, Handler(handler))
}

// Put 添加 PUT 方法
func (engine *WebEngine) Put(pattern string, handler func(*Context) IResponse) {
	engine.addRoute(http.MethodPut, pattern, Handler(handler))
}

// Options 添加 OPTIONS 方法
func (engine *WebEngine) Options(pattern string, handler func(*Context) IResponse) {
	engine.addRoute(http.MethodOptions, pattern, Handler(handler))
}

// ServeHTTP 实现 http.Handler 接口
func (engine *WebEngine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if route, params := engine.router.GetRoute(r.Method, r.URL.Path); route != nil {
		ctx := NewContext(w, r, params)
		engine.handlers[r.Method][route.Pattern].Invoke(ctx).Invoke(ctx)
	} else {
		ctx := NewContext(w, r, params)
		String(http.StatusNotFound, "404 Not Found: %s", r.URL.Path).Invoke(ctx)
	}
}

// Run 启动 Web 服务
func (engine *WebEngine) Run(addr string) error {
	return http.ListenAndServe(addr, engine)
}
