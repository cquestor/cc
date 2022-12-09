package cc_test

import (
	"fmt"
	"testing"

	"github.com/cquestor/cc"
)

func TestMain(m *testing.M) {
	// c := cc.Default()

	// c.GET("/", func(c *cc.Context) {
	// })

	// c.Run()
	type user struct {
		Name string
	}
	result, err := cc.JWTToken(user{Name: "nihao"}, "cquestor")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	fmt.Println(cc.JWTCheck(result, "cquestor"))
	var test user
	err = cc.JWTParse(result, &test)
	if err != nil {
		panic(err)
	}
	fmt.Println(test.Name)
}
