package cc

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type session struct {
	db       *sql.DB
	dialect  dialect
	refTable *schema
	clause   clause
	sql      strings.Builder
	sqlVars  []any
}

func newSession(db *sql.DB, dialect dialect) *session {
	return &session{db: db, dialect: dialect}
}

func (s *session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause{}
}

func (s *session) DB() *sql.DB {
	return s.db
}

func (s *session) Raw(sql string, values ...any) *session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	dbLog(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		Error(err.Error())
	}
	return
}

func (s *session) QueryRow() *sql.Row {
	defer s.Clear()
	dbLog(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	dbLog(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		Error(err.Error())
	}
	return
}

func (s *session) Model(value any) *session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = parseTable(value, s.dialect)
	}
	return s
}

func (s *session) RefTable() *schema {
	if s.refTable == nil {
		Error("Model is not set")
	}
	return s.refTable
}

func (s *session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

func (s *session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.refTable.Name)).Exec()
	return err
}

func (s *session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	row.Scan(&tmp)
	return tmp == s.refTable.Name
}

func (s *session) Insert(values ...any) (int64, error) {
	recordValues := make([]any, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(VALUES, recordValues...)
	sql, vars := s.clause.Build(INSERT, VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *session) Find(values any) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	s.clause.Set(SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []any
		for _, name := range table.SourceNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *session) First(value any) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("not found")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

func (s *session) Update(kv ...any) (int64, error) {
	m, ok := kv[0].(map[string]any)
	if !ok {
		m = make(map[string]any)
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(UPDATE, WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *session) Delete() (int64, error) {
	s.clause.Set(DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(DELETE, WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *session) Count() (int64, error) {
	s.clause.Set(COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(COUNT, WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *session) Limit(num int) *session {
	s.clause.Set(LIMIT, num)
	return s
}

func (s *session) Where(desc string, args ...any) *session {
	var vars []any
	s.clause.Set(WHERE, append(append(vars, desc), args...)...)
	return s
}

func (s *session) OrderBy(desc string) *session {
	s.clause.Set(ORDERBY, desc)
	return s
}
