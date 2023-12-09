package qry

import (
	"reflect"
	"testing"
)

func Test_parseColumn(t *testing.T) {
	t.Run("returns empty name and false when tag is empty", func(t *testing.T) {
		tag := reflect.StructTag("")
		name, gen := parseColumn(tag)
		if name != "" || gen {
			t.Error("Expected empty name and false")
		}
	})

	t.Run("returns column name and false when tag has no generated", func(t *testing.T) {
		tag := reflect.StructTag(`column:"col1"`)
		name, gen := parseColumn(tag)
		if name != "col1" || gen {
			t.Error("Expected column name and false")
		}
	})

	t.Run("returns column name and true when tag has generated", func(t *testing.T) {
		tag := reflect.StructTag(`column:"col1,generated"`)
		name, gen := parseColumn(tag)
		if name != "col1" || !gen {
			t.Error("Expected column name and true")
		}
	})
}
