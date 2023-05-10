// ORM framework
package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Engine ORM引擎
type Engine struct {
	DB *sql.DB
}

// Session 数据库会话
type Session struct {
	db          *sql.DB
	table       []string
	sql         *strings.Builder
	storeInsert *StoreInsert
	storeWhere  []*StoreWhere
	storeSet    []*StoreSet
	storeLimit  *StoreLimit
	storeOrder  []*StoreOrder
}

// StoreInsert 插入缓存
type StoreInsert struct {
	elem       reflect.Type
	fields     []string
	values     [][]any
	prepareStr *strings.Builder
	execs      []any
}

// StoreWhere where缓存
type StoreWhere struct {
	prepareStr string
	exec       any
}

// StoreSet set缓存
type StoreSet struct {
	prepareStr string
	exec       any
}

// StoreLimit limit缓存
type StoreLimit struct {
	offset int
	count  int
}

// StoreOrder orderby缓存
type StoreOrder struct {
	name string
	desc bool
}

// NewEngine 构造ORM引擎
func NewEngine(source string) (*Engine, error) {
	db, err := connect(source)
	if err != nil {
		return nil, err
	}
	return &Engine{
		DB: db,
	}, nil
}

// NewSession 构造数据库会话
func (engine *Engine) NewSession() *Session {
	return &Session{
		db:    engine.DB,
		table: make([]string, 0),
		sql:   &strings.Builder{},
		storeInsert: &StoreInsert{
			fields:     make([]string, 0),
			values:     make([][]any, 0),
			prepareStr: &strings.Builder{},
			execs:      make([]any, 0),
		},
		storeWhere: make([]*StoreWhere, 0),
		storeSet:   make([]*StoreSet, 0),
		storeLimit: &StoreLimit{},
		storeOrder: make([]*StoreOrder, 0),
	}
}

// SetMaxOpenConns 设置打开的最大连接数
func (engine *Engine) SetMaxOpenConns(v int) {
	engine.DB.SetMaxOpenConns(v)
}

// SetMaxIdleConns 设置池中最大空闲连接数，即保留连接以备下次使用
func (engine *Engine) SetMaxIdleConns(v int) {
	engine.DB.SetMaxIdleConns(v)
}

// Close 关闭数据库连接
func (engine *Engine) Close() {
	engine.DB.Close()
}

// GetTx 获取事务
func (session *Session) Begin() (*CTx, error) {
	tx, err := session.db.Begin()
	if err != nil {
		return nil, err
	}
	return &CTx{
		tx:    tx,
		table: make([]string, 0),
		sql:   &strings.Builder{},
		storeInsert: &StoreInsert{
			fields:     make([]string, 0),
			values:     make([][]any, 0),
			prepareStr: &strings.Builder{},
			execs:      make([]any, 0),
		},
		storeWhere: make([]*StoreWhere, 0),
		storeSet:   make([]*StoreSet, 0),
		storeLimit: &StoreLimit{},
		storeOrder: make([]*StoreOrder, 0),
	}, nil
}

// Table 设置表格名
func (session *Session) Table(name string) *Session {
	session.table = append(session.table, name)
	return session
}

// Where 添加 Where 子句
func (session *Session) Where(field, flag string, v any) *Session {
	session.storeWhere = append(session.storeWhere, &StoreWhere{prepareStr: fmt.Sprintf("%s %s ?", field, flag), exec: v})
	return session
}

// Equal 相等 Where 子句
func (session *Session) Equal(field string, v any) *Session {
	return session.Where(field, "=", v)
}

// Unequal 不相等 Where 子句
func (session *Session) Unequal(field string, v any) *Session {
	return session.Where(field, "!=", v)
}

// Set 添加 Set 子句
func (session *Session) Set(field string, v any) *Session {
	session.storeSet = append(session.storeSet, &StoreSet{prepareStr: fmt.Sprintf("%s = ?", field), exec: v})
	return session
}

// Limit 添加 Limit 子句
func (session *Session) Limit(count int, offset ...int) *Session {
	if len(offset) < 1 {
		offset = append(offset, 0)
	}
	session.storeLimit.offset = offset[0]
	session.storeLimit.count = count
	return session
}

// Order 添加 Order By 子句
func (session *Session) Order(name string, desc ...bool) *Session {
	if len(desc) < 1 {
		desc = append(desc, false)
	}
	session.storeOrder = append(session.storeOrder, &StoreOrder{name: name, desc: desc[0]})
	return session
}

