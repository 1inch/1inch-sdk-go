package fusion

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
)

func TestGenerateSalt(t *testing.T) {
	// Save the original function
	originalBigIntMaxFunc := random_number_generation.BigIntMaxFunc

	// Monkey patch the function
	random_number_generation.BigIntMaxFunc = func(max *big.Int) (*big.Int, error) {
		return big.NewInt(123456), nil
	}

	// Restore the original function after the test
	defer func() {
		random_number_generation.BigIntMaxFunc = originalBigIntMaxFunc
	}()

	tests := []struct {
		name      string
		extension *Extension
		expected  string
		expectErr bool
	}{
		{
			name: "Generate salt when extension is not empty",
			extension: &Extension{
				MakerAssetSuffix: "suffix1",
				TakerAssetSuffix: "suffix2",
				MakingAmountData: "data1",
				TakingAmountData: "data2",
				Predicate:        "predicate",
				MakerPermit:      "permit",
				PreInteraction:   "pre",
				PostInteraction:  "post",
				CustomData:       "custom",
			},
			expected:  "180431178743033967347942937469468920088249224033532329",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expected, err := BigIntFromString(tc.expected)
			require.NoError(t, err)

			result, err := tc.extension.GenerateSalt()
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expected, result)
			}
		})
	}
}

func TestNewExtension(t *testing.T) {
	tests := []struct {
		name              string
		params            ExtensionParams
		expectedExtension *Extension
		expectErr         bool
		errMsg            string
	}{
		{
			name: "Valid parameters",
			params: ExtensionParams{
				SettlementContract: "0x5678",
				AuctionDetails: &AuctionDetails{
					StartTime:       0,
					Duration:        0,
					InitialRateBump: 0,
					Points:          nil,
					GasCost:         GasCostConfigClassFixed{},
				},
				PostInteractionData: &SettlementPostInteractionData{
					Whitelist: []WhitelistItem{},
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(0),
						Receiver: common.Address{},
					},
					BankFee:            big.NewInt(0),
					ResolvingStartTime: big.NewInt(0),
					CustomReceiver:     common.Address{},
				},
				Asset:  "0x1234",
				Permit: "0x3456",

				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				Predicate:        "0x1234",
				PreInteraction:   "pre",
			},
			expectedExtension: &Extension{
				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				MakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
				TakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
				Predicate:        "0x1234",
				MakerPermit:      "0x00000000000000000000000000000000000012343456",
				PreInteraction:   "pre",
				PostInteraction:  "0x00000000000000000000000000000000000056780000000000",
			},
			expectErr: false,
		},
		{
			name: "Valid parameters 2",
			params: ExtensionParams{
				SettlementContract: "0x0500000000000000000000000000000000000000",
				AuctionDetails: &AuctionDetails{
					StartTime:       0,
					Duration:        0,
					InitialRateBump: 0,
					Points:          nil,
					GasCost:         GasCostConfigClassFixed{},
				},
				PostInteractionData: &SettlementPostInteractionData{
					Whitelist: []WhitelistItem{},
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(0),
						Receiver: common.Address{},
					},
					BankFee:            big.NewInt(0),
					ResolvingStartTime: big.NewInt(0),
					CustomReceiver:     common.Address{},
				},
				Asset:  "0x1234",
				Permit: "0x03",

				MakerAssetSuffix: "0x01",
				TakerAssetSuffix: "0x02",
				Predicate:        "0x07",
				PreInteraction:   "0x09",
			},
			expectedExtension: &Extension{
				MakerAssetSuffix: "0x01",
				TakerAssetSuffix: "0x02",
				MakingAmountData: "0x05000000000000000000000000000000000000000000000000000000000000000000000000",
				TakingAmountData: "0x05000000000000000000000000000000000000000000000000000000000000000000000000",
				Predicate:        "0x07",
				MakerPermit:      "0x000000000000000000000000000000000000123403",
				PreInteraction:   "0x09",
				PostInteraction:  "0x05000000000000000000000000000000000000000000000000",
			},
			expectErr: false,
		},
		{
			name: "Invalid MakerAssetSuffix",
			params: ExtensionParams{
				MakerAssetSuffix: "invalid",
				TakerAssetSuffix: "0x1234",
				Predicate:        "0x1234",
				PreInteraction:   "pre",
			},
			expectErr: true,
			errMsg:    "MakerAssetSuffix must be valid hex string",
		},
		{
			name: "Invalid TakerAssetSuffix",
			params: ExtensionParams{
				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "invalid",
				Predicate:        "0x1234",
				PreInteraction:   "pre",
			},
			expectErr: true,
			errMsg:    "TakerAssetSuffix must be valid hex string",
		},
		{
			name: "Invalid Predicate",
			params: ExtensionParams{
				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				Predicate:        "invalid",
				PreInteraction:   "pre",
			},
			expectErr: true,
			errMsg:    "Predicate must be valid hex string",
		},
		{
			name: "CustomData not supported",
			params: ExtensionParams{
				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				Predicate:        "0x1234",
				PreInteraction:   "pre",
				CustomData:       "0x1234",
			},
			expectErr: true,
			errMsg:    "CustomData is not currently supported",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ext, err := NewExtension(tc.params)
			if tc.expectErr {
				require.Error(t, err)
				assert.Equal(t, tc.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, ext)
				assert.Equal(t, tc.expectedExtension.MakerAssetSuffix, ext.MakerAssetSuffix)
				assert.Equal(t, tc.expectedExtension.TakerAssetSuffix, ext.TakerAssetSuffix)
				assert.Equal(t, tc.expectedExtension.MakingAmountData, ext.MakingAmountData)
				assert.Equal(t, tc.expectedExtension.TakingAmountData, ext.TakingAmountData)
				assert.Equal(t, tc.expectedExtension.Predicate, ext.Predicate)
				assert.Equal(t, tc.expectedExtension.MakerPermit, ext.MakerPermit)
				assert.Equal(t, tc.expectedExtension.PreInteraction, ext.PreInteraction)
				assert.Equal(t, tc.expectedExtension.PostInteraction, ext.PostInteraction)
			}
		})
	}
}

