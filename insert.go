package qry

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func MustInsert(db Querier, ptr interface{}) {
	if err := Insert(db, ptr); err != nil {
		panic(err)
	}
}

func Insert(db Querier, ptr interface{}) (err error) {
	v, err := dereference(ptr)
	if err != nil {
		return err
	}
	switch v.Kind() {
	case reflect.Struct:
		return insertOne(db, v)
	case reflect.Slice:
		if err := validateStructType(v.Type().Elem()); err != nil {
			return err
		}
		if v.Len() == 0 {
			return fmt.Errorf("slice len: %d == 0", v.Len())
		}
		return insertMany(db, sliceItems(v))
	}
	return fmt.Errorf("invalid kind: %s", v.Kind())
}

func insertOne(db Querier, model reflect.Value) (err error) {
	table, err := parseModel(model)
	if err != nil {
		return err
	}
	insertCols, scanCol := splitCols(table.columns)
	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		table.name, strings.Join(names(insertCols), ","), sqlValues(len(insertCols), 0),
	)

	if scanCol.name != "" {
		sql += fmt.Sprintf(" RETURNING %s", scanCol.name)
		return db.QueryRow(context.Background(), sql, pointers(insertCols)...).Scan(scanCol.structFieldPtr)
	}
	_, err = db.Exec(context.Background(), sql, values(insertCols)...)
	return err
}

func insertMany(db Querier, models []reflect.Value) (err error) {
	table, err := parseModel(reflect.New(models[0].Type()))
	if err != nil {
		return err
	}
	insertCols, _ := splitCols(table.columns)
	sql := fmt.Sprintf(
		"INSERT INTO %s (%s)",
		table.name, strings.Join(names(insertCols), ","),
	)

	var vals []interface{}
	var gens []column
	for _, m := range models {
		tab, err := parseModel(m)
		if err != nil {
			return err
		}
		insert, scan := splitCols(tab.columns)
		vals = append(vals, values(insert)...)
		gens = append(gens, scan)
	}

	sql += fmt.Sprintf(" VALUES %s", sqlValues(len(insertCols), len(models)))
	if len(gens) == 0 {
		_, err = db.Exec(context.Background(), sql, vals...)
		return err
	}

	sql += fmt.Sprintf(" RETURNING %s", gens[0].name)
	rows, err := db.Query(context.Background(), sql, vals...)
	if err != nil {
		return err
	}

	ptrs := pointers(gens)
	for i := 0; rows.Next(); i++ {
		if i+1 > len(ptrs) {
			rows.Close()
			return fmt.Errorf("%d rows > %d serials", i+1, len(ptrs))
		}
		if err := rows.Scan(ptrs[i]); err != nil {
			rows.Close()
			return err
		}
	}
	return rows.Err()
}

func splitCols(cols []column) (res []column, generated column) {
	res = make([]column, 0, len(cols))
	for _, col := range cols {
		if col.generated {
			generated = col
			continue
		}
		res = append(res, col)
	}
	return
}
