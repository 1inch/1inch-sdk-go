package orderbook

//
//import (
//	"bytes"
//	"math/big"
//	"os"
//	"testing"
//	"time"
//
//	"github.com/1inch/1inch-sdk-go/internal/helpers/constants/addresses"
//	"github.com/1inch/1inch-sdk-go/internal/helpers/constants/amounts"
//	"github.com/1inch/1inch-sdk-go/internal/helpers/constants/chains"
//	"github.com/1inch/1inch-sdk-go/internal/helpers/constants/tokens"
//
//	"github.com/ethereum/go-ethereum/ethclient"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//
//	"github.com/1inch/1inch-sdk-go/client/models"
//
//	"github.com/1inch/1inch-sdk-go/helpers"
//)
//
//func TestTrim0x(t *testing.T) {
//	testcases := []struct {
//		description string
//		input       string
//		expected    string
//	}{
//		{
//			description: "String starts with 0x",
//			input:       "0xabcdef",
//			expected:    "abcdef",
//		},
//		{
//			description: "String does not start with 0x",
//			input:       "abcdef",
//			expected:    "abcdef",
//		},
//		{
//			description: "Empty string",
//			input:       "",
//			expected:    "",
//		},
//		{
//			description: "String is just 0x",
//			input:       "0x",
//			expected:    "",
//		},
//		{
//			description: "String starts with 0X (uppercase)",
//			input:       "0Xabcdef",
//			expected:    "0Xabcdef", // note: the function only trims lowercase "0x"
//		},
//	}
//
//	for _, tc := range testcases {
//		t.Run(tc.description, func(t *testing.T) {
//			result := Trim0x(tc.input)
//			require.Equal(t, tc.expected, result)
//		})
//	}
//}
//
//func TestCumulativeSum(t *testing.T) {
//	testcases := []struct {
//		description  string
//		initial      int
//		values       []int
//		expectedSums []int
//	}{
//		{
//			description:  "Initial value is 0",
//			initial:      0,
//			values:       []int{5, 10, 15},
//			expectedSums: []int{5, 15, 30},
//		},
//		{
//			description:  "Initial value is 5",
//			initial:      5,
//			values:       []int{5, 10, 15},
//			expectedSums: []int{10, 20, 35},
//		},
//		{
//			description:  "No values passed",
//			initial:      5,
//			values:       []int{},
//			expectedSums: []int{},
//		},
//		{
//			description:  "Negative values",
//			initial:      -5,
//			values:       []int{5, -10, 15},
//			expectedSums: []int{0, -10, 5},
//		},
//	}
//
//	for _, tc := range testcases {
//		t.Run(tc.description, func(t *testing.T) {
//			sumFunc := CumulativeSum(tc.initial)
//
//			for i, value := range tc.values {
//				result := sumFunc(value)
//				require.Equal(t, tc.expectedSums[i], result)
//			}
//		})
//	}
//}
//
//func TestGenerateSalt(t *testing.T) {
//	t.Run("Generated salt is unique", func(t *testing.T) {
//		salt1 := GenerateSalt()
//		time.Sleep(1 * time.Millisecond) // Sleep for a millisecond to ensure a different time
//		salt2 := GenerateSalt()
//		require.NotEqual(t, salt1, salt2)
//	})
//
//	t.Run("Generated salt has expected length", func(t *testing.T) {
//		salt := GenerateSalt()
//		// Since we're using UnixNano / Millisecond, it should be a long string but not as long as nano time.
//		require.True(t, len(salt) > 5 && len(salt) < 20)
//	})
//}
//
//func TestCreateLimitOrder(t *testing.T) {
//
//	staticSalt := "100000000"
//	mockGenerateSaltFunction := func() string {
//		return staticSalt
//	}
//
//	tests := []struct {
//		name          string
//		orderRequest  models.CreateOrderParams
//		interactions  []string // TODO Revisit this to make it more encapsulated
//		mockBigInt    func(string) (*big.Int, error)
//		expectedOrder *models.Order
//		expectError   bool
//		expectedError string
//	}{
//		{
//			name: "happy path",
//			orderRequest: models.CreateOrderParams{
//				chainId:      chains.Polygon,
//				PrivateKey:   os.Getenv("WALLET_KEY"),
//				MakerAsset:   "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
//				TakerAsset:   "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619",
//				MakingAmount: "1000000",
//				TakingAmount: "1000000000",
//				Maker:        "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
//				Taker:        "0x0000000000000000000000000000000000000000",
//			},
//			interactions: []string{"0x", "0x", "0x", "0x", "0xbf15fcd8000000000000000000000000a5eb255ef45dfb48b5d133d08833def69871691d000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000242cc2878d0071150dff0000000000000050c5df26654b5efbdd0c54a062dfa6012933defe00000000000000000000000000000000000000000000000000000000", "0x", "0x", "0x"},
//			expectedOrder: &models.Order{
//				OrderHash: "0xdc9344cfa6d3b4da5a2ad3283e02826d3f569b4472443390d3e1cfe86cacd13f",
//				Signature: "0x317ed3e021851542deeafb4897ef091b010317772b7299477121d0f46cdd32cf1403429b13d2337b459c7a982ac71144ceaad88dd08d5b7c7b8abbe1618070ab1b",
//				Data: models.OrderData{
//					MakerAsset:    "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
//					TakerAsset:    "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619",
//					MakingAmount:  "1000000",
//					TakingAmount:  "1000000000",
//					Salt:          "100000000",
//					Maker:         "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
//					AllowedSender: "0x0000000000000000000000000000000000000000",
//					Receiver:      "0x0000000000000000000000000000000000000000",
//					Offsets:       "4421431254442149611168492388118363282642987198110904030635476664713216",
//					Interactions:  "0xbf15fcd8000000000000000000000000a5eb255ef45dfb48b5d133d08833def69871691d000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000242cc2878d0071150dff0000000000000050c5df26654b5efbdd0c54a062dfa6012933defe00000000000000000000000000000000000000000000000000000000",
//				},
//			},
//			expectError: false,
//		},
//		{
//			name: "empty maker asset",
//			orderRequest: models.CreateOrderParams{
//				chainId:      chains.Polygon,
//				PrivateKey:   os.Getenv("WALLET_KEY"),
//				MakerAsset:   "",
//				TakerAsset:   "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619",
//				MakingAmount: "1000000",
//				TakingAmount: "1000000000",
//				Maker:        "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
//				Taker:        "0x0000000000000000000000000000000000000000",
//			},
//			interactions:  []string{"0x", "0x", "0x", "0x", "0xbf15fcd8000000000000000000000000a5eb255ef45dfb48b5d133d08833def69871691d000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000242cc2878d0071150dff0000000000000050c5df26654b5efbdd0c54a062dfa6012933defe00000000000000000000000000000000000000000000000000000000", "0x", "0x", "0x"},
//			expectError:   true,
//			expectedError: "error hashing typed data",
//		},
//		{
//			name: "empty taker asset",
//			orderRequest: models.CreateOrderParams{
//				chainId:      chains.Polygon,
//				PrivateKey:   os.Getenv("WALLET_KEY"),
//				MakerAsset:   "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
//				TakerAsset:   "",
//				MakingAmount: "1000000",
//				TakingAmount: "1000000000",
//				Maker:        "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
//				Taker:        "0x0000000000000000000000000000000000000000",
//			},
//			interactions:  []string{"0x", "0x", "0x", "0x", "0xbf15fcd8000000000000000000000000a5eb255ef45dfb48b5d133d08833def69871691d000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000242cc2878d0071150dff0000000000000050c5df26654b5efbdd0c54a062dfa6012933defe00000000000000000000000000000000000000000000000000000000", "0x", "0x", "0x"},
//			expectError:   true,
//			expectedError: "error hashing typed data",
//		},
//		{
//			name: "invalid private key",
//			orderRequest: models.CreateOrderParams{
//				chainId:      chains.Polygon,
//				PrivateKey:   "invalid_private_key", // non-hex or short length key
//				MakerAsset:   "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
//				TakerAsset:   "0x7ceb23fd6bc0add59e62ac25578270cff1b9f619",
//				MakingAmount: "1000000",
//				TakingAmount: "1000000000",
//				Maker:        "0x2c9b2DBdbA8A9c969Ac24153f5C1c23CB0e63914",
//				Taker:        "0x0000000000000000000000000000000000000000",
//			},
//			interactions:  []string{"0x", "0x", "0x", "0x", "0xbf15fcd8000000000000000000000000a5eb255ef45dfb48b5d133d08833def69871691d000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000242cc2878d0071150dff0000000000000050c5df26654b5efbdd0c54a062dfa6012933defe00000000000000000000000000000000000000000000000000000000", "0x", "0x", "0x"},
//			expectError:   true,
//			expectedError: "error converting private key to ECDSA",
//		},
//	}
//
//	// Save the original salt generation function and return it later
//	originalGenerateSalt := GenerateSalt
//	GenerateSalt = mockGenerateSaltFunction
//	defer func() {
//		GenerateSalt = originalGenerateSalt
//	}()
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			result, err := CreateLimitOrderMessage(tc.orderRequest, tc.interactions)
//
//			if tc.expectError {
//				assert.Error(t, err)
//				if tc.expectedError != "" {
//					assert.Contains(t, err.Error(), tc.expectedError)
//				}
//			} else {
//				require.NoError(t, err)
//				// Validate all the fields in the order data to be as expected
//				assert.Equal(t, tc.expectedOrder.OrderHash, result.OrderHash, "Order hash does not match expected value")
//				assert.Equal(t, tc.expectedOrder.Signature, result.Signature, "Signature does not match expected value")
//				// Compare the data fields individually or as a whole
//				assert.Equal(t, tc.expectedOrder.Data, result.Data, "Order data does not match expected value")
//			}
//		})
//	}
//}
//
//func TestConfirmTradeWithUser(t *testing.T) {
//
//	order := &models.Order{
//		Data: models.OrderData{
//			MakerAsset:   tokens.EthereumUsdc,
//			TakerAsset:   tokens.EthereumDai,
//			MakingAmount: amounts.Ten6 + "1",
//			TakingAmount: amounts.Ten18,
//			Maker:        addresses.Vitalik,
//		},
//	}
//
//	tests := []struct {
//		name           string
//		userInput      string
//		expectedResult bool
//		expectedOutput string
//	}{
//		{
//			name:           "User inputs 'y'",
//			userInput:      "y\n",
//			expectedResult: true,
//		},
//		{
//			name:           "User inputs 'Y'",
//			userInput:      "Y\n",
//			expectedResult: true,
//		},
//		{
//			name:           "User inputs 'n'",
//			userInput:      "n\n",
//			expectedResult: false,
//		},
//		{
//			name:           "User inputs nothing",
//			userInput:      "\n",
//			expectedResult: false,
//		},
//		{
//			name:           "User inputs other text",
//			userInput:      "other\n",
//			expectedResult: false,
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			ethClient, err := ethclient.Dial(os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"))
//			require.NoError(t, err)
//			reader := bytes.NewBufferString(tc.userInput)
//			writer := helpers.NoOpPrinter{}
//			result, err := confirmLimitOrderWithUser(order, ethClient, reader, writer)
//
//			assert.NoError(t, err)
//			assert.Equal(t, tc.expectedResult, result)
//		})
//	}
//}
//
//func TestConcatenateInteractions(t *testing.T) {
//
//	tests := []struct {
//		name           string
//		interactions   []string
//		expectedResult string
//	}{
//		{
//			name:           "Empty slice",
//			interactions:   []string{},
//			expectedResult: "0x",
//		},
//		{
//			name:           "Single element without prefix",
//			interactions:   []string{"abcdef"},
//			expectedResult: "0xabcdef",
//		},
//		{
//			name:           "Single element with prefix",
//			interactions:   []string{"0x123456"},
//			expectedResult: "0x123456",
//		},
//		{
//			name:           "Multiple elements mixed prefixes",
//			interactions:   []string{"0xabcdef", "123456", "0x7890"},
//			expectedResult: "0xabcdef1234567890",
//		},
//		{
//			name:           "Multiple elements all with prefix",
//			interactions:   []string{"0xabcdef", "0x123456", "0x7890"},
//			expectedResult: "0xabcdef1234567890",
//		},
//		{
//			name:           "Multiple elements none with prefix",
//			interactions:   []string{"abcdef", "123456", "7890"},
//			expectedResult: "0xabcdef1234567890",
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			result := concatenateInteractions(tc.interactions)
//			assert.Equal(t, tc.expectedResult, result)
//		})
//	}
//}
