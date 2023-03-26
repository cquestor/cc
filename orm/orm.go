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
	table       string
	sql         *strings.Builder
	storeInsert *StoreInsert
}

// StoreInsert 插入缓存
type StoreInsert struct {
	elem       reflect.Type
	fields     []string
	values     [][]any
	prepareStr *strings.Builder
	execs      []any
}

// NewEngine 构造ORM引擎
func NewEngine(source string) *Engine {
	db, err := connect(source)
	if err != nil {
		panic(err)
	}
	return &Engine{
		DB: db,
	}
}

// NewSession 构造数据库会话
func (engine *Engine) NewSession() *Session {
	return &Session{
		db:  engine.DB,
		sql: &strings.Builder{},
		storeInsert: &StoreInsert{
			fields:     make([]string, 0),
			values:     make([][]any, 0),
			prepareStr: &strings.Builder{},
			execs:      make([]any, 0),
		},
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

// Table 设置表格名
func (session *Session) Table(name string) *Session {
	session.table = name
	return session
}

// Insert 插入数据
func (session *Session) Insert(v ...any) error {
	defer session.Reset()
	if session.table == "" {
		return fmt.Errorf("table name is empty, forgot set it?")
	}
	for _, each := range v {
		fields, values, _type := parseObject(each)
		if err := session._insert(fields, values, _type); err != nil {
			return err
		}
	}
	session.sql.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", session.table, strings.Join(session.storeInsert.fields, ", "), session.storeInsert.prepareStr.String()))
	stmt, err := session.db.Prepare(session.sql.String())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(session.storeInsert.execs...)
	return err
}

// _insert 插入一条数据
func (session *Session) _insert(fields []string, values []any, _type reflect.Type) error {
	if session.storeInsert.elem == nil {
		session.storeInsert.elem = _type
	}
	if session.storeInsert.elem != _type {
		return fmt.Errorf("all inserted objects must be of the same type (%s) %s", session.storeInsert.elem.Name(), _type.Name())
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

// Reset 重置会话
func (session *Session) Reset() {
	session.table = ""
	session.storeInsert.Clear()
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
