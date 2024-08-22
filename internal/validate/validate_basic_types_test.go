package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckString(t *testing.T) {
	testCases := []struct {
		description string
		input       interface{}
		expectError bool
	}{
		{
			description: "True",
			input:       "a",
		},
		{
			description: "True 2",
			input:       "12",
		},
		{
			description: "Must fail",
			input:       nil,
			expectError: true,
		},
		{
			description: "Must fail 2",
			input: struct {
				A string
			}{
				A: "a",
			},
			expectError: true,
		},
		{
			description: "Must fail 3",
			input:       []string{"1", "2"},
			expectError: true,
		},
		{
			description: "Must fail 4",
			input:       12,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckString(tc.input, "testValue")
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckBoolean(t *testing.T) {
	testCases := []struct {
		description string
		input       interface{}
		expectError bool
	}{
		{
			description: "True",
			input:       true,
		},
		{
			description: "False",
			input:       false,
		},
		{
			description: "Must fail",
			input:       nil,
			expectError: true,
		},
		{
			description: "Must fail 2",
			input: struct {
				A string
			}{
				A: "a",
			},
			expectError: true,
		},
		{
			description: "Must fail 3",
			input:       []string{"1", "2"},
			expectError: true,
		},
		{
			description: "Must fail 4",
			input:       12,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckBoolean(tc.input, "testValue")
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