func TestDecodeExtensionPure(t *testing.T) {
	tests := []struct {
		name          string
		hexInput      string
		expected      *Extension
		expectingErr  bool
		errorContains string
	}{
		{
			name:     "Successful Decoding",
			hexInput: "0000007c00000063000000620000004d0000004c00000027000000020000000101020500000000000000000000000000000000000000000000000000000000000000000000000005000000000000000000000000000000000000000000000000000000000000000000000000070000000000000000000000000000000000001234030905000000000000000000000000000000000000000000000000",
			expected: &Extension{
				MakerAssetSuffix: "0x01",
				TakerAssetSuffix: "0x02",
				MakingAmountData: "05000000000000000000000000000000000000000000000000000000000000000000000000",
				TakingAmountData: "05000000000000000000000000000000000000000000000000000000000000000000000000",
				Predicate:        "0x07",
				MakerPermit:      "0x000000000000000000000000000000000000123403",
				PreInteraction:   "0x09",
				PostInteraction:  "0x05000000000000000000000000000000000000000000000000",
			},
			expectingErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert hex string to bytes
			data, err := hexToBytes(tt.hexInput)
			if err != nil {
				t.Fatalf("Failed to convert hex to bytes: %v", err)
			}

			// Decode the data
			decoded, err := DecodeExtensionPure(data)
			require.NoError(t, err)

			if tt.expectingErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errorContains != "" && !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s' but got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !extensionsEqual(decoded, tt.expected) {
					t.Errorf("Decoded Extension does not match expected.\nGot: %+v\nExpected: %+v", printSelectedFields(decoded), printSelectedFields(tt.expected))
				}
			}
		})
	}
}

