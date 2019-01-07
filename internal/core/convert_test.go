package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToStringSlice(t *testing.T) {
	tests := []struct {
		scenario string
		in       interface{}
		expected []string
	}{
		{
			scenario: "convert a slice of strings returns a slice of strings",
			in:       []string{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			scenario: "convert a slice containing strings returns a slice of strings",
			in:       []interface{}{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			scenario: "convert a non slice value returns nil",
			in:       "hello world",
			expected: nil,
		},
		{
			scenario: "convert a slice of values returns a slice of strings",
			in:       []interface{}{"hello", 42},
			expected: []string{"hello", "42"},
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			result := ConvertToStringSlice(test.in)
			assert.Equal(t, test.expected, result)
		})
	}
}
