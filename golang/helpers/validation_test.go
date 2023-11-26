package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEthereumAddress(t *testing.T) {
	testcases := []struct {
		description string
		address     string
		expecting   bool
	}{
		{
			description: "Valid address with lowercase letters",
			address:     "0x1234567890abcdef1234567890abcdef12345678",
			expecting:   true,
		},
		{
			description: "Valid address with mixed case letters",
			address:     "0x1234567890ABCDEF1234567890abcdef12345678",
			expecting:   true,
		},
		{
			description: "Invalid address without 0x prefix",
			address:     "1234567890abcdef1234567890abcdef12345678",
			expecting:   false,
		},
		{
			description: "Invalid address too short",
			address:     "0x12345",
			expecting:   false,
		},
		{
			description: "Invalid address too long",
			address:     "0x1234567890abcdef1234567890abcdef1234567890",
			expecting:   false,
		},
		{
			description: "Invalid empty address",
			address:     "",
			expecting:   false,
		},
		{
			description: "Invalid address with non-hex characters",
			address:     "0xGHIJKL7890abcdef1234567890abcdef12345678",
			expecting:   false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			result := IsEthereumAddress(tc.address)
			assert.Equal(t, tc.expecting, result, tc.description)
		})
	}
}
