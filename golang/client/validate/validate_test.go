package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/svanas/1inch-sdk/golang/helpers/consts/chains"
)

func TestIsEthereumAddress(t *testing.T) {
	testcases := []struct {
		description string
		address     string
		expectError bool
	}{
		{
			description: "Valid address with lowercase letters",
			address:     "0x1234567890abcdef1234567890abcdef12345678",
		},
		{
			description: "Valid address with mixed case letters",
			address:     "0x1234567890ABCDEF1234567890abcdef12345678",
		},
		{
			description: "Invalid address without 0x prefix",
			address:     "1234567890abcdef1234567890abcdef12345678",
			expectError: true,
		},
		{
			description: "Invalid address too short",
			address:     "0x12345",
			expectError: true,
		},
		{
			description: "Invalid address too long",
			address:     "0x1234567890abcdef1234567890abcdef1234567890",
			expectError: true,
		},
		{
			description: "Invalid empty address",
			address:     "",
			expectError: true,
		},
		{
			description: "Invalid address with non-hex characters",
			address:     "0xGHIJKL7890abcdef1234567890abcdef12345678",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckEthereumAddress(tc.address, "")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestBigInt(t *testing.T) {
	testcases := []struct {
		description string
		value       string
		expectError bool
	}{
		{
			description: "Valid big integer within uint256 range",
			value:       "1234567890",
		},
		{
			description: "Value exceeding uint256 range",
			value:       "115792089237316195423570985008687907853269984665640564039457584007913129639936",
			expectError: true,
		},
		{
			description: "Empty value",
			value:       "",
			expectError: true,
		},
		{
			description: "Invalid numeric string",
			value:       "123abc456",
			expectError: true,
		},
		{
			description: "Maximum uint256 value",
			value:       "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckBigInt(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestChainId(t *testing.T) {
	testcases := []struct {
		description  string
		value        int
		variableName string
		expectError  bool
	}{
		{
			description: "Valid chain id - Ethereum",
			value:       chains.Ethereum,
		},
		{
			description: "Valid chain id - Polygon",
			value:       chains.Polygon,
		},
		{
			description: "Invalid chain id",
			value:       999999,
			expectError: true,
		},
		{
			description: "Chain id is required",
			value:       0,
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckChainId(tc.value, "testChainId")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestPrivateKey(t *testing.T) {
	testcases := []struct {
		description string
		address     string
		expectError bool
	}{
		{
			description: "Valid private key",
			address:     "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1",
		},
		{
			description: "Invalid private key with special characters",
			address:     "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b!",
			expectError: true,
		},
		{
			description: "Invalid private key with short length",
			address:     "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5",
			expectError: true,
		},
		{
			description: "Invalid private key with long length",
			address:     "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1a2b3",
			expectError: true,
		},
		{
			description: "Empty private key",
			address:     "",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckPrivateKey(tc.address, "testPrivateKey")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestSlippage(t *testing.T) {
	testcases := []struct {
		description string
		value       float32
		expectError bool
	}{
		{
			description: "Valid slippage value - lower boundary",
			value:       0.01,
		},
		{
			description: "Valid slippage value - upper boundary",
			value:       50,
		},
		{
			description: "Valid slippage value - mid range",
			value:       25,
		},
		{
			description: "Invalid slippage value - empty",
			value:       0,
			expectError: true,
		},
		{
			description: "Invalid slippage value - below lower boundary",
			value:       -1,
			expectError: true,
		},
		{
			description: "Invalid slippage value - above upper boundary",
			value:       51,
			expectError: true,
		},
		{
			description: "Valid slippage value with decimals",
			value:       12.34,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckSlippage(tc.value, "testSlippage")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}
