package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// connect 连接数据库
func connect(source string) (*sql.DB, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// genPrepare 生成占位符
func genPrepare(v int) []string {
	temp := make([]string, v)
	for i := 0; i < v; i++ {
		temp[i] = "?"
	}
	return temp
}

// addWheres 添加 where 子句
func addWheres(sql *strings.Builder, wheres []*StoreWhere, execs *[]any) {
	if len(wheres) > 0 {
		sql.WriteString(" WHERE ")
	}
	for i, e := range wheres {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(e.prepareStr)
		*execs = append(*execs, e.exec)
	}
}

// addSets 添加 set 子句
func addSets(sql *strings.Builder, sets []*StoreSet, execs *[]any) {
	if len(sets) > 0 {
		sql.WriteString(" SET ")
	}
	for i, e := range sets {
		if i > 0 {
			sql.WriteString(", ")
		}
		sql.WriteString(e.prepareStr)
		*execs = append(*execs, e.exec)
	}
}

// addLimit 添加 limit 语句
func addLimit(sql *strings.Builder, limit *StoreLimit, execs *[]any) {
	if limit.count > 0 {
		sql.WriteString(" LIMIT ?, ?")
		*execs = append(*execs, limit.offset)
		*execs = append(*execs, limit.count)
	}
}

// addOrders 添加 order by 语句
func addOrders(sql *strings.Builder, orders []*StoreOrder, execs *[]any) {
	if len(orders) > 0 {
		sql.WriteString(" ORDER BY ")
	}
	for i, e := range orders {
		if i > 0 {
			sql.WriteString(", ")
		}
		if e.desc {
			sql.WriteString(fmt.Sprintf("%s DESC", e.name))
		} else {
			sql.WriteString(fmt.Sprintf("%s ASC", e.name))
		}
	}
}

// checkInsertObject 检查插入元素
func checkInsertObject(v any) (any, error) {
	t := reflect.TypeOf(v)
	e := reflect.ValueOf(v)
	if t.Kind() == reflect.Ptr {
		return checkInsertObject(e.Elem().Interface())
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("all insert objects must be a struct, not %s", t.Kind())
	}
	return v, nil
}

// checkSelectObject 检查查询元素
func checkSelectObject(v any) (reflect.Type, bool, bool, error) {
	var isSlice bool
	var isPointer bool
	t := reflect.ValueOf(v).Elem().Type()
	if t.Kind() == reflect.Slice {
		isSlice = true
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		isPointer = true
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, false, false, fmt.Errorf("the select target must be a ptr of struct, not %s", t.Kind())
	}
	return t, isSlice, isPointer, nil
}
