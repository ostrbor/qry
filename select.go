package qry

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func MustSelect(db Querier, ptr interface{}, where string, params ...interface{}) {
	if err := Select(db, ptr, where, params...); err != nil {
		panic(err)
	}
}

func Select(db Querier, ptr interface{}, where string, params ...interface{}) (err error) {
	v, err := dereference(ptr)
	if err != nil {
		return err
	}
	switch v.Kind() {
	case reflect.Struct:
		return selectOne(db, v, where, params...)
	case reflect.Slice:
		elemType := v.Type().Elem()
		if err := validateStructType(elemType); err != nil {
			return err
		}
		if v.Len() != 0 {
			return fmt.Errorf("slice len: %d != 0", v.Len())
		}
		return selectMany(db, v, where, params...)
	}
	return fmt.Errorf("invalid kind: %s", v.Kind())
}

func selectOne(db Querier, model reflect.Value, where string, params ...interface{}) error {
	table, err := parseModel(model)
	if err != nil {
		return err
	}
	sql := buildSelectSQL(table, where)
	row := db.QueryRow(context.Background(), sql, params...)
	return row.Scan(pointers(table.columns)...)
}

func selectMany(db Querier, models reflect.Value, where string, params ...interface{}) error {
	elemType := models.Type().Elem()
	table, err := parseModel(reflect.New(elemType))
	if err != nil {
		return err
	}
	sql := buildSelectSQL(table, where)
	rows, err := db.Query(context.Background(), sql, params...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		elemPtr := reflect.New(elemType)
		elem := elemPtr.Elem()
		table, err = parseModel(elem)
		if err != nil {
			return err
		}
		if err := rows.Scan(pointers(table.columns)...); err != nil {
			return err
		}
		models.Set(reflect.Append(models, elemPtr.Elem()))
	}
	return rows.Err()
}

func buildSelectSQL(table dbTable, where string) string {
	return fmt.Sprintf(
		"SELECT %s FROM %s %s",
		strings.Join(names(table.columns), ","), table.name, where,
	)
}
