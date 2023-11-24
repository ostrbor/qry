package qry

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"reflect"
	"strings"
)

var errNil = errors.New("nil is not allowed")

func errKind(a, b reflect.Kind) error {
	return fmt.Errorf("unexpected kind: %s != %s", a, b)
}

// dereference returns underlying pointer value.
func dereference(v interface{}) (ptrVal reflect.Value, err error) {
	if v == nil {
		return ptrVal, errNil
	}
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		return ptrVal, errKind(val.Kind(), reflect.Ptr)
	}
	return val.Elem(), nil
}

// validateStructType checks if the given type represents a struct.
func validateStructType(t reflect.Type) (err error) {
	if t.Kind() != reflect.Struct {
		return errKind(t.Kind(), reflect.Struct)
	}
	return nil
}

// sliceItems extracts individual items from a reflect.Value representing a slice.
// todo return err?
func sliceItems(slice reflect.Value) (items []reflect.Value) {
	// todo remove check?
	if slice.Kind() != reflect.Slice {
		return nil
	}
	items = make([]reflect.Value, 0, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		items = append(items, slice.Index(i))
	}
	return items
}

// sqlValues generates $ placeholders for the VALUES clause in SQL.
// The placeholders are arranged in rows, and the number of rows is determined by totalRows.
// Each row contains phsPerRow placeholders.
// Example: ($1,$2),($3,$4) for phsPerRow=2 and totalRows=2.
func sqlValues(phsPerRow, totalRows int) string {
	// todo return err if phsPerRow < 1?
	// todo return err if totalRows < 1?
	if totalRows < 1 {
		totalRows = 1
	}

	rows := make([]string, 0, totalRows)
	// phs is used to collect placeholders for each row before joining them.
	phs := make([]string, 0, phsPerRow)
	for i := 1; i <= phsPerRow*totalRows; i++ {
		ph := fmt.Sprintf("$%d", i)
		phs = append(phs, ph)
		if len(phs) == phsPerRow {
			row := fmt.Sprintf("(%s)", strings.Join(phs, ","))
			rows = append(rows, row)
			phs = phs[:0]
		}
	}
	return strings.Join(rows, ",")
}

type Querier interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}
