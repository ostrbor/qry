package qry

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	tableTag     = "table"
	columnTag    = "column"
	generatedTag = "generated"
)

type dbTable struct {
	name    string
	columns []column
}

type column struct {
	name           string
	generated      bool // indicates if the column is generated in database (serial or identity)
	structFieldPtr interface{}
}

func names(cols []column) []string {
	res := make([]string, 0, len(cols))
	for _, c := range cols {
		res = append(res, c.name)
	}
	return res
}

func pointers(cols []column) []interface{} {
	res := make([]interface{}, 0, len(cols))
	for _, c := range cols {
		res = append(res, c.structFieldPtr)
	}
	return res
}

func values(cols []column) (res []interface{}) {
	for _, col := range cols {
		// todo check if can be done without reflect
		res = append(res, reflect.ValueOf(col.structFieldPtr).Elem().Interface())
	}
	return res
}

// model can be struct or pointer to struct.
func parseModel(model reflect.Value) (t dbTable, err error) {
	s := reflect.Indirect(model)
	if s.Kind() != reflect.Struct {
		return t, errKind(s.Kind(), reflect.Struct)
	}

	var tableName string
	var cols []column
	var generatedCount, tableCount int
	for i := 0; i < s.NumField(); i++ {
		f := s.Type().Field(i)
		colName, generated := parseColumn(f.Tag)
		if colName == "" {
			continue
		}
		cols = append(cols, column{
			name:           colName,
			generated:      generated,
			structFieldPtr: s.Field(i).Addr().Interface(),
		})
		if generated {
			generatedCount++
		}
		tab := f.Tag.Get(tableTag)
		if tab != "" {
			tableName = tab
			tableCount++
		}
	}

	if len(cols) == 0 {
		return t, errors.New("missing 'column' in tags")
	}
	if tableCount == 0 {
		return t, errors.New("missing 'table' in tags")
	} else if tableCount > 1 {
		return t, fmt.Errorf("must be one 'table', got %d", tableCount)
	}
	if generatedCount > 1 {
		return t, fmt.Errorf("must be one 'generated', got %d", generatedCount)
	}

	t.name = tableName
	t.columns = cols
	return t, nil
}

// parseColumn is a function that processes a struct tag from a Go struct field.
// It specifically looks for the 'column' tag and extracts two pieces of information:
// 1. The name of the database column that the struct field maps to.
// 2. A boolean flag indicating whether the database column is auto-generated (e.g., an auto-incrementing ID).
func parseColumn(t reflect.StructTag) (name string, gen bool) {
	c := t.Get(columnTag)
	name = strings.Split(c, ",")[0]
	gen = strings.Contains(c, generatedTag)
	return name, gen
}
