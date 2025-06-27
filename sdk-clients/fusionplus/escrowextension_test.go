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
			expected:  "180431909497609865807168059378624320943465639784996571",
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

//TODO commenting out until Fusion is working again
//func TestNewExtension(t *testing.T) {
//	tests := []struct {
//		name      string
//		params    EscrowExtensionParams
//		expected  *EscrowExtension
//		expectErr bool
//		errMsg    string
//	}{
//		{
//			name: "Valid parameters with Escrow",
//			params: EscrowExtensionParams{
//				ExtensionParams: fusion.ExtensionParams{
//					SettlementContract: "0x5678",
//					AuctionDetails: &fusion.AuctionDetails{
//						StartTime:       0,
//						Duration:        0,
//						InitialRateBump: 0,
//						Points:          nil,
//						GasCost:         fusion.GasCostConfigClassFixed{},
//					},
//					PostInteractionData: &fusion.SettlementPostInteractionData{
//						Whitelist: []fusion.WhitelistItem{},
//						IntegratorFee: &fusion.IntegratorFee{
//							Ratio:    big.NewInt(0),
//							Receiver: common.Address{},
//						},
//						BankFee:            big.NewInt(0),
//						ResolvingStartTime: big.NewInt(0),
//						CustomReceiver:     common.Address{},
//					},
//				},
//			},
//			expected: &EscrowExtension{
//				Extension: fusion.Extension{
//					MakerAssetSuffix: "0x1234",
//					TakerAssetSuffix: "0x1234",
//					MakingAmountData: "0x0000000000000000000000000000000000005678",
//					TakingAmountData: "0x0000000000000000000000000000000000005678",
//					Predicate:        "0x1234",
//					MakerPermit:      "0x00000000000000000000000000000000000012343456",
//					PreInteraction:   "0xpre",
//					PostInteraction:  "0x0000000000000000000000000000000000005678",
//				},
//			},
//			expectErr: false,
//		},
//		{
//			name: "Invalid MakerAssetSuffix",
//			params: EscrowExtensionParams{
//				ExtensionParams: fusion.ExtensionParams{
//					SettlementContract: "0x5678",
//					MakerAssetSuffix:   "invalid",
//					TakerAssetSuffix:   "0x1234",
//					Predicate:          "0x1234",
//					PreInteraction:     "pre",
//				},
//			},
//			expectErr: true,
//			errMsg:    "MakerAssetSuffix must be valid hex string",
//		},
//		{
//			name: "Invalid TakerAssetSuffix",
//			params: EscrowExtensionParams{
//				ExtensionParams: fusion.ExtensionParams{
//					SettlementContract: "0x5678",
//					MakerAssetSuffix:   "0x1234",
//					TakerAssetSuffix:   "invalid",
//					Predicate:          "0x1234",
//					PreInteraction:     "pre",
//				},
//			},
//			expectErr: true,
//			errMsg:    "TakerAssetSuffix must be valid hex string",
//		},
//		{
//			name: "Invalid Predicate",
//			params: EscrowExtensionParams{
//				ExtensionParams: fusion.ExtensionParams{
//					SettlementContract: "0x5678",
//					MakerAssetSuffix:   "0x1234",
//					TakerAssetSuffix:   "0x1234",
//					Predicate:          "invalid",
//					PreInteraction:     "pre",
//				},
//			},
//			expectErr: true,
//			errMsg:    "Predicate must be valid hex string",
//		},
//		{
//			name: "CustomData not supported",
//			params: EscrowExtensionParams{
//				ExtensionParams: fusion.ExtensionParams{
//					SettlementContract: "0x5678",
//					MakerAssetSuffix:   "0x1234",
//					TakerAssetSuffix:   "0x1234",
//					Predicate:          "0x1234",
//					PreInteraction:     "pre",
//					CustomData:         "0x1234",
//				},
//			},
//			expectErr: true,
//			errMsg:    "CustomData is not currently supported",
//		},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			ext, err := NewEscrowExtension(tc.params)
//			if tc.expectErr {
//				require.Error(t, err)
//				assert.Equal(t, tc.errMsg, err.Error())
//			} else {
//				require.NoError(t, err)
//				assert.NotNil(t, ext)
//				assert.Equal(t, tc.expected.MakerAssetSuffix, ext.MakerAssetSuffix)
//				assert.Equal(t, tc.expected.TakerAssetSuffix, ext.TakerAssetSuffix)
//				assert.Equal(t, tc.expected.Predicate, ext.Predicate)
//				assert.Equal(t, tc.expected.PreInteraction, ext.PreInteraction)
//				assert.Equal(t, tc.expected.PostInteraction, ext.PostInteraction)
//			}
//		})
//	}
//}

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
