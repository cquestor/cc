package cc

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type H map[string]any

type HandlerFunc func(*Context)

type engine struct {
	*routerGroup
	router  *router
	groups  []*routerGroup
	db      *sql.DB
	dialect dialect
}

type routerGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *routerGroup
	engine      *engine
}

var (
	AppConfig     appConfig = appConfig{}
	AppConfigPath string    = "application.yaml"
	sigint                  = make(chan os.Signal, 1)
	flag_database           = make(chan *sql.DB, 1)
)

func init() {
	fmt.Println(`   __________                  __            `)
	fmt.Println(`  / ____/ __ \__  _____  _____/ /_____  _____`)
	fmt.Println(` / /   / / / / / / / _ \/ ___/ __/ __ \/ ___/`)
	fmt.Println(`/ /___/ /_/ / /_/ /  __(__  ) /_/ /_/ / /    `)
	fmt.Println(`\____/\___\_\__,_/\___/____/\__/\____/_/     `)
	fmt.Println(`                                             `)
}

func New() *engine {
	newEngine := &engine{router: newRouter()}
	newEngine.routerGroup = &routerGroup{engine: newEngine}
	newEngine.groups = []*routerGroup{newEngine.routerGroup}
	return newEngine
}

func Default() *engine {
	newEngine := New()
	newEngine.Use(Logger(), Recovery())
	return newEngine
}

func (group *routerGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *routerGroup) Group(prefix string) *routerGroup {
	parentEngine := group.engine
	newGroup := &routerGroup{
		prefix: prefix,
		parent: group,
		engine: parentEngine,
	}
	parentEngine.groups = append(parentEngine.groups, newGroup)
	return newGroup
}

func (group *routerGroup) addRoute(method, pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *routerGroup) GET(pattern string, handler HandlerFunc) {
	group.engine.addRoute("GET", pattern, handler)
}

func (group *routerGroup) POST(pattern string, handler HandlerFunc) {
	group.engine.addRoute("POST", pattern, handler)
}

func (engine *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r, engine.db, engine.dialect)
	c.handlers = middlewares
	engine.router.handle(c)
}

func (engine *engine) Run() {
	startTime := time.Now()
	AppConfig.Default()
	ok, err := AppConfig.IsExist(AppConfigPath)
	if err != nil {
		panic(err.Error())
	}
	if ok {
		AppConfig.Parse(AppConfigPath)
		Info("Using application config from " + AppConfigPath)
	}
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", AppConfig.Server.Port),
		Handler: engine,
	}
	go serverStop(engine, &srv)
	go databaseInit(engine)
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		log.Fatalf(err.Error())
		sigint <- syscall.SIGTERM
	}
	engine.db = <-flag_database
	Infof("Listening and Servicing on :%d in %v", AppConfig.Server.Port, time.Since(startTime))
	err = srv.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		Errorf("Server error: %v", err)
	} else {
		Infof("Successfully shutdown: %v", err)
	}
}

func serverStop(engine *engine, srv *http.Server) {
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	data := <-sigint
	Info("Received signal: " + data.String())
	if engine.db != nil {
		if err := engine.db.Close(); err != nil {
			Error("Failed to close database")
		} else {
			Info("Close database success")
		}
	}
	if err := srv.Shutdown(context.Background()); err != nil {
		Errorf("HTTP server shutdown: %v", err)
	}
}

func databaseInit(engine *engine) {
	if !AppConfig.Database.Logger {
		dbLogger.SetOutput(io.Discard)
	}
	if AppConfig.Database.DriverName != "" && AppConfig.Database.Url != "" {
		Info("Database config found, try to connect...")
		if !strings.Contains(AppConfig.Database.Url, "timeout") {
			if !strings.Contains(AppConfig.Database.Url, "?") {
				AppConfig.Database.Url += "?timeout=2s"
			} else {
				AppConfig.Database.Url += "&timeout=2s"
			}
		}
		db, err := sql.Open(AppConfig.Database.DriverName, AppConfig.Database.Url)
		if err != nil {
			panic(err.Error())
		}
		if err = db.Ping(); err != nil {
			Error("Connect database failed")
			panic(err.Error())
		}
		dial, ok := getDialect(AppConfig.Database.DriverName)
		if !ok {
			Errorf("Dialect %s not found", AppConfig.Database.DriverName)
			sigint <- syscall.SIGTERM
		}
		engine.dialect = dial
		Info("Connect database success")
		flag_database <- db
	} else {
		Info("Database config not found, ignore")
		flag_database <- nil
	}
}
