package cc

import (
	"go/ast"
	"reflect"
	"strings"
)

type field struct {
	Name   string
	Type   string
	Tag    string
	Source string
}

type schema struct {
	Model       any
	Name        string
	Fields      []*field
	FieldNames  []string
	SourceNames []string
	fieldMap    map[string]*field
}

func (schema *schema) GetField(name string) *field {
	return schema.fieldMap[name]
}

func parseTable(dest any, d dialect) *schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &schema{
		Model:    dest,
		Name:     camel2Case(modelType.Name()),
		fieldMap: make(map[string]*field),
	}
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &field{
				Name:   camel2Case(p.Name),
				Type:   d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
				Source: p.Name,
			}
			if v, ok := p.Tag.Lookup("cc"); ok {
				infos := strings.Split(v, ",")
				field.Name = infos[0]
				if len(infos) > 1 {
					field.Tag = strings.Join(infos[1:], " ")
				}
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, field.Name)
			schema.SourceNames = append(schema.SourceNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (schema *schema) RecordValues(dest any) []any {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []any
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Source).Interface())
	}
	return fieldValues
}