// Update 修改数据
func (session *Session) Update() error {
	defer session.Reset()
	if len(session.table) < 1 {
		return fmt.Errorf("table name is empty, forgot set it?")
	}
	session.sql.WriteString(fmt.Sprintf("UPDATE %s", session.table[0]))
	execs := make([]any, 0)
	addSets(session.sql, session.storeSet, &execs)
	addWheres(session.sql, session.storeWhere, &execs)
	return session._exec(execs...)
}

// Delete 删除数据
func (session *Session) Delete() error {
	defer session.Reset()
	if len(session.table) < 1 {
		return fmt.Errorf("table name is empty, forgot set it?")
	}
	session.sql.WriteString(fmt.Sprintf("DELETE FROM %s", session.table[0]))
	execs := make([]any, 0)
	addWheres(session.sql, session.storeWhere, &execs)
	return session._exec(execs...)
}

// Insert 插入数据
func (session *Session) Insert(v ...any) error {
	defer session.Reset()
	if len(session.table) < 1 {
		return fmt.Errorf("table name is empty, forgot set it?")
	}
	for _, each := range v {
		each, err := checkInsertObject(each)
		if err != nil {
			return err
		}
		fields, values, _type := parseInsertObject(each)
		if err := session._insert(fields, values, _type); err != nil {
			return err
		}
	}
	session.sql.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", session.table[0], strings.Join(session.storeInsert.fields, ", "), session.storeInsert.prepareStr.String()))
	return session._exec(session.storeInsert.execs...)
}

// _insert 插入一条数据
func (session *Session) _insert(fields []string, values []any, _type reflect.Type) error {
	if session.storeInsert.elem == nil {
		session.storeInsert.elem = _type
	}
	if session.storeInsert.elem != _type {
		return fmt.Errorf("all inserted objects must be of the same type (%s) <=> (%s)", session.storeInsert.elem.Name(), _type.Name())
	}
	session.storeInsert.fields = fields
	session.storeInsert.values = append(session.storeInsert.values, values)
	session.storeInsert.execs = append(session.storeInsert.execs, values...)
	if session.storeInsert.prepareStr.String() != "" {
		session.storeInsert.prepareStr.WriteString(", ")
	}
	session.storeInsert.prepareStr.WriteString(fmt.Sprintf("(%s)", strings.Join(genPrepare(len(values)), ", ")))
	return nil
}

// Select 查询数据
func (session *Session) Select(v any) (int, error) {
	defer session.Reset()
	if len(session.table) < 1 {
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
		rows, err := session._select(v, _type)
		if err != nil {
			return 0, err
		}
		n, values, err := session._select_fetch(_type, rows)
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
	session.Limit(1)
	rows, err := session._select([]any{v}, _type)
	if err != nil {
		return 0, err
	}
	return session._select_one(v, rows)
}

// _select 查询数据
func (session *Session) _select(v any, t reflect.Type) (*sql.Rows, error) {
	execs := make([]any, 0)
	fields := parseSelectObject(t)
	session.sql.WriteString(fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), strings.Join(session.table, ", ")))
	addWheres(session.sql, session.storeWhere, &execs)
	addOrders(session.sql, session.storeOrder, &execs)
	addLimit(session.sql, session.storeLimit, &execs)
	stmt, err := session.db.Prepare(session.sql.String())
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
func (session *Session) _select_one(v any, rows *sql.Rows) (int, error) {
	defer rows.Close()
	if !rows.Next() {
		return 0, nil
	}
	if err := scanOneObject(v, rows); err != nil {
		return 0, err
	}
	return 1, nil
}

// _select_fetch 查询多条数据
func (session *Session) _select_fetch(t reflect.Type, rows *sql.Rows) (int, []reflect.Value, error) {
	defer rows.Close()
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
func (session *Session) _exec(execs ...any) error {
	stmt, err := session.db.Prepare(session.sql.String())
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(execs...)
	return err
}

// Reset 重置会话
func (session *Session) Reset() {
	session.table = session.table[:0]
	session.storeInsert.Clear()
	session.storeWhere = session.storeWhere[:0]
	session.storeSet = session.storeSet[:0]
	session.storeLimit.Clear()
	session.sql.Reset()
}

// Clear 清除插入缓存
func (store *StoreInsert) Clear() {
	store.elem = nil
	store.fields = store.fields[:0]
	store.values = store.values[:0]
	store.prepareStr.Reset()
	store.execs = store.execs[:0]
}

// Clear 重置 limit 缓存
func (store *StoreLimit) Clear() {
	store.offset = 0
	store.count = 0
}
