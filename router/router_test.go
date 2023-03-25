package router_test

import (
	"testing"

	"github.com/cquestor/cc/router"
)

func TestRouter(t *testing.T) {
	router := router.NewRouter()
	// 绝对路由
	t.Run("absolute", func(t *testing.T) {
		router.AddRoute("GET", "/index")
		pattern, params := router.GetRoute("GET", "/index")
		if pattern != "/index" {
			t.Fatalf("absolute route parse error: %s %+v\n", pattern, params)
		}
	})
	// 动态路由
	t.Run("dynamic", func(t *testing.T) {
		router.AddRoute("GET", "/user/:name")
		pattern, params := router.GetRoute("GET", "/user/admin")
		if pattern != "/user/:name" || params["name"] != "admin" {
			t.Fatalf("dynamic route parse error: %s %+v\n", pattern, params)
		}
	})
	// 通配符路由
	t.Run("wild", func(t *testing.T) {
		router.AddRoute("GET", "/static/*file")
		pattern, params := router.GetRoute("GET", "/static/test.jpg")
		if pattern != "/static/*file" || params["file"] != "test.jpg" {
			t.Fatalf("wild route parse error: %s %+v\n", pattern, params)
		}
	})
}
