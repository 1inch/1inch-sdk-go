package orderbook

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"1inch-sdk-golang/helpers/consts/chains"
)

func TestTrim0x(t *testing.T) {
	testcases := []struct {
		description string
		input       string
		expected    string
	}{
		{
			description: "String starts with 0x",
			input:       "0xabcdef",
			expected:    "abcdef",
		},
		{
			description: "String does not start with 0x",
			input:       "abcdef",
			expected:    "abcdef",
		},
		{
			description: "Empty string",
			input:       "",
			expected:    "",
		},
		{
			description: "String is just 0x",
			input:       "0x",
			expected:    "",
		},
		{
			description: "String starts with 0X (uppercase)",
			input:       "0Xabcdef",
			expected:    "0Xabcdef", // note: the function only trims lowercase "0x"
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			result := Trim0x(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestCumulativeSum(t *testing.T) {
	testcases := []struct {
		description  string
		initial      int
		values       []int
		expectedSums []int
	}{
		{
			description:  "Initial value is 0",
			initial:      0,
			values:       []int{5, 10, 15},
			expectedSums: []int{5, 15, 30},
		},
		{
			description:  "Initial value is 5",
			initial:      5,
			values:       []int{5, 10, 15},
			expectedSums: []int{10, 20, 35},
		},
		{
			description:  "No values passed",
			initial:      5,
			values:       []int{},
			expectedSums: []int{},
		},
		{
			description:  "Negative values",
			initial:      -5,
			values:       []int{5, -10, 15},
			expectedSums: []int{0, -10, 5},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			sumFunc := CumulativeSum(tc.initial)

			for i, value := range tc.values {
				result := sumFunc(value)
				require.Equal(t, tc.expectedSums[i], result)
			}
		})
	}
}

func TestGenerateSalt(t *testing.T) {
	t.Run("Generated salt is unique", func(t *testing.T) {
		salt1 := GenerateSalt()
		time.Sleep(1 * time.Millisecond) // Sleep for a millisecond to ensure a different time
		salt2 := GenerateSalt()
		require.NotEqual(t, salt1, salt2)
	})

	t.Run("Generated salt has expected length", func(t *testing.T) {
		salt := GenerateSalt()
		// Since we're using UnixNano / Millisecond, it should be a long string but not as long as nano time.
		require.True(t, len(salt) > 5 && len(salt) < 20)
	})
}

func TestCreateLimitOrder(t *testing.T) {

	staticSalt := "100000000"
	mockGenerateSaltFunction := func() string {
		return staticSalt
	}

	tests := []struct {
		name          string
		orderRequest  OrderRequest
		chainId       int
		key           string
		mockBigInt    func(string) (*big.Int, error)
		expectedOrder *Order
		expectError   bool
		expectedError string
	}{
		{
			name: "happy path",
			orderRequest: OrderRequest{
				FromToken:    "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
				ToToken:      "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619",
				MakingAmount: 1000000,
				TakingAmount: 1000000000,
				SourceWallet: "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
				//AllowedSender: "0x0000000000000000000000000000000000000000",
				Receiver: "0x0000000000000000000000000000000000000000",
				//Offsets:       "0",
				//Interactions:  "0x",
			},
			chainId: chains.Polygon,
			key:     os.Getenv("WALLET_KEY"),
			expectedOrder: &Order{
				OrderHash: "0xaba6be89f39d5c6fce46648caaa974a0bc31a842b157b29a99f05d4c2fa7b781",
				Signature: "0xfc1e704bcd719c076396f2e4795501aec5f50607ef8ec123d0f9306f9420b8543bcee9e934a86c8ecf1a57723463108568743b4c8da03f699a1127a46112e8371c",
				Data: OrderData{
					MakerAsset:    "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
					TakerAsset:    "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619",
					MakingAmount:  "1000000",
					TakingAmount:  "1000000000",
					Salt:          "100000000",
					Maker:         "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
					AllowedSender: "0x0000000000000000000000000000000000000000",
					Receiver:      "0x0000000000000000000000000000000000000000",
					Offsets:       "0",
					Interactions:  "0x",
				},
			},
			expectError: false,
		},
		{
			name: "empty maker asset",
			orderRequest: OrderRequest{
				FromToken:    "",
				ToToken:      "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619",
				MakingAmount: 1000000,
				TakingAmount: 1000000000,
				SourceWallet: "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
				//AllowedSender: "0x0000000000000000000000000000000000000000",
				Receiver: "0x0000000000000000000000000000000000000000",
				//Offsets:       "0",
				//Interactions:  "0x",
			},
			chainId:       chains.Polygon,
			key:           os.Getenv("WALLET_KEY"),
			expectError:   true,
			expectedError: "error hashing typed data",
		},
		{
			name: "empty taker asset",
			orderRequest: OrderRequest{
				FromToken:    "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
				ToToken:      "",
				MakingAmount: 1000000,
				TakingAmount: 1000000000,
				SourceWallet: "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
				//AllowedSender: "0x0000000000000000000000000000000000000000",
				Receiver: "0x0000000000000000000000000000000000000000",
				//Offsets:       "0",
				//Interactions:  "0x",
			},
			chainId:       chains.Polygon,
			key:           os.Getenv("WALLET_KEY"),
			expectError:   true,
			expectedError: "error hashing typed data",
		},
		{
			name: "invalid private key",
			orderRequest: OrderRequest{
				FromToken:    "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
				ToToken:      "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619",
				MakingAmount: 1000000,
				TakingAmount: 1000000000,
				SourceWallet: "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
				//AllowedSender: "0x0000000000000000000000000000000000000000",
				Receiver: "0x0000000000000000000000000000000000000000",
				//Offsets:       "0",
				//Interactions:  "0x",
			},
			chainId:       chains.Polygon,
			key:           "invalid_private_key", // non-hex or short length key
			expectError:   true,
			expectedError: "error converting private key to ECDSA",
		},
	}

	// Save the original salt generation function and return it later
	originalGenerateSalt := GenerateSalt
	GenerateSalt = mockGenerateSaltFunction
	defer func() {
		GenerateSalt = originalGenerateSalt
	}()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateLimitOrder(tc.orderRequest, tc.chainId, tc.key)

			if tc.expectError {
				assert.Error(t, err)
				if tc.expectedError != "" {
					assert.Contains(t, err.Error(), tc.expectedError)
				}
			} else {
				require.NoError(t, err)
				// Validate all the fields in the order data to be as expected
				assert.Equal(t, tc.expectedOrder.OrderHash, result.OrderHash, "Order hash does not match expected value")
				assert.Equal(t, tc.expectedOrder.Signature, result.Signature, "Signature does not match expected value")
				// Compare the data fields individually or as a whole
				assert.Equal(t, tc.expectedOrder.Data, result.Data, "Order data does not match expected value")
			}
		})
	}
}
