package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
)

func TestIsEthereumAddressRequired(t *testing.T) {
	testcases := []struct {
		description string
		address     string
		expectError bool
	}{
		{
			description: "Invalid address - empty",
			address:     "",
			expectError: true,
		},
		{
			description: "Valid address",
			address:     "0x1234567890abcdef1234567890abcdef12345678",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckEthereumAddressRequired(tc.address, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestIsEthereumAddress(t *testing.T) {
	testcases := []struct {
		description string
		address     string
		expectError bool
	}{
		{
			description: "Valid address - empty",
			address:     "",
		},
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
			description: "Invalid address with non-hex characters",
			address:     "0xGHIJKL7890abcdef1234567890abcdef12345678",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckEthereumAddress(tc.address, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestBigIntRequired(t *testing.T) {
	testcases := []struct {
		description string
		address     string
		expectError bool
	}{
		{
			description: "Invalid big int - empty",
			address:     "",
			expectError: true,
		},
		{
			description: "Valid big int",
			address:     "1",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckBigIntRequired(tc.address, "testValue")
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
		description   string
		value         string
		expectedError string
	}{
		{
			description: "Valid big integer - empty",
			value:       "",
		},
		{
			description: "Valid big integer within uint256 range",
			value:       "1234567890",
		},
		{
			description: "Maximum uint256 value",
			value:       "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		},
		{
			description:   "Value exceeding uint256 range",
			value:         "115792089237316195423570985008687907853269984665640564039457584007913129639936",
			expectedError: "too big to fit in uint256",
		},
		{
			description:   "Invalid numeric string",
			value:         "123abc456",
			expectedError: "not a valid value",
		},
		{
			description:   "Negative uint256 value",
			value:         "-1",
			expectedError: "must be a positive value",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckBigInt(tc.value, "testValue")
			if tc.expectedError != "" {
				require.Contains(t, err.Error(), tc.expectedError, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestChainIdRequired(t *testing.T) {
	testcases := []struct {
		description string
		value       int
		expectError bool
	}{
		{
			description: "Invalid chain id - zero",
			value:       0,
			expectError: true,
		},
		{
			description: "Valid chain id - Ethereum",
			value:       chains.Ethereum,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckChainIdRequired(tc.value, "testValue")
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
		description string
		value       int
		expectError bool
	}{
		{
			description: "Valid chain id - zero",
			value:       0,
		},
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
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckChainId(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestPrivateKeyRequired(t *testing.T) {
	testcases := []struct {
		description string
		address     string
		expectError bool
	}{
		{
			description: "Invalid empty private key",
			address:     "",
			expectError: true,
		},
		{
			description: "Valid private key",
			address:     "a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckPrivateKeyRequired(tc.address, "testValue")
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
			description: "Valid empty private key",
			address:     "",
		},
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
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckPrivateKey(tc.address, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestApprovalType(t *testing.T) {
	testcases := []struct {
		description  string
		approvalType int
		expectError  bool
	}{
		{
			description:  "Valid approval type 0 (PermitIfPossible)",
			approvalType: 0,
		},
		{
			description:  "Valid approval type 1 (PermitAlways)",
			approvalType: 1,
		},
		{
			description:  "Valid approval type 2 (ApprovalAlways)",
			approvalType: 2,
		},
		{
			description:  "Invalid approval type",
			approvalType: 3,
			expectError:  true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckApprovalType(tc.approvalType, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestSlippageRequired(t *testing.T) {
	testcases := []struct {
		description string
		value       float32
		expectError bool
	}{
		{
			description: "Invalid slippage value - zero",
			value:       0,
			expectError: true,
		},
		{
			description: "Valid slippage value",
			value:       1,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckSlippageRequired(tc.value, "testValue")
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
			description: "Valid slippage value - zero",
			value:       0,
		},
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
			err := CheckSlippage(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckPage(t *testing.T) {
	testcases := []struct {
		description string
		value       float32
		expectError bool
	}{
		{
			description: "Valid page value - zero",
			value:       0,
		},
		{
			description: "Valid page value - positive integer",
			value:       1,
		},
		{
			description: "Invalid page value - negative integer",
			value:       -1,
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckPage(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckLimit(t *testing.T) {
	testcases := []struct {
		description string
		value       float32
		expectError bool
	}{
		{
			description: "Valid limit value - zero",
			value:       0,
		},
		{
			description: "Valid limit value - positive integer",
			value:       1,
		},
		{
			description: "Invalid limit value - negative integer",
			value:       -1,
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckLimit(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckCheckStatusesInts(t *testing.T) {
	testcases := []struct {
		description string
		values      []float32
		expectError bool
	}{
		{
			description: "Valid status int - nil",
			values:      nil,
		},
		{
			description: "Valid status int - 1",
			values:      []float32{1},
		},
		{
			description: "Valid status int - 2",
			values:      []float32{2},
		},
		{
			description: "Valid status int - 3",
			values:      []float32{3},
		},
		{
			description: "Valid status int combo - 1 and 2",
			values:      []float32{1, 2},
		},
		{
			description: "Invalid status int combo - 1 twice",
			values:      []float32{1, 1},
			expectError: true,
		},
		{
			description: "Invalid status int combo - 0",
			values:      []float32{0},
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckStatusesInts(tc.values, "testValues")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckCheckStatusesStrings(t *testing.T) {
	testcases := []struct {
		description string
		values      []string
		expectError bool
	}{
		{
			description: "Valid status int - nil",
			values:      nil,
		},
		{
			description: "Valid status int - 1",
			values:      []string{"1"},
		},
		{
			description: "Valid status int - 2",
			values:      []string{"2"},
		},
		{
			description: "Valid status int - 3",
			values:      []string{"3"},
		},
		{
			description: "Valid status int combo - 1 and 2",
			values:      []string{"1", "2"},
		},
		{
			description: "Invalid status int combo - 1 twice",
			values:      []string{"1", "1"},
			expectError: true,
		},
		{
			description: "Invalid status int combo - 0",
			values:      []string{"0"},
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckStatusesStrings(tc.values, "testValues")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckSortBy(t *testing.T) {
	testcases := []struct {
		description string
		value       string
		expectError bool
	}{
		{
			description: "Valid sortBy parameter - empty",
			value:       "",
		},
		{
			description: "Valid sortBy parameter - createDateTime",
			value:       "createDateTime",
		},
		{
			description: "Valid sortBy parameter - takerRate",
			value:       "takerRate",
		},
		{
			description: "Valid sortBy parameter - makerRate",
			value:       "makerRate",
		},
		{
			description: "Valid sortBy parameter - makerAmount",
			value:       "makerAmount",
		},
		{
			description: "Valid sortBy parameter - takerAmount",
			value:       "takerAmount",
		},
		{
			description: "Invalid sortBy parameter - random",
			value:       "random",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckSortBy(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckOrderHashRequired(t *testing.T) {
	testcases := []struct {
		description string
		value       string
		expectError bool
	}{
		{
			description: "Invalid sortBy parameter - empty",
			value:       "",
			expectError: true,
		},
		{
			description: "Valid sortBy parameter - createDateTime",
			value:       "0x123",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckOrderHashRequired(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckOrderHash(t *testing.T) {
	testcases := []struct {
		description string
		value       string
		expectError bool
	}{
		{
			description: "Valid parameter - empty",
			value:       "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckOrderHash(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckProtocols(t *testing.T) {
	testcases := []struct {
		description string
		value       string
		expectError bool
	}{
		{
			description: "Valid protocols string - empty",
			value:       "",
		},
		{
			description: "Valid protocols string - one protocol",
			value:       "UNISWAP",
		},
		{
			description: "Valid protocols string - two protocols",
			value:       "UNISWAP,UNISWAP_V2",
		},
		{
			description: "Valid protocols string - many protocols",
			value:       "UNISWAP,UNISWAP_V2,UNISWAP_V3",
		},
		{
			description: "Invalid protocols string - trailing comma",
			value:       "UNISWAP,UNISWAP_V2,",
			expectError: true,
		},
		{
			description: "Invalid protocols string - whitespace",
			value:       "UNISWAP, UNISWAP_V2,",
			expectError: true,
		},
		{
			description: "Invalid protocols string - special characters",
			value:       "UNISWAP,UNISWAP*_V2",
			expectError: true,
		},
		{
			description: "Invalid protocols string - duplicate protocols",
			value:       "UNISWAP,UNISWAP",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckProtocols(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckFee(t *testing.T) {
	testcases := []struct {
		description string
		value       float32
		expectError bool
	}{
		{
			description: "Valid fee - empty/min",
			value:       0,
		},
		{
			description: "Valid fee - max",
			value:       3,
		},
		{
			description: "Invalid fee - negative",
			value:       -0.01,
			expectError: true,
		},
		{
			description: "Invalid fee - too high",
			value:       3.0001,
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckFee(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckFloat32NonNegativeWhole(t *testing.T) {
	testcases := []struct {
		description string
		value       float32
		expectError bool
	}{
		{
			description: "Valid float32 - zero",
			value:       0,
		},
		{
			description: "Valid float32 - positive",
			value:       1,
		},
		{
			description: "Invalid float32 - negative",
			value:       -1,
			expectError: true,
		},
		{
			description: "Invalid float32 - decimal",
			value:       0.5,
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckFloat32NonNegativeWhole(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckConnectorTokens(t *testing.T) {
	testcases := []struct {
		description string
		value       string
		expectError bool
	}{
		{
			description: "Valid connector tokens string - empty",
			value:       "",
		},
		{
			description: "Valid connector tokens string - ethereum address",
			value:       "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		},
		{
			description: "Valid connector tokens string - two ethereum addresses",
			value:       "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48,0x6b175474e89094c44da98b954eedeac495271d0f",
		},
		{
			description: "Valid connector tokens string - many ethereum addresses",
			value:       "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48,0x6b175474e89094c44da98b954eedeac495271d0f",
		},
		{
			description: "Invalid connector tokens string - trailing comma",
			value:       "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48,0x6b175474e89094c44da98b954eedeac495271d0f,",
			expectError: true,
		},
		{
			description: "Invalid connector tokens string - whitespace",
			value:       "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48, 0x6b175474e89094c44da98b954eedeac495271d0f,",
			expectError: true,
		},
		{
			description: "Invalid connector tokens string - invalid characters",
			value:       "0xP0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			expectError: true,
		},
		{
			description: "Invalid connector tokens string - wrong length address",
			value:       "0xEA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			expectError: true,
		},
		{
			description: "Invalid connector tokens string - duplicate addresses",
			value:       "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48,0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckConnectorTokens(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestCheckPermitHash(t *testing.T) {
	testcases := []struct {
		description string
		value       string
		expectError bool
	}{
		{
			description: "Valid permit hash - empty",
			value:       "",
		},
		{
			description: "Valid permit hash",
			value:       "0x00000000000000000000000050c5df26654b5efbdd0c54a062dfa6012933defe0000000000000000000000001111111254eeb25477b68fb85ed929f73a960582ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000000000000000000000000000000000000000000065c6e1b4000000000000000000000000000000000000000000000000000000000000001bc5d700e912fc92bdc59fd1a4963278199a0fc69e95aa8e4b6a3b2bc7387bc2c86137bb710e3dbd3ab0b02762c805902bf9e5f23548ef8431c22dbb7db800d523",
		},
		{
			description: "Invalid permit hash - invalid character",
			value:       "0xT6487c5cb7c4b20202d34117abc57b1c7d91570e100d8a16eced3dbbe8b22eee41339c0772fc6affe3bb8b2a72e9e3e2e2b061d351936c1d534fcfe8d336073d1b",
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckPermitHash(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}

func TestExpireAfter(t *testing.T) {
	testcases := []struct {
		description string
		value       int64
		expectError bool
	}{
		{
			description: "Valid timestamp - empty",
			value:       0,
		},
		{
			description: "Valid timestamp - year 2030",
			value:       1897205247,
		},
		{
			description: "Invalid timestamp - past",
			value:       1707791247,
			expectError: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			err := CheckExpireAfter(tc.value, "testValue")
			if tc.expectError {
				require.Error(t, err, fmt.Sprintf("%s should have caused an error", tc.description))
			} else {
				require.NoError(t, err, fmt.Sprintf("%s should not have caused an error", tc.description))
			}
		})
	}
}
