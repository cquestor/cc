package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// CTx 事务
type CTx struct {
	tx          *sql.Tx
	table       []string
	sql         *strings.Builder
	storeInsert *StoreInsert
	storeWhere  []*StoreWhere
	storeSet    []*StoreSet
	storeLimit  *StoreLimit
	storeOrder  []*StoreOrder
	lastExec    sql.Result
}

func (tx *CTx) LastExec() sql.Result {
	return tx.lastExec
}

// Table 设置表格名
func (tx *CTx) Table(name string) *CTx {
	tx.table = append(tx.table, "`"+name+"`")
	return tx
}

// Where 添加 Where 子句
func (tx *CTx) Where(field, flag string, v any) *CTx {
	tx.storeWhere = append(tx.storeWhere, &StoreWhere{prepareStr: fmt.Sprintf("%s%s?", field, flag), exec: v})
	return tx
}

// Equal 相等 Where 子句
func (tx *CTx) Equal(field string, v any) *CTx {
	return tx.Where(field, "=", v)
}

// Unequal 不相等 Where 子句
func (tx *CTx) Unequal(field string, v any) *CTx {
	return tx.Where(field, "!=", v)
}

// Set 添加 Set 子句
func (tx *CTx) Set(field string, v any) *CTx {
	tx.storeSet = append(tx.storeSet, &StoreSet{prepareStr: fmt.Sprintf("%s=?", field), exec: v})
	return tx
}

// Limit 添加 Limit 子句
func (tx *CTx) Limit(count int, offset ...int) *CTx {
	if len(offset) < 1 {
		offset = append(offset, 0)
	}
	tx.storeLimit.offset = offset[0]
	tx.storeLimit.count = count
	return tx
}

// Order 添加 Order By 子句
func (tx *CTx) Order(name string, desc ...bool) *CTx {
	if len(desc) < 1 {
		desc = append(desc, false)
	}
	tx.storeOrder = append(tx.storeOrder, &StoreOrder{name: name, desc: desc[0]})
	return tx
}

// Update 修改数据
func (tx *CTx) Update() error {
	defer tx.Reset()
	if len(tx.table) < 1 {
		return fmt.Errorf("table name is empty, forgot set it?")
	}
	tx.sql.WriteString(fmt.Sprintf("UPDATE %s", tx.table[0]))
	execs := make([]any, 0)
	addSets(tx.sql, tx.storeSet, &execs)
	addWheres(tx.sql, tx.storeWhere, &execs)
	return tx._exec(execs...)
}

// Delete 删除数据
func (tx *CTx) Delete() error {
	defer tx.Reset()
	if len(tx.table) < 1 {
		return fmt.Errorf("table name is empty, forgot set it?")
	}
	tx.sql.WriteString(fmt.Sprintf("DELETE FROM %s", tx.table[0]))
	execs := make([]any, 0)
	addWheres(tx.sql, tx.storeWhere, &execs)
	return tx._exec(execs...)
}

// Insert 插入数据
func (tx *CTx) Insert(v ...any) error {
	defer tx.Reset()
	if len(tx.table) < 1 {
		return fmt.Errorf("table name is empty, forgot set it?")
	}
	for _, each := range v {
		each, err := checkInsertObject(each)
		if err != nil {
			return err
		}
		fields, values, _type := parseInsertObject(each)
		if err := tx._insert(fields, values, _type); err != nil {
			return err
		}
	}
	tx.sql.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tx.table[0], strings.Join(tx.storeInsert.fields, ", "), tx.storeInsert.prepareStr.String()))
	return tx._exec(tx.storeInsert.execs...)
}

