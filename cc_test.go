package cc_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	_ "embed"

	"github.com/cquestor/cc"
	_ "github.com/go-sql-driver/mysql"
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

	c.Get("/", func(ctx *cc.Context) cc.Response {
		return cc.String(http.StatusOK, "success")
	})

	c.Get("/hello", func(ctx *cc.Context) cc.Response {
		return cc.Html(http.StatusOK, []byte("<h1>Hello CC!</h1>"))
	})

	c.Get("/panic", func(ctx *cc.Context) cc.Response {
		panic("recovery test")
	})

	c.Run(cc.CAppConfig(content))
}