func printSelectedFields(ext *Extension) string {
	selectedFields := map[string]string{
		"MakerAssetSuffix": strings.TrimPrefix(ext.MakerAssetSuffix, "0x"),
		"TakerAssetSuffix": strings.TrimPrefix(ext.TakerAssetSuffix, "0x"),
		"MakingAmountData": strings.TrimPrefix(ext.MakingAmountData, "0x"),
		"TakingAmountData": strings.TrimPrefix(ext.TakingAmountData, "0x"),
		"Predicate":        strings.TrimPrefix(ext.Predicate, "0x"),
		"MakerPermit":      strings.TrimPrefix(ext.MakerPermit, "0x"),
		"PreInteraction":   strings.TrimPrefix(ext.PreInteraction, "0x"),
		"PostInteraction":  strings.TrimPrefix(ext.PostInteraction, "0x"),
	}

	jsonData, err := json.MarshalIndent(selectedFields, "", "  ")
	if err != nil {
		return fmt.Sprint("Error marshalling to JSON:", err)
	}
	return string(jsonData)
}

func TestConvertToOrderbookExtensionPure(t *testing.T) {
	tests := []struct {
		name                       string
		fusionExtension            Extension
		expectedOrderbookExtension *orderbook.ExtensionPure
		expectErr                  bool
		errMsg                     string
	}{
		{
			name: "Valid parameters",
			fusionExtension: Extension{
				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				MakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
				TakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
				Predicate:        "0x1234",
				MakerPermit:      "0x00000000000000000000000000000000000012343456",
				PreInteraction:   "pre",
				PostInteraction:  "0x00000000000000000000000000000000000056780000000000",
			},
			expectedOrderbookExtension: &orderbook.ExtensionPure{
				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				MakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
				TakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
				Predicate:        "0x1234",
				MakerPermit:      "0x00000000000000000000000000000000000012343456",
				PreInteraction:   "pre",
				PostInteraction:  "0x00000000000000000000000000000000000056780000000000",
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ext := tc.fusionExtension.ConvertToOrderbookExtensionPure()
			assert.NotNil(t, ext)
			assert.Equal(t, tc.expectedOrderbookExtension.MakerAssetSuffix, ext.MakerAssetSuffix)
			assert.Equal(t, tc.expectedOrderbookExtension.TakerAssetSuffix, ext.TakerAssetSuffix)
			assert.Equal(t, tc.expectedOrderbookExtension.MakingAmountData, ext.MakingAmountData)
			assert.Equal(t, tc.expectedOrderbookExtension.TakingAmountData, ext.TakingAmountData)
			assert.Equal(t, tc.expectedOrderbookExtension.Predicate, ext.Predicate)
			assert.Equal(t, tc.expectedOrderbookExtension.MakerPermit, ext.MakerPermit)
			assert.Equal(t, tc.expectedOrderbookExtension.PreInteraction, ext.PreInteraction)
			assert.Equal(t, tc.expectedOrderbookExtension.PostInteraction, ext.PostInteraction)
		})
	}
}

var asset = "0xBAb2C3d4e5f67890123456789AbcDEf123456789"
var permit = "9999999999999999999999"
var fullAuctionDetails = &AuctionDetails{
	StartTime:       1,
	Duration:        2,
	InitialRateBump: 3,
	Points:          []AuctionPointClassFixed{{Coefficient: 4, Delay: 5}},
	GasCost:         GasCostConfigClassFixed{GasBumpEstimate: 6, GasPriceEstimate: 7},
}

var fullPostInteractionData = &SettlementPostInteractionData{
	Whitelist: []WhitelistItem{
		{
			AddressHalf: "a1b2c3d4e5f678901234",
			Delay:       big.NewInt(8),
		},
	},
	IntegratorFee: &IntegratorFee{
		Ratio:    big.NewInt(9),
		Receiver: common.HexToAddress("0xB1B2C3D4E5F67890123456789ABCDEF123456789"),
	},
	BankFee:            big.NewInt(10),
	ResolvingStartTime: big.NewInt(11),
	CustomReceiver:     common.HexToAddress("0xC1B2C3D4E5F67890123456789ABCDEF123456789"),
}