// _insert 插入一条数据
func (tx *CTx) _insert(fields []string, values []any, _type reflect.Type) error {
	if tx.storeInsert.elem == nil {
		tx.storeInsert.elem = _type
	}
	if tx.storeInsert.elem != _type {
		return fmt.Errorf("all inserted objects must be of the same type (%s) <=> (%s)", tx.storeInsert.elem.Name(), _type.Name())
	}
	tx.storeInsert.fields = fields
	tx.storeInsert.values = append(tx.storeInsert.values, values)
	tx.storeInsert.execs = append(tx.storeInsert.execs, values...)
	if tx.storeInsert.prepareStr.String() != "" {
		tx.storeInsert.prepareStr.WriteString(", ")
	}
	tx.storeInsert.prepareStr.WriteString(fmt.Sprintf("(%s)", strings.Join(genPrepare(len(values)), ", ")))
	return nil
}

// Select 查询数据
func (tx *CTx) Select(v any) (int, error) {
	defer tx.Reset()
	if len(tx.table) < 1 {
		return 0, fmt.Errorf("table name is empty, forgot set it?")
	}
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("the select target object must be a ptr, not %s", t.Kind())
	}
	_type, isSlice, isPointer, err := checkSelectObject(v)
	if err != nil {
		return 0, err
	}
	if isSlice {
		rows, err := tx._select(v, _type)
		if err != nil {
			return 0, err
		}
		n, values, err := tx._select_fetch(_type, rows)
		if err != nil {
			return 0, err
		}
		_value := reflect.ValueOf(v)
		for _, each := range values {
			if isPointer {
				_value.Elem().Set(reflect.Append(_value.Elem(), each))
			} else {
				_value.Elem().Set(reflect.Append(_value.Elem(), each.Elem()))
			}
		}
		return n, nil
	}
	tx.Limit(1)
	rows, err := tx._select([]any{v}, _type)
	if err != nil {
		return 0, err
	}
	return tx._select_one(v, rows)
}

// _select 查询数据
func (tx *CTx) _select(v any, t reflect.Type) (*sql.Rows, error) {
	execs := make([]any, 0)
	fields := parseSelectObject(t)
	tx.sql.WriteString(fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), strings.Join(tx.table, ", ")))
	addWheres(tx.sql, tx.storeWhere, &execs)
	addOrders(tx.sql, tx.storeOrder, &execs)
	addLimit(tx.sql, tx.storeLimit, &execs)
	stmt, err := tx.tx.Prepare(tx.sql.String())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(execs...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// _select_one 查询单条数据
func (tx *CTx) _select_one(v any, rows *sql.Rows) (int, error) {
	if !rows.Next() {
		return 0, nil
	}
	if err := scanOneObject(v, rows); err != nil {
		return 0, err
	}
	return 1, nil
}

// _select_fetch 查询多条数据
func (tx *CTx) _select_fetch(t reflect.Type, rows *sql.Rows) (int, []reflect.Value, error) {
	effectCount := 0
	values := make([]reflect.Value, 0)
	if !rows.Next() {
		return 0, nil, nil
	} else {
		effectCount++
		if value, err := scanFetchObject(t, rows); err != nil {
			return 0, nil, err
		} else {
			values = append(values, value)
		}
		for rows.Next() {
			effectCount++
			if value, err := scanFetchObject(t, rows); err != nil {
				return 0, nil, err
			} else {
				values = append(values, value)
			}
		}
	}
	return effectCount, values, nil
}

// _exec 执行 sql 语句
func (tx *CTx) _exec(execs ...any) error {
	stmt, err := tx.tx.Prepare(tx.sql.String())
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(execs...)
	tx.lastExec = res
	return err
}

// Reset 重置会话
func (tx *CTx) Reset() {
	tx.table = tx.table[:0]
	tx.storeInsert.Clear()
	tx.storeWhere = tx.storeWhere[:0]
	tx.storeSet = tx.storeSet[:0]
	tx.storeLimit.Clear()
	tx.sql.Reset()
}

// Commit 提交事务
func (tx *CTx) Commit() error {
	return tx.tx.Commit()
}

// Rollback 回滚事务
func (tx *CTx) Rollback() error {
	return tx.tx.Rollback()
}
