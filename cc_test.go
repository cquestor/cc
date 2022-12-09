package cc_test

import (
	"fmt"
	"testing"

	"github.com/cquestor/cc"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestMain(m *testing.M) {
	c := cc.Default()

	c.Use(cc.Cors(nil))

	c.GET("/", func(c *cc.Context) {
		c.HTML(200, "<h1>Hello CC!</h1>")
	})

	c.POST("/test", func(c *cc.Context) {
		var user User
		c.Bind(&user)
		fmt.Println(user)
		c.JSON(200, cc.H{
			"msg": "success",
		})
	})

	c.Run()
}
