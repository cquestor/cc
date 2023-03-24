package cc_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/cquestor/cc"
)

type Test struct{}

func TestRoute(t *testing.T) {
	c := cc.New()
	c.Get("/", func(ctx *cc.Context) cc.IResponse {
		return cc.String(http.StatusOK, "你好!")
	})
	c.Get("/html", func(ctx *cc.Context) cc.IResponse {
		return cc.Html(http.StatusOK, []byte("<h1>Hello World!</h1>"))
	})
	c.Run(":9999")
}

func TestJson(t *testing.T) {
	type User struct {
		User string `json:"name"`
	}
	user := User{User: "chen"}
	data, _ := json.Marshal(user)
	fmt.Println(string(data))
}
