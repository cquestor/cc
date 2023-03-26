package orm_test

import (
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
	if err := session.Table("account").Insert(Account{Name: "admin", Age: 100}, Account{Name: "chen", Age: 23}); err != nil {
		t.Fatal(err)
	}
}
