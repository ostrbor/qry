package qry

import (
	"reflect"
	"testing"
)

func Test_sqlValues(t *testing.T) {
	p := sqlValues(2, 2)
	if p != "($1,$2),($3,$4)" {
		t.Error("failed 2 times row with 2 vals")
	}

	p = sqlValues(1, 3)
	if p != "($1),($2),($3)" {
		t.Error("failed 3 times row with 1 val")
	}

	p = sqlValues(2, 0)
	if p != "($1,$2)" {
		t.Error("failed 0 times row with 2 vals")
	}
}

func Test_batch(t *testing.T) {
	sliceType := reflect.TypeOf([]int{})

	s := reflect.MakeSlice(sliceType, 10, 10)
	if len(batch(sliceItems(s), 1)) != 10 {
		t.Error("must be 10 batch, each with 1 item")
	}

	s = reflect.MakeSlice(sliceType, 10, 10)
	if len(batch(sliceItems(s), 2)) != 5 {
		t.Error("must be 5 batch, each with 2 items")
	}

	s = reflect.MakeSlice(sliceType, 10, 10)
	if len(batch(sliceItems(s), 10)) != 1 {
		t.Error("must be 1 group with 10 items")
	}

	s = reflect.MakeSlice(sliceType, 10, 10)
	if len(batch(sliceItems(s), 0)) != 1 {
		t.Error("must be 1 group with 10 items")
	}

	s = reflect.MakeSlice(sliceType, 10, 10)
	if len(batch(sliceItems(s), 100)) != 1 {
		t.Error("must be 1 group with 10 items")
	}

	s = reflect.MakeSlice(sliceType, 3, 3)
	b := batch(sliceItems(s), 2)
	if len(b) != 2 {
		t.Error("must be 2 batch")
	} else if len(b[1]) != 1 {
		t.Error("last group must have 1 item")
	}
}
