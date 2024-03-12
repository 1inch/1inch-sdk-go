package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSubset(t *testing.T) {
	testCases := []struct {
		description string
		sliceA      []int
		sliceB      []int
		expected    bool
	}{
		{
			description: "sliceA is a subset of sliceB",
			sliceA:      []int{1, 2},
			sliceB:      []int{1, 2, 3, 4},
			expected:    true,
		},
		{
			description: "sliceA is not a subset of sliceB",
			sliceA:      []int{1, 2, 5},
			sliceB:      []int{1, 2, 3, 4},
			expected:    false,
		},
		{
			description: "sliceA is empty",
			sliceA:      []int{},
			sliceB:      []int{1, 2, 3, 4},
			expected:    true,
		},
		{
			description: "both slices are empty",
			sliceA:      []int{},
			sliceB:      []int{},
			expected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := IsSubset(tc.sliceA, tc.sliceB)
			assert.Equal(t, tc.expected, result, fmt.Sprintf("%v: expected %v, got %v", tc.description, tc.expected, result))
		})
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		description string
		value       int
		slice       []int
		expected    bool
	}{
		{
			description: "value is present in the slice",
			value:       1,
			slice:       []int{1, 2, 3},
			expected:    true,
		},
		{
			description: "value is not present in the slice",
			value:       4,
			slice:       []int{1, 2, 3},
			expected:    false,
		},
		{
			description: "slice is empty",
			value:       1,
			slice:       []int{},
			expected:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := Contains(tc.value, tc.slice)
			assert.Equal(t, tc.expected, result, fmt.Sprintf("%v: expected %v, got %v", tc.description, tc.expected, result))
		})
	}
}
