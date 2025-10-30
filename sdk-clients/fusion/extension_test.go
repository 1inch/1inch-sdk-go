package fusion

import (
	"math/big"
	"testing"

	"github.com/1inch/1inch-sdk-go/internal/bigint"
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
					AuctionFees: &FeesIntegratorAndResolver{
						Resolver:   ResolverFee{},
						Integrator: IntegratorFee{},
					},
					ResolvingStartTime: big.NewInt(0),
					CustomReceiver:     common.Address{},
				},
				Asset:  "0x1234",
				Permit: "0x3456",

				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				Predicate:        "0x1234",
				PreInteraction:   "0x5678",

				Surplus: &SurplusParams{
					EstimatedTakerAmount: big.NewInt(1),
					ProtocolFee:          FromPercent(1, GetDefaultBase()),
				},

				ResolvingStartTime: big.NewInt(0),
			},
			expectedExtension: &Extension{
				MakerAssetSuffix: "0x1234",
				TakerAssetSuffix: "0x1234",
				MakingAmountData: "0x000000000000000000000000000000000000567800000000000000000000000000000000000000000000006400",
				TakingAmountData: "0x000000000000000000000000000000000000567800000000000000000000000000000000000000000000006400",
				Predicate:        "0x1234",
				MakerPermit:      "0x00000000000000000000000000000000000012343456",
				PreInteraction:   "0x5678",
				PostInteraction:  "0x000000000000000000000000000000000000567800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000000000000000101",
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
					AuctionFees: &FeesIntegratorAndResolver{
						Resolver:   ResolverFee{},
						Integrator: IntegratorFee{},
					},
					ResolvingStartTime: big.NewInt(0),
					CustomReceiver:     common.Address{},
				},
				Asset:  "0x1234",
				Permit: "0x03",

				MakerAssetSuffix: "0x01",
				TakerAssetSuffix: "0x02",
				Predicate:        "0x07",
				PreInteraction:   "0x09",

				Surplus: &SurplusParams{
					EstimatedTakerAmount: big.NewInt(1),
					ProtocolFee:          FromPercent(1, GetDefaultBase()),
				},

				ResolvingStartTime: big.NewInt(0),
			},
			expectedExtension: &Extension{
				MakerAssetSuffix: "0x01",
				TakerAssetSuffix: "0x02",
				MakingAmountData: "0x050000000000000000000000000000000000000000000000000000000000000000000000000000000000006400",
				TakingAmountData: "0x050000000000000000000000000000000000000000000000000000000000000000000000000000000000006400",
				Predicate:        "0x07",
				MakerPermit:      "0x000000000000000000000000000000000000123403",
				PreInteraction:   "0x09",
				PostInteraction:  "0x050000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000000000000000101",
			},
			expectErr: false,
		},
		{
			name: "Invalid SettlementContract",
			params: ExtensionParams{
				SettlementContract: "invalid",
				MakerAssetSuffix:   "0x1234",
				TakerAssetSuffix:   "0x1234",
				Predicate:          "0x1234",
				PreInteraction:     "0x5678",
			},
			expectErr: true,
			errMsg:    "Settlement contract must be valid hex string",
		},
		{
			name: "Invalid MakerAssetSuffix",
			params: ExtensionParams{
				SettlementContract: "0x9012",
				MakerAssetSuffix:   "invalid",
				TakerAssetSuffix:   "0x1234",
				Predicate:          "0x1234",
				PreInteraction:     "0x5678",
			},
			expectErr: true,
			errMsg:    "MakerAssetSuffix must be valid hex string",
		},
		{
			name: "Invalid TakerAssetSuffix",
			params: ExtensionParams{
				SettlementContract: "0x9012",
				MakerAssetSuffix:   "0x1234",
				TakerAssetSuffix:   "invalid",
				Predicate:          "0x1234",
				PreInteraction:     "0x5678",
			},
			expectErr: true,
			errMsg:    "TakerAssetSuffix must be valid hex string",
		},
		{
			name: "Invalid Predicate",
			params: ExtensionParams{
				SettlementContract: "0x9012",
				MakerAssetSuffix:   "0x1234",
				TakerAssetSuffix:   "0x1234",
				Predicate:          "invalid",
				PreInteraction:     "0x5678",
			},
			expectErr: true,
			errMsg:    "Predicate must be valid hex string",
		},
		{
			name: "CustomData not supported",
			params: ExtensionParams{
				SettlementContract: "0x9012",
				MakerAssetSuffix:   "0x1234",
				TakerAssetSuffix:   "0x1234",
				Predicate:          "0x1234",
				PreInteraction:     "0x5678",
				CustomData:         "0x1234",
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

func TestConvertToOrderbookExtension(t *testing.T) {
	tests := []struct {
		name                       string
		fusionExtension            Extension
		expectedOrderbookExtension *orderbook.Extension
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
			expectedOrderbookExtension: &orderbook.Extension{
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
			ext := tc.fusionExtension.ConvertToOrderbookExtension()
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

func TestBuildAmountGetterData(t *testing.T) {
	tests := []struct {
		name             string
		details          *AuctionDetails
		detailsFull      *Details
		whitelist        []WhitelistItem
		forAmountGetters bool
		expected         string
	}{
		{
			name: "basic auction details with forAmountGetters true",
			detailsFull: &Details{
				Auction: &AuctionDetails{
					GasCost: GasCostConfigClassFixed{
						GasBumpEstimate:  0,
						GasPriceEstimate: 0,
					},
					StartTime:       1673548149,
					Duration:        180,
					InitialRateBump: 50000,
					Points: []AuctionPointClassFixed{
						{
							Delay:       12,
							Coefficient: 20000,
						},
					},
				},
				ResolvingStartTime: big.NewInt(1673548139),
			},
			whitelist: []WhitelistItem{
				{
					AddressHalf: "bb839cbe05303d7705fa",
					Delay:       big.NewInt(0),
				},
			},
			forAmountGetters: true,
			expected:         "0x0000000000000063c051750000b400c35001004e20000c00000000006401bb839cbe05303d7705fa",
		},
		{
			name: "basic auction details with forAmountGetters false",
			detailsFull: &Details{
				Auction: &AuctionDetails{
					GasCost: GasCostConfigClassFixed{
						GasBumpEstimate:  0,
						GasPriceEstimate: 0,
					},
					StartTime:       1673548149,
					Duration:        180,
					InitialRateBump: 50000,
					Points: []AuctionPointClassFixed{
						{
							Delay:       12,
							Coefficient: 20000,
						},
					},
				},
				ResolvingStartTime: big.NewInt(1673548139),
			},
			whitelist: []WhitelistItem{
				{
					AddressHalf: "bb839cbe05303d7705fa",
					Delay:       big.NewInt(0),
				},
			},
			forAmountGetters: false,
			expected:         "0x00000000006463c0516b01bb839cbe05303d7705fa0000",
		},
		{
			name: "with fees",
			detailsFull: &Details{
				Auction: &AuctionDetails{
					StartTime:       1673548149,
					Duration:        180,
					InitialRateBump: 50000,
					Points:          []AuctionPointClassFixed{{Coefficient: 20000, Delay: 12}},
				},
				FeesIntAndRes: &FeesIntegratorAndResolver{
					Integrator: IntegratorFee{
						Integrator: "0x0000000000000000000000000000000000000001",
						Protocol:   "0x0000000000000000000000000000000000000002",
						Fee:        FromPercent(1, GetDefaultBase()),
						Share:      FromPercent(50, GetDefaultBase()),
					},
				},
				Whitelist: []AuctionWhitelistItem{
					{Address: common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"), AllowFrom: big.NewInt(0)},
				},
				ResolvingStartTime: big.NewInt(1673548139),
			},
			whitelist: []WhitelistItem{
				{
					AddressHalf: "bb839cbe05303d7705fa",
					Delay:       big.NewInt(0),
				},
			},
			forAmountGetters: true,
			expected:         "0x0000000000000063c051750000b400c35001004e20000c03e83200006401bb839cbe05303d7705fa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &BuildAmountGetterDataParams{
				AuctionDetails:     tt.detailsFull.Auction,
				ResolvingStartTime: tt.detailsFull.ResolvingStartTime,
				PostInteractionData: &SettlementPostInteractionData{
					Whitelist:          tt.whitelist,
					ResolvingStartTime: tt.detailsFull.ResolvingStartTime,
					CustomReceiver:     common.Address{},
					AuctionFees:        tt.detailsFull.FeesIntAndRes,
				},
			}

			gotHex, err := BuildAmountGetterData(params, tt.forAmountGetters)
			require.NoError(t, err)
			require.Equal(t, tt.expected, gotHex)
		})
	}
}
