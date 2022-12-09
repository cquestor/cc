package cc_test

import (
	"testing"

	"github.com/cquestor/cc"
)

func TestMain(m *testing.M) {
	c := cc.Default()

	c.GET("/", func(c *cc.Context) {
	})

	c.Run()
}
