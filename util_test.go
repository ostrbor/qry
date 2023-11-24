package qry

import "testing"

func TestSQLValues(t *testing.T) {
	tests := []struct {
		name          string
		phsPerRow     int
		totalRows     int
		expectedValue string
	}{
		{"SingleRowSinglePlaceholder", 1, 1, "($1)"},
		{"SingleRowMultiplePlaceholders", 3, 1, "($1,$2,$3)"},
		{"MultipleRowsSinglePlaceholder", 1, 3, "($1),($2),($3)"},
		{"MultipleRowsMultiplePlaceholders", 2, 2, "($1,$2),($3,$4)"},
		{"ZeroPlaceholders", 0, 3, ""},
		{"ZeroRows", 2, 0, "($1,$2)"},
		{"ZeroPlaceholdersZeroRows", 0, 0, ""},
		{"OneRow", 3, 1, "($1,$2,$3)"},
		{"OnePlaceholder", 1, 5, "($1),($2),($3),($4),($5)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sqlValues(tt.phsPerRow, tt.totalRows)
			if result != tt.expectedValue {
				t.Errorf("expected: %s, got: %s", tt.expectedValue, result)
			}
		})
	}
}
