package xplorersbot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArrayContainsItem(t *testing.T) {
	assert := require.New(t)
	tests := []struct {
		name     string
		array    []string
		item     string
		expected bool
	}{
		{
			name:     "array contains item - should return true",
			array:    []string{"kubernetes", "aws", "celebrate"},
			item:     "kubernetes",
			expected: true,
		},
		{
			name:     "array does not contain item - should return false",
			array:    []string{"kubernetes", "aws", "celebrate"},
			item:     "someword",
			expected: false,
		},
	}
	for _, test_case := range tests {
		t.Run(test_case.name, func(t *testing.T) {
			result := ArrayContainsItem(test_case.array, test_case.item)
			assert.Equal(test_case.expected, result)
		})
	}

}
