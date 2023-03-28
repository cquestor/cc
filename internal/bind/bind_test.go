package bind_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/cquestor/cc/internal/bind/binder"
)

func TestParse(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		binder := binder.GetBinder(reflect.String)
		var value string = "12"
		var uintTest uint
		var intTest int
		var floatTest float32
		var boolTest bool
		var stringTest string
		if err := binder.Bind(value, &uintTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(uintTest)
		if err := binder.Bind(value, &intTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(intTest)
		if err := binder.Bind(value, &floatTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(floatTest)
		if err := binder.Bind(value, &boolTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(boolTest)
		if err := binder.Bind(value, &stringTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(stringTest)
	})
	t.Run("number", func(t *testing.T) {
		binder := binder.GetBinder(reflect.Int)
		var value int = 12
		var stringTest string
		var boolTest bool
		var intTest int
		if err := binder.Bind(value, &stringTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(stringTest)
		if err := binder.Bind(value, &boolTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(boolTest)
		if err := binder.Bind(value, &intTest); err != nil {
			t.Fatal(intTest)
		}
		fmt.Println(intTest)
	})
	t.Run("bool", func(t *testing.T) {
		binder := binder.GetBinder(reflect.Bool)
		var value bool = true
		var intTest int
		var floatTest float32
		var stringTest string
		var boolTest bool
		if err := binder.Bind(value, &intTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(intTest)
		if err := binder.Bind(value, &floatTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(floatTest)
		if err := binder.Bind(value, &stringTest); err != nil {
			t.Fatal(stringTest)
		}
		fmt.Println(stringTest)
		if err := binder.Bind(value, &boolTest); err != nil {
			t.Fatal(err)
		}
		fmt.Println(boolTest)
	})
}
