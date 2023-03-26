package orm_test

import (
	"fmt"
	"testing"

	"github.com/cquestor/cc/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	Id   int    `cc:"id,auto_increment"`
	Name string `cc:"name"`
	Age  int    `cc:"age"`
}

func TestData(t *testing.T) {
	data := orm.NewEngine("root:software@tcp(localhost:3306)/test")
	defer data.Close()
	session := data.NewSession()
	t.Run("insert", func(t *testing.T) {
		if err := session.Table("account").Insert(Account{Name: "admin", Age: 100}, &Account{Name: "chen", Age: 23}); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("update", func(t *testing.T) {
		if err := session.Table("account").Equal("name", "admin").Set("age", 23).Update(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("tx", func(t *testing.T) {
		tx, err := session.GetTx()
		if err != nil {
			t.Fatal(err)
		}
		if err := tx.Table("account").Equal("name", "admin").Set("name", "张三").Update(); err != nil {
			tx.Rollback()
			t.Fatal(err)
		}
		if err := tx.Table("account").Equal("name", "chen").Set("name", "李四").Update(); err != nil {
			tx.Rollback()
			t.Fatal(err)
		}
		if err := tx.Commit(); err != nil {
			tx.Rollback()
			t.Fatal(err)
		}
	})
	t.Run("select", func(t *testing.T) {
		var account []Account
		if n, err := session.Table("account").Equal("age", 23).Order("id", true).Select(&account); err != nil {
			t.Fatal(err)
		} else {
			fmt.Println("select count:", n)
			fmt.Println(account)
		}
	})
	t.Run("delete", func(t *testing.T) {
		if err := session.Table("account").Equal("age", 23).Delete(); err != nil {
			t.Fatal(err)
		}
	})
}
