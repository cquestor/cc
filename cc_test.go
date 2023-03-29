package cc_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	_ "embed"

	"github.com/cquestor/cc"
	"github.com/cquestor/cc/middleware"
)

func TestConfig(t *testing.T) {
	config := cc.AppConfig{}
	if err := config.ParseFile("1.txt"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("不存在")
		} else {
			t.Fatal(err)
		}
	}
}

//go:embed test.yaml
var content []byte

func TestRun(t *testing.T) {
	c := cc.New()

	cors := middleware.CorsMiddleware{}
	clog := middleware.CLogMiddleware{}

	c.Before(cors.Instance(), clog.Instance())

	c.Get("/", func(ctx *cc.Context) cc.Response {
		return cc.String(http.StatusOK, "success")
	})

	c.Get("/hello", func(ctx *cc.Context) cc.Response {
		return cc.Html(http.StatusOK, []byte("<h1>Hello CC!</h1>"))
	})

	c.Get("/panic", func(ctx *cc.Context) cc.Response {
		panic("recovery test")
	})

	age := 10

	user := c.Group("/user")
	user.Before(func(ctx *cc.Context) cc.Response {
		age += 90
		return nil
	})
	user.Get("age", func(ctx *cc.Context) cc.Response {
		return cc.String(http.StatusOK, "%d岁\n", age)
	})

	c.Run(cc.CAppConfig(content))
}
