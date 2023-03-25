// Copyright 2023 cquestor. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Engine Web引擎
type Engine struct {
	router map[string]func()
}

// New 构造Engine
func New() *Engine {
	return &Engine{
		router: make(map[string]func()),
	}
}

// Run 启动 Web Server
func (engine *Engine) Run(port int) {
	listenAddr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:    listenAddr,
		Handler: engine,
		// TODO 自定义错误日志
		ErrorLog:     log.Default(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go engine.shutdown(server, quit, done)
	log.Println("Server is ready to handle requests at", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v \n", listenAddr, err)
	}
	<-done
	log.Println("Server stopped")
}

// shutdown 服务关闭处理
func (engine *Engine) shutdown(server *http.Server, quit <-chan os.Signal, done chan<- struct{}) {
	<-quit
	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Cound not gracefully shutdown the server: %v \n", err)
	}
	// TODO do something, such as close database connection...
	close(done)
}

// ServeHTTP 实现 http.Handler 接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
