package validate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckStringRequired(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:  "Non-empty string",
			input: "hello",
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckStringRequired(tc.input, "testValue")
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCheckString(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Non-empty string",
			input: "hello",
		},
		{
			name:  "Empty string",
			input: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckString(tc.input, "testValue")
			require.NoError(t, err)
		})
	}
}

func TestCheckBoolean(t *testing.T) {
	tests := []struct {
		name  string
		input bool
	}{
		{
			name:  "True",
			input: true,
		},
		{
			name:  "False",
			input: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := CheckBoolean(tc.input, "testValue")
			require.NoError(t, err)
		})
	}
}
