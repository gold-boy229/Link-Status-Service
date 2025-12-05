package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortStrings(t *testing.T) {
	tests := []struct {
		name           string
		inputData      []string
		expectedResult []string
	}{
		{
			name:           "test a, b, c",
			inputData:      []string{"a", "b", "c"},
			expectedResult: []string{"a", "b", "c"},
		},
		{
			name:           "test a, c, b",
			inputData:      []string{"a", "c", "b"},
			expectedResult: []string{"a", "b", "c"},
		},
		{
			name:           "test b, a, c",
			inputData:      []string{"b", "a", "c"},
			expectedResult: []string{"a", "b", "c"},
		},
		{
			name:           "test b, c, a",
			inputData:      []string{"b", "c", "a"},
			expectedResult: []string{"a", "b", "c"},
		},
		{
			name:           "test c, a, b",
			inputData:      []string{"c", "a", "b"},
			expectedResult: []string{"a", "b", "c"},
		},
		{
			name:           "test c, b, a",
			inputData:      []string{"c", "b", "a"},
			expectedResult: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortStrings(tt.inputData)
			assert.EqualValues(t, tt.expectedResult, result)
		})
	}
}
