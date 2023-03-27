package cc_test

import (
	"fmt"
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
	c.Run(cc.CAppConfig(content))
}
