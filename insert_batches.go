package qry

import (
	"fmt"
	"reflect"
)

func InsertInBatches(db Querier, ptr interface{}, batchSize int) (err error) {
	v, err := dereference(ptr)
	if err != nil {
		return err
	}
	t := v.Type()
	if t.Kind() != reflect.Slice {
		return errKind(t.Kind(), reflect.Slice)
	}
	if err := validateStructType(t.Elem()); err != nil {
		return err
	}
	if v.Len() == 0 {
		return fmt.Errorf("slice len: %d == 0", v.Len())
	}
	for _, b := range batch(sliceItems(v), batchSize) {
		err := insertMany(db, b)
		if err != nil {
			return err
		}
	}
	return nil
}

// batch groups a slice of items into batches of a specified size.
// It takes a slice of items and an integer size,
// returning a 2D slice where each inner slice contains at most 'size' items.
// If the input slice is empty, it returns nil.
// Example:
//
//	input: [1, 2, 3, 4, 5], size: 2
//	output: [[1, 2], [3, 4], [5]]
func batch(items []reflect.Value, size int) [][]reflect.Value {
	length := len(items)
	if length == 0 {
		return nil
	}
	if size <= 0 || size > length {
		size = length
	}

	estimatedCapacity := (length + size - 1) / size
	batches := make([][]reflect.Value, 0, estimatedCapacity)
	b := make([]reflect.Value, 0, size)
	for _, item := range items {
		if len(b) == size {
			batches = append(batches, b)
			b = make([]reflect.Value, 0, size)
		}
		b = append(b, item)
	}

	if len(b) > 0 {
		batches = append(batches, b)
	}

	return batches
}