func TestFromExtension(t *testing.T) {
	tests := []struct {
		name              string
		params            ExtensionParams
		expectedExtension *Extension
		expectErr         bool
		errMsg            string
	}{
		{
			name: "Valid parameters",
			params: ExtensionParams{
				SettlementContract:  "0xAAB2C3d4E5F67890123456789abcdef123456789",
				AuctionDetails:      fullAuctionDetails,
				PostInteractionData: fullPostInteractionData,
				Asset:               asset,
				Permit:              permit,

				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				Predicate:        "0x1234",
				PreInteraction:   "pre",
			},
			expectedExtension: &Extension{
				SettlementContract:  "0xAAB2C3d4E5F67890123456789abcdef123456789",
				AuctionDetails:      fullAuctionDetails,
				PostInteractionData: fullPostInteractionData,
				Asset:               asset,
				Permit:              permit,

				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				MakingAmountData: "0xAAB2C3d4E5F67890123456789abcdef12345678900000600000007000000010000020000030000040005",
				TakingAmountData: "0xAAB2C3d4E5F67890123456789abcdef12345678900000600000007000000010000020000030000040005",
				Predicate:        "0x1234",
				MakerPermit:      fmt.Sprintf("%s%s", asset, permit),
				PreInteraction:   "pre",
				PostInteraction:  "0xAAB2C3d4E5F67890123456789abcdef1234567890000000a0009b1b2c3d4e5f67890123456789abcdef123456789c1b2c3d4e5f67890123456789abcdef1234567890000000ba1b2c3d4e5f67890123400080f",
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ext, err := NewExtension(tc.params)
			require.NoError(t, err)

			limitOrderExtensionPure := ext.ConvertToOrderbookExtensionPure()
			decodedExtension, err := FromLimitOrderExtensionPure(limitOrderExtensionPure)
			require.NoError(t, err)

			assert.NotNil(t, ext)
			assert.Equal(t, tc.expectedExtension.SettlementContract, decodedExtension.SettlementContract)
			assert.Equal(t, tc.expectedExtension.AuctionDetails, decodedExtension.AuctionDetails)
			assert.Equal(t, tc.expectedExtension.PostInteractionData, decodedExtension.PostInteractionData)
			//assert.Equal(t, tc.expectedExtension.Asset, decodedExtension.Asset)
			//assert.Equal(t, tc.expectedExtension.Permit, decodedExtension.Permit)

			assert.Equal(t, tc.expectedExtension.MakerAssetSuffix, decodedExtension.MakerAssetSuffix)
			assert.Equal(t, tc.expectedExtension.TakerAssetSuffix, decodedExtension.TakerAssetSuffix)
			assert.Equal(t, tc.expectedExtension.MakingAmountData, decodedExtension.MakingAmountData)
			assert.Equal(t, tc.expectedExtension.TakingAmountData, decodedExtension.TakingAmountData)
			assert.Equal(t, tc.expectedExtension.Predicate, decodedExtension.Predicate)
			assert.Equal(t, tc.expectedExtension.MakerPermit, decodedExtension.MakerPermit)
			assert.Equal(t, tc.expectedExtension.PreInteraction, decodedExtension.PreInteraction)
			assert.Equal(t, tc.expectedExtension.PostInteraction, decodedExtension.PostInteraction)
		})
	}
}

func extensionsEqual(a, b *Extension) bool {
	return strings.TrimPrefix(a.MakerAssetSuffix, "0x") == strings.TrimPrefix(b.MakerAssetSuffix, "0x") &&
		strings.TrimPrefix(a.TakerAssetSuffix, "0x") == strings.TrimPrefix(b.TakerAssetSuffix, "0x") &&
		strings.TrimPrefix(a.MakingAmountData, "0x") == strings.TrimPrefix(b.MakingAmountData, "0x") &&
		strings.TrimPrefix(a.TakingAmountData, "0x") == strings.TrimPrefix(b.TakingAmountData, "0x") &&
		strings.TrimPrefix(a.Predicate, "0x") == strings.TrimPrefix(b.Predicate, "0x") &&
		strings.TrimPrefix(a.MakerPermit, "0x") == strings.TrimPrefix(b.MakerPermit, "0x") &&
		strings.TrimPrefix(a.PreInteraction, "0x") == strings.TrimPrefix(b.PreInteraction, "0x") &&
		strings.TrimPrefix(a.PostInteraction, "0x") == strings.TrimPrefix(b.PostInteraction, "0x")
	// strings.TrimPrefix(a.CustomData, "0x") == strings.TrimPrefix(b.CustomData, "0x")
}

// hexToBytes converts a hexadecimal string to a byte slice.
func hexToBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// contains checks if the substring is present in the string.
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
