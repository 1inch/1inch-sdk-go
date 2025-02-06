package fusionplus

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/1inch/1inch-sdk-go/internal/bigint"
	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/fusion"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		extension *EscrowExtension
		expected  string
		expectErr bool
	}{
		{
			name: "Generate salt when extension is not empty",
			extension: &EscrowExtension{
				Extension: fusion.Extension{
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
			},
			expected:  "180431178743033967347942937469468920088249224033532329",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expected, err := bigint.FromString(tc.expected)
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
		name      string
		params    EscrowExtensionParams
		expected  *EscrowExtension
		expectErr bool
		errMsg    string
	}{
		{
			name: "Valid parameters with Escrow",
			params: EscrowExtensionParams{
				ExtensionParams: fusion.ExtensionParams{
					SettlementContract: "0x5678",
					AuctionDetails: &fusion.AuctionDetails{
						StartTime:       0,
						Duration:        0,
						InitialRateBump: 0,
						Points:          nil,
						GasCost:         fusion.GasCostConfigClassFixed{},
					},
					PostInteractionData: &fusion.SettlementPostInteractionData{
						Whitelist: []fusion.WhitelistItem{},
						IntegratorFee: &fusion.IntegratorFee{
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
			},
			expected: &EscrowExtension{
				Extension: fusion.Extension{
					MakerAssetSuffix: "0x1234",
					TakerAssetSuffix: "0x1234",
					MakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
					TakingAmountData: "0x00000000000000000000000000000000000056780000000000000000000000000000000000",
					Predicate:        "0x1234",
					MakerPermit:      "0x00000000000000000000000000000000000012343456",
					PreInteraction:   "pre",
					PostInteraction:  "0x00000000000000000000000000000000000056780000000000",
				},
			},
			expectErr: false,
		},
		{
			name: "Invalid MakerAssetSuffix",
			params: EscrowExtensionParams{
				ExtensionParams: fusion.ExtensionParams{
					SettlementContract: "0x5678",
					MakerAssetSuffix:   "invalid",
					TakerAssetSuffix:   "0x1234",
					Predicate:          "0x1234",
					PreInteraction:     "pre",
				},
			},
			expectErr: true,
			errMsg:    "MakerAssetSuffix must be valid hex string",
		},
		{
			name: "Invalid TakerAssetSuffix",
			params: EscrowExtensionParams{
				ExtensionParams: fusion.ExtensionParams{
					SettlementContract: "0x5678",
					MakerAssetSuffix:   "0x1234",
					TakerAssetSuffix:   "invalid",
					Predicate:          "0x1234",
					PreInteraction:     "pre",
				},
			},
			expectErr: true,
			errMsg:    "TakerAssetSuffix must be valid hex string",
		},
		{
			name: "Invalid Predicate",
			params: EscrowExtensionParams{
				ExtensionParams: fusion.ExtensionParams{
					SettlementContract: "0x5678",
					MakerAssetSuffix:   "0x1234",
					TakerAssetSuffix:   "0x1234",
					Predicate:          "invalid",
					PreInteraction:     "pre",
				},
			},
			expectErr: true,
			errMsg:    "Predicate must be valid hex string",
		},
		{
			name: "CustomData not supported",
			params: EscrowExtensionParams{
				ExtensionParams: fusion.ExtensionParams{
					SettlementContract: "0x5678",
					MakerAssetSuffix:   "0x1234",
					TakerAssetSuffix:   "0x1234",
					Predicate:          "0x1234",
					PreInteraction:     "pre",
					CustomData:         "0x1234",
				},
			},
			expectErr: true,
			errMsg:    "CustomData is not currently supported",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ext, err := NewEscrowExtension(tc.params)
			if tc.expectErr {
				require.Error(t, err)
				assert.Equal(t, tc.errMsg, err.Error())
			} else {
				require.NoError(t, err)
				assert.NotNil(t, ext)
				assert.Equal(t, tc.expected.MakerAssetSuffix, ext.MakerAssetSuffix)
				assert.Equal(t, tc.expected.TakerAssetSuffix, ext.TakerAssetSuffix)
				assert.Equal(t, tc.expected.Predicate, ext.Predicate)
				assert.Equal(t, tc.expected.PreInteraction, ext.PreInteraction)
				assert.Equal(t, tc.expected.PostInteraction, ext.PostInteraction)
				assert.Equal(t, tc.expected.CustomData, ext.CustomData)
			}
		})
	}
}

// TestDecodeEscrowExtension contains all unit tests for the DecodeEscrowExtension function.
func TestDecodeEscrowExtension(t *testing.T) {
	tests := []struct {
		name          string
		hexInput      string
		expected      *EscrowExtension
		expectingErr  bool
		errorContains string
	}{
		{
			name:     "Full decode",
			hexInput: "0x0000016b0000005e0000005e0000005e0000005e0000002f0000000000000000fb2809a5314473e1165f6b58018e20ed8f07b84000b8460000222c6656b88f0000b401e0da00ba01009000b8460024fb2809a5314473e1165f6b58018e20ed8f07b84000b8460000222c6656b88f0000b401e0da00ba01009000b8460024fb2809a5314473e1165f6b58018e20ed8f07b8406656b877b09498030ae3416b66dc0000db05a6a504f04d92e79d00000c989d73cf0bd5f83b660000d18bd45f0b94f54a968f0000d61b892b2ad6249011850000d0847e80c0b823a65ce70000901f8f650d76dcc657d1000038ad1723a873d05effcbdc57dcf7d00458d6a8c763558d5af7522bf6ad2d3e253d000000000000000000000000000000000000000000000000000000000000a4b1000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000064000000000000000000000000000000c80000000000000003000000020000000100000004000000030000000200000001",
			expected: &EscrowExtension{
				Extension: fusion.Extension{
					SettlementContract: "0xfb2809a5314473e1165f6b58018e20ed8f07b840",
					AuctionDetails: &fusion.AuctionDetails{
						StartTime:       1716959375,
						InitialRateBump: 123098,
						Duration:        180,
						Points: []fusion.AuctionPointClassFixed{
							{
								Coefficient: 47617,
								Delay:       144,
							},
							{
								Coefficient: 47174,
								Delay:       36,
							},
						},
						GasCost: fusion.GasCostConfigClassFixed{
							GasBumpEstimate:  47174,
							GasPriceEstimate: 8748,
						},
					},
					PostInteractionData: &fusion.SettlementPostInteractionData{
						Whitelist: []fusion.WhitelistItem{
							{
								AddressHalf: "b09498030ae3416b66dc",
								Delay:       big.NewInt(0),
							},
							{
								AddressHalf: "db05a6a504f04d92e79d",
								Delay:       big.NewInt(0),
							},
							{
								AddressHalf: "0c989d73cf0bd5f83b66",
								Delay:       big.NewInt(0),
							},
							{
								AddressHalf: "d18bd45f0b94f54a968f",
								Delay:       big.NewInt(0),
							},
							{
								AddressHalf: "d61b892b2ad624901185",
								Delay:       big.NewInt(0),
							},
							{
								AddressHalf: "d0847e80c0b823a65ce7",
								Delay:       big.NewInt(0),
							},
							{
								AddressHalf: "901f8f650d76dcc657d1",
								Delay:       big.NewInt(0),
							},
						},
						BankFee:            nil,
						ResolvingStartTime: big.NewInt(1716959351),
					},
					MakerAssetSuffix: "0x",
					TakerAssetSuffix: "0x",
					MakingAmountData: "0xfb2809a5314473e1165f6b58018e20ed8f07b84000b8460000222c6656b88f0000b401e0da00ba01009000b8460024",
					TakingAmountData: "0xfb2809a5314473e1165f6b58018e20ed8f07b84000b8460000222c6656b88f0000b401e0da00ba01009000b8460024",
					Predicate:        "0x",
					MakerPermit:      "0x",
					PreInteraction:   "0x",
					PostInteraction:  "0xfb2809a5314473e1165f6b58018e20ed8f07b8406656b877b09498030ae3416b66dc0000db05a6a504f04d92e79d00000c989d73cf0bd5f83b660000d18bd45f0b94f54a968f0000d61b892b2ad6249011850000d0847e80c0b823a65ce70000901f8f650d76dcc657d1000038",
					CustomData:       "0x",
				},
				HashLock: &HashLock{
					Value: "0xed17b7cc09d7a0ba79bce96c0f0ec59d15e63bceeeae147ed230cff89689ce5c",
				},
				DstChainId:       42161,
				DstToken:         common.HexToAddress("0x0000000000000000000000000000000000000001"),
				SrcSafetyDeposit: "100",
				DstSafetyDeposit: "200",
				TimeLocks: TimeLocks{
					DstCancellation:       3,
					DstPublicWithdrawal:   2,
					DstWithdrawal:         1,
					SrcPublicCancellation: 4,
					SrcCancellation:       3,
					SrcPublicWithdrawal:   2,
					SrcWithdrawal:         1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert hex string to bytes
			data := hexutil.MustDecode(tt.hexInput)

			// Decode the data
			decoded, err := DecodeEscrowExtension(data)
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
				assert.Equal(t, tt.expected.SettlementContract, decoded.SettlementContract)
				assert.Equal(t, tt.expected.AuctionDetails, decoded.AuctionDetails)
				assert.Equal(t, tt.expected.PostInteractionData, decoded.PostInteractionData)
				assert.Equal(t, tt.expected.Asset, decoded.Asset)
				assert.Equal(t, tt.expected.Permit, decoded.Permit)
				assert.Equal(t, tt.expected.MakerAssetSuffix, decoded.MakerAssetSuffix)
				assert.Equal(t, tt.expected.TakerAssetSuffix, decoded.TakerAssetSuffix)
				assert.Equal(t, tt.expected.MakingAmountData, decoded.MakingAmountData)
				assert.Equal(t, tt.expected.TakingAmountData, decoded.TakingAmountData)
				assert.Equal(t, tt.expected.Predicate, decoded.Predicate)
				assert.Equal(t, tt.expected.MakerPermit, decoded.MakerPermit)
				assert.Equal(t, tt.expected.PreInteraction, decoded.PreInteraction)
				assert.Equal(t, tt.expected.PostInteraction, decoded.PostInteraction)
				assert.Equal(t, tt.expected.TimeLocks, decoded.TimeLocks)
			}
		})
	}
}

func TestEncodeExtraData(t *testing.T) {
	tests := []struct {
		name            string
		expectedEncoded string
		extraData       *EscrowExtraData
		expectingErr    bool
		errorContains   string
	}{
		{
			name: "Encode without any other data",
			extraData: &EscrowExtraData{
				HashLock: &HashLock{
					Value: "ad1723a873d05effcbdc57dcf7d00458d6a8c763558d5af7522bf6ad2d3e253d",
				},
				DstChainId:       42161,
				DstToken:         common.HexToAddress("0x0000000000000000000000000000000000000001"),
				SrcSafetyDeposit: big.NewInt(100),
				DstSafetyDeposit: big.NewInt(200),
				TimeLocks: &TimeLocks{
					DstCancellation:       3,
					DstPublicWithdrawal:   2,
					DstWithdrawal:         1,
					SrcPublicCancellation: 4,
					SrcCancellation:       3,
					SrcPublicWithdrawal:   2,
					SrcWithdrawal:         1,
				},
			},
			expectedEncoded: "ad1723a873d05effcbdc57dcf7d00458d6a8c763558d5af7522bf6ad2d3e253d000000000000000000000000000000000000000000000000000000000000a4b1000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000064000000000000000000000000000000c80000000000000003000000020000000100000004000000030000000200000001",
			expectingErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			encoded, err := encodeExtraData(tt.extraData)
			require.NoError(t, err)

			require.Equal(t, tt.expectedEncoded, fmt.Sprintf("%x", encoded))
		})
	}
}

// contains checks if the substring is present in the string.
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
