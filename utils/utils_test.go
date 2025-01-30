package utils

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Capitalized",
			input:    "Hello",
			expected: []string{"hello"},
		},
		{
			name:     "Trailing capital",
			input:    "HellO",
			expected: []string{"hello"},
		},
		{
			name:     "Whitespace tail",
			input:    "Hello ",
			expected: []string{"hello"},
		},
		{
			name:     "Whitespace around",
			input:    "   hello   end  ",
			expected: []string{"hello", "end"},
		},
	}

	for _, test := range tests {
		res := CleanInput(test.input)
		if !reflect.DeepEqual(res, test.expected) {
			t.Errorf("failed test: %s", test.name)
		}
	}
}
