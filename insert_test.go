package qry

import "testing"

func Test_splitCols(t *testing.T) {
	t.Run("returns empty slice and empty column when input is empty", func(t *testing.T) {
		var cols []column
		res, gen := splitCols(cols)
		if len(res) != 0 || (column{} != gen) {
			t.Error("Expected empty slice and empty column")
		}
	})

	t.Run("returns slice with non-generated columns and generated column when input has generated column", func(t *testing.T) {
		cols := []column{
			{name: "col1", generated: false},
			{name: "col2", generated: true},
			{name: "col3", generated: false},
		}
		res, gen := splitCols(cols)
		if len(res) != 2 || gen.name != "col2" {
			t.Error("Expected slice with non-generated columns and generated column")
		}
	})

	t.Run("returns slice with all columns and empty column when input has no generated column", func(t *testing.T) {
		cols := []column{
			{name: "col1", generated: false},
			{name: "col2", generated: false},
			{name: "col3", generated: false},
		}
		res, gen := splitCols(cols)
		if len(res) != 3 || (column{} != gen) {
			t.Error("Expected slice with all columns and empty column")
		}
	})
}
