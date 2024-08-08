package fusion

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	web3_provider "github.com/1inch/1inch-sdk-go/internal/web3-provider"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

var (
	publicAddress = "0x737baD27cF1374AE2af29C49Bb6D6007D5CD67EE"
	privateKey    = "0f3edf983ac636a65a842ce7c78d9aa706d3b113b37e265ba6b02d758e70b3d0"
)

const (
	usdc    = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	wmatic  = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	amount  = "1000000000000000000"
	chainId = 137
)

func TestCreateFusionOrderData(t *testing.T) {
	tests := []struct {
		name                        string
		chainId                     uint64
		privateKey                  string
		orderParams                 OrderParams
		additionalParams            AdditionalParams
		auctionStartTime            uint32
		nonce                       *big.Int
		resolverStartTime           int64
		baseSaltValue               string
		serializedQuoteData         string
		serializedPreparedOrderData string
		serializedLimitOrderData    string
		data                        string
	}{
		{
			name:       "Successful order creation",
			chainId:    chainId,
			privateKey: privateKey,
			orderParams: OrderParams{
				WalletAddress:    publicAddress,
				FromTokenAddress: wmatic,
				ToTokenAddress:   usdc,
				Amount:           amount,
				Receiver:         "0x0000000000000000000000000000000000000000",
				Preset:           "fast",
			},
			auctionStartTime:            1718671900,
			nonce:                       big.NewInt(887174712009),
			resolverStartTime:           1718671883,
			baseSaltValue:               "35020243109857195061155306569",
			serializedQuoteData:         `{"feeToken":"0x3c499c542cef5e3811e1192ce70d8cc03d5c3359","fromTokenAmount":"1000000000000000000","presets":{"fast":{"allowMultipleFills":false,"allowPartialFills":false,"auctionDuration":180,"auctionEndAmount":"538946","auctionStartAmount":"557310","bankFee":"0","estP":100,"exclusiveResolver":null,"gasCost":{"gasBumpEstimate":0,"gasPriceEstimate":"0"},"initialRateBump":340757,"points":[],"startAuctionIn":17,"tokenFee":"18366"},"medium":{"allowMultipleFills":true,"allowPartialFills":true,"auctionDuration":360,"auctionEndAmount":"538946","auctionStartAmount":"576251","bankFee":"0","estP":100,"exclusiveResolver":null,"gasCost":{"gasBumpEstimate":0,"gasPriceEstimate":"0"},"initialRateBump":692202,"points":[{"coefficient":681533,"delay":6},{"coefficient":340757,"delay":6}],"startAuctionIn":17,"tokenFee":"18366"},"slow":{"allowMultipleFills":true,"allowPartialFills":true,"auctionDuration":600,"auctionEndAmount":"538946","auctionStartAmount":"581432","bankFee":"0","estP":100,"exclusiveResolver":null,"gasCost":{"gasBumpEstimate":0,"gasPriceEstimate":"0"},"initialRateBump":788335,"points":[{"coefficient":681533,"delay":81},{"coefficient":340757,"delay":6}],"startAuctionIn":17,"tokenFee":"18366"}},"prices":{"usd":{"fromToken":"0.57493897","toToken":"0.9995015368854032"}},"quoteId":"55c3f478-b176-448c-b968-656c19b9c04a","recommended_preset":"fast","settlementAddress":"0xfb2809a5314473e1165f6b58018e20ed8f07b840","suggested":true,"toTokenAmount":"575677","volume":{"usd":{"fromToken":"0.57493897","toToken":"0.57539"}},"whitelist":["0x46fd018b32a9315ef5b4c0866635457d36ab318d","0xc1b19a08c2798c6930b8f3a44b7b0d08f4e198b8","0x0000000000000000000000000000000000000000","0xad3b67bca8935cb510c8d18bd45f0b94f54a968f","0x0000000000000000000000000000000000000000","0x0000000000000000000000000000000000000000","0x62f861201db5fdc04c48c976bf098c4dba0a061d","0x0000000000000000000000000000000000000000"]}`,
			serializedPreparedOrderData: `{"order":{"FusionExtension":{"MakerAssetSuffix":"","TakerAssetSuffix":"","MakingAmountData":"0xfb2809A5314473E1165f6B58018E20ed8F07B840000000000000006670da1c0000b4053315","TakingAmountData":"0xfb2809A5314473E1165f6B58018E20ed8F07B840000000000000006670da1c0000b4053315","Predicate":"","MakerPermit":"","PreInteraction":"","PostInteraction":"0xfb2809A5314473E1165f6B58018E20ed8F07B8406670da0bc0866635457d36ab318d0000f3a44b7b0d08f4e198b80000000000000000000000000000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d000000000000000000000000000040","CustomData":""},"Inner":{"makerAsset":"0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270","takerAsset":"0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359","makingAmount":"1000000000000000000","takingAmount":"538946","salt":"712810ef08aca692b6d59c49fc131590b1edc52d382c2a9684cae76e49ca45bf","maker":"0x737baD27cF1374AE2af29C49Bb6D6007D5CD67EE","allowedSender":"","receiver":"0x0000000000000000000000000000000000000000","makerTraits":"0x8a0000000000000000000000ce8fbbcac9006670dad000000000000000000000","extension":"357969f7ed9a797c95a9da11fc131590b1edc52d382c2a9684cae76e49ca45bf"},"SettlementExtension":"0xfb2809a5314473e1165f6b58018e20ed8f07b840","OrderInfo":{"maker":"0x737baD27cF1374AE2af29C49Bb6D6007D5CD67EE","makerAsset":"0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270","makerTraits":"","makingAmount":"1000000000000000000","receiver":"0x0000000000000000000000000000000000000000","salt":"","takerAsset":"0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359","takingAmount":"538946"},"AuctionDetails":{"startTime":1718671900,"duration":180,"initialRateBump":340757,"points":[],"gasCost":{"gasBumpEstimate":0,"gasPriceEstimate":0}},"PostInteractionData":{"Whitelist":[{"AddressHalf":"c0866635457d36ab318d","Delay":0},{"AddressHalf":"f3a44b7b0d08f4e198b8","Delay":0},{"AddressHalf":"00000000000000000000","Delay":0},{"AddressHalf":"d18bd45f0b94f54a968f","Delay":0},{"AddressHalf":"00000000000000000000","Delay":0},{"AddressHalf":"00000000000000000000","Delay":0},{"AddressHalf":"c976bf098c4dba0a061d","Delay":0},{"AddressHalf":"00000000000000000000","Delay":0}],"IntegratorFee":{"Ratio":0,"Receiver":"0x0000000000000000000000000000000000000000"},"BankFee":0,"ResolvingStartTime":1718671883,"CustomReceiver":"0x0000000000000000000000000000000000000000"},"Extra":{"UnwrapWETH":false,"Nonce":887174712009,"Permit":"","AllowPartialFills":false,"AllowMultipleFills":false,"OrderExpirationDelay":0,"EnablePermit2":false,"Source":""}},"hash":"0xe635531055466f92fdf64222d3e6d5ff18cda08c1a87b28c6853095d50699574","quoteId":"55c3f478-b176-448c-b968-656c19b9c04a"}`,
			serializedLimitOrderData:    `{"orderHash":"0xe635531055466f92fdf64222d3e6d5ff18cda08c1a87b28c6853095d50699574","signature":"0xa1cb6463f2e9126fe24e5b8f1f0bb3762ed588fc0e8c7186cfa81f19806127cd21a37b8c9ee812429a2449f926736d32b1e2108f7aae8f5e96802a2d35e242781b","data":{"makerAsset":"0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270","takerAsset":"0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359","makingAmount":"1000000000000000000","takingAmount":"538946","salt":"0x9a042bfb67cf14b0a1a98c4ae5d6295e2c08820","maker":"0x737baD27cF1374AE2af29C49Bb6D6007D5CD67EE","allowedSender":"0x0000000000000000000000000000000000000000","receiver":"0x0000000000000000000000000000000000000000","makerTraits":"0x8a0000000000000000000000ce8fbbcac9006670dad000000000000000000000","extension":"0x000000c30000004a0000004a0000004a0000004a000000250000000000000000fb2809A5314473E1165f6B58018E20ed8F07B840000000000000006670da1c0000b4053315fb2809A5314473E1165f6B58018E20ed8F07B840000000000000006670da1c0000b4053315fb2809A5314473E1165f6B58018E20ed8F07B8406670da0bc0866635457d36ab318d0000f3a44b7b0d08f4e198b80000000000000000000000000000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d000000000000000000000000000040"}}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var quote GetQuoteOutputFixed
			err := json.Unmarshal([]byte(tc.serializedQuoteData), &quote)
			require.NoError(t, err)

			var expectedOrdrebookOrder orderbook.Order
			err = json.Unmarshal([]byte(tc.serializedLimitOrderData), &expectedOrdrebookOrder)
			require.NoError(t, err)

			zero := big.NewInt(0)
			var expectedPreparedOrder PreparedOrder
			err = json.Unmarshal([]byte(tc.serializedPreparedOrderData), &expectedPreparedOrder)
			require.NoError(t, err)
			for _, whitelist := range expectedPreparedOrder.Order.PostInteractionData.Whitelist {
				if whitelist.Delay != nil && whitelist.Delay.Cmp(zero) == 0 {
					whitelist.Delay = zero
				}
			}

			baseSaltValue, err := BigIntFromString(tc.baseSaltValue)
			require.NoError(t, err)

			originalRandBigIntFunc := random_number_generation.BigIntMaxFunc
			first := true
			random_number_generation.BigIntMaxFunc = func(b *big.Int) (*big.Int, error) {
				if first {
					first = false
					return tc.nonce, nil
				} else {
					return baseSaltValue, nil
				}
			}

			// Monkey patch custom start time value
			originalTimeNowFunc := timeNow
			timeNow = func() int64 {
				return tc.resolverStartTime
			}

			// Monkey patch custom start time value
			originalCalcAuctionStartTimeFunc := CalcAuctionStartTimeFunc
			CalcAuctionStartTimeFunc = func(u uint32, u2 uint32) uint32 {
				return tc.auctionStartTime
			}

			wallet, err := web3_provider.DefaultWalletOnlyProvider(privateKey, tc.chainId)
			require.NoError(t, err)

			preparedOrder, orderbookOrder, err := CreateFusionOrderData(quote, tc.orderParams, wallet, tc.chainId)
			require.NoError(t, err)
			timeNow = originalTimeNowFunc
			CalcAuctionStartTimeFunc = originalCalcAuctionStartTimeFunc
			random_number_generation.BigIntMaxFunc = originalRandBigIntFunc

			assert.Equal(t, expectedOrdrebookOrder, *orderbookOrder)
			assert.Equal(t, expectedPreparedOrder, *preparedOrder)

		})
	}
}

func TestGetPreset(t *testing.T) {
	customPreset := &PresetClass{
		AllowMultipleFills: true,
		AllowPartialFills:  true,
		AuctionDuration:    10.0,
		AuctionEndAmount:   "1000",
		AuctionStartAmount: "500",
		BankFee:            "5",
		EstP:               0.1,
		ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
		GasCost: GasCostConfigClass{
			GasBumpEstimate:  1.0,
			GasPriceEstimate: "100",
		},
		InitialRateBump: 0.2,
		Points: []AuctionPointClass{
			{Coefficient: 1.0, Delay: 2.0},
		},
		StartAuctionIn: 1.0,
		TokenFee:       "1",
	}

	fastPreset := PresetClass{
		AllowMultipleFills: false,
		AllowPartialFills:  false,
		AuctionDuration:    20.0,
		AuctionEndAmount:   "2000",
		AuctionStartAmount: "1000",
		BankFee:            "10",
		EstP:               0.2,
		ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
		GasCost: GasCostConfigClass{
			GasBumpEstimate:  2.0,
			GasPriceEstimate: "200",
		},
		InitialRateBump: 0.4,
		Points: []AuctionPointClass{
			{Coefficient: 2.0, Delay: 4.0},
		},
		StartAuctionIn: 2.0,
		TokenFee:       "2",
	}

	mediumPreset := PresetClass{
		AllowMultipleFills: true,
		AllowPartialFills:  false,
		AuctionDuration:    30.0,
		AuctionEndAmount:   "3000",
		AuctionStartAmount: "1500",
		BankFee:            "15",
		EstP:               0.3,
		ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
		GasCost: GasCostConfigClass{
			GasBumpEstimate:  3.0,
			GasPriceEstimate: "300",
		},
		InitialRateBump: 0.6,
		Points: []AuctionPointClass{
			{Coefficient: 3.0, Delay: 6.0},
		},
		StartAuctionIn: 3.0,
		TokenFee:       "3",
	}

	slowPreset := PresetClass{
		AllowMultipleFills: false,
		AllowPartialFills:  true,
		AuctionDuration:    40.0,
		AuctionEndAmount:   "4000",
		AuctionStartAmount: "2000",
		BankFee:            "20",
		EstP:               0.4,
		ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
		GasCost: GasCostConfigClass{
			GasBumpEstimate:  4.0,
			GasPriceEstimate: "400",
		},
		InitialRateBump: 0.8,
		Points: []AuctionPointClass{
			{Coefficient: 4.0, Delay: 8.0},
		},
		StartAuctionIn: 4.0,
		TokenFee:       "4",
	}

	presets := QuotePresetsClass{
		Custom: customPreset,
		Fast:   fastPreset,
		Medium: mediumPreset,
		Slow:   slowPreset,
	}

	tests := []struct {
		name       string
		presetType GetQuoteOutputRecommendedPreset
		expected   *PresetClass
		expectErr  bool
	}{
		{
			name:       "Get Custom Preset",
			presetType: Custom,
			expected:   customPreset,
			expectErr:  false,
		},
		{
			name:       "Get Fast Preset",
			presetType: Fast,
			expected:   &fastPreset,
			expectErr:  false,
		},
		{
			name:       "Get Medium Preset",
			presetType: Medium,
			expected:   &mediumPreset,
			expectErr:  false,
		},
		{
			name:       "Get Slow Preset",
			presetType: Slow,
			expected:   &slowPreset,
			expectErr:  false,
		},
		{
			name:       "Unknown Preset Type",
			presetType: GetQuoteOutputRecommendedPreset("Unknown"),
			expected:   nil,
			expectErr:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getPreset(presets, tc.presetType)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateAuctionDetails(t *testing.T) {
	tests := []struct {
		name                 string
		preset               *PresetClass
		additionalWaitPeriod float32
		expected             *AuctionDetails
		expectErr            bool
	}{
		{
			name: "Valid Preset",
			preset: &PresetClass{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    60.0,
				AuctionEndAmount:   "1000",
				AuctionStartAmount: "500",
				BankFee:            "5",
				EstP:               0.1,
				ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
				GasCost: GasCostConfigClass{
					GasBumpEstimate:  1.0,
					GasPriceEstimate: "100",
				},
				InitialRateBump: 2,
				Points: []AuctionPointClass{
					{Coefficient: 1.0, Delay: 2.0},
				},
				StartAuctionIn: 5.0,
				TokenFee:       "1",
			},
			additionalWaitPeriod: 10.0,
			expected: &AuctionDetails{
				StartTime:       CalcAuctionStartTimeFunc(5, 10),
				Duration:        60,
				InitialRateBump: 2,
				Points: []AuctionPointClassFixed{
					{Coefficient: 1, Delay: 2},
				},
				GasCost: GasCostConfigClassFixed{
					GasBumpEstimate:  1,
					GasPriceEstimate: 100,
				},
			},
			expectErr: false,
		},
		{
			name: "Invalid Gas Price Estimate",
			preset: &PresetClass{
				AllowMultipleFills: true,
				AllowPartialFills:  true,
				AuctionDuration:    60.0,
				AuctionEndAmount:   "1000",
				AuctionStartAmount: "500",
				BankFee:            "5",
				EstP:               0.1,
				ExclusiveResolver:  map[string]interface{}{"resolver": "value"},
				GasCost: GasCostConfigClass{
					GasBumpEstimate:  1.0,
					GasPriceEstimate: "invalid",
				},
				InitialRateBump: 0.2,
				Points: []AuctionPointClass{
					{Coefficient: 1.0, Delay: 2.0},
				},
				StartAuctionIn: 5.0,
				TokenFee:       "1",
			},
			additionalWaitPeriod: 10.0,
			expected:             nil,
			expectErr:            true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateAuctionDetails(tc.preset, tc.additionalWaitPeriod)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestBpsToRatioFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected *big.Int
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: big.NewInt(0),
		},
		{
			name:     "Zero input",
			input:    big.NewInt(0),
			expected: big.NewInt(0),
		},
		{
			name:     "Positive input",
			input:    big.NewInt(5),
			expected: big.NewInt(50), // 5 * 100_000 / 10_000
		},
		{
			name:     "Negative input",
			input:    big.NewInt(-5),
			expected: big.NewInt(-50), // -5 * 100_000 / 10_000
		},
		{
			name:     "Large input",
			input:    big.NewInt(100_000),
			expected: big.NewInt(1_000_000), // 100_000 * 100_000 / 10_000
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := bpsToRatioFormat(tc.input)
			require.NotNil(t, result)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCreateMakerTraits(t *testing.T) {
	tests := []struct {
		name        string
		details     Details
		extraParams ExtraParams
		expected    *orderbook.MakerTraits
		expectErr   bool
		expectedErr error
	}{
		{
			name: "Valid Maker Traits",
			details: Details{
				Auction: &AuctionDetails{
					StartTime: 1000,
					Duration:  2000,
				},
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.Address{},
					},
					BankFee: big.NewInt(200),
				},
			},
			extraParams: ExtraParams{
				Nonce:                big.NewInt(1),
				Permit:               "permit",
				AllowPartialFills:    true,
				AllowMultipleFills:   true,
				OrderExpirationDelay: 3000,
				EnablePermit2:        true,
				Source:               "source",
				unwrapWeth:           true,
			},
			expected: &orderbook.MakerTraits{
				AllowedSender:       "",
				Expiry:              6000,
				Nonce:               1,
				Series:              0,
				NoPartialFills:      false,
				NeedPostinteraction: true,
				NeedPreinteraction:  false,
				NeedEpochCheck:      false,
				HasExtension:        true,
				ShouldUsePermit2:    true,
				ShouldUnwrapWeth:    true,
				AllowPartialFills:   true,
				AllowMultipleFills:  true,
			},
			expectErr: false,
		},
		{
			name: "Invalid Maker Traits - No Nonce",
			details: Details{
				Auction: &AuctionDetails{
					StartTime: 1000,
					Duration:  2000,
				},
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.Address{},
					},
					BankFee: big.NewInt(200),
				},
			},
			extraParams: ExtraParams{
				Nonce:                big.NewInt(0),
				Permit:               "permit",
				AllowPartialFills:    false,
				AllowMultipleFills:   false,
				OrderExpirationDelay: 3000,
				EnablePermit2:        true,
				Source:               "source",
				unwrapWeth:           true,
			},
			expected:    nil,
			expectErr:   true,
			expectedErr: errors.New("nonce required, when partial fill or multiple fill disallowed"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateMakerTraits(tc.details, tc.extraParams)
			if tc.expectErr {
				require.Error(t, err)
				require.Equal(t, tc.expectedErr, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateSettlementPostInteractionData(t *testing.T) {
	tests := []struct {
		name        string
		details     Details
		orderInfo   FusionOrderV4
		expected    *SettlementPostInteractionData
		expectErr   bool
		expectedErr error
	}{
		{
			name: "Valid Details and Order Info with Resolving Start Time",
			details: Details{
				ResolvingStartTime: big.NewInt(1622548800), // Example timestamp
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
						AllowFrom: big.NewInt(1622548800),
					},
				},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expected: &SettlementPostInteractionData{
				Whitelist: []WhitelistItem{
					{
						AddressHalf: "00000000000000000002",
						Delay:       big.NewInt(0),
					},
				},
				IntegratorFee: &IntegratorFee{
					Ratio:    big.NewInt(100),
					Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				},
				BankFee:            big.NewInt(200),
				ResolvingStartTime: big.NewInt(1622548800),
				CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
			},
			expectErr: false,
		},
		{
			name: "Valid Details and Order Info with non-zero Delay",
			details: Details{
				ResolvingStartTime: big.NewInt(1622548800), // Example timestamp
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
						AllowFrom: big.NewInt(1622549800),
					},
				},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expected: &SettlementPostInteractionData{
				Whitelist: []WhitelistItem{
					{
						AddressHalf: "00000000000000000002",
						Delay:       big.NewInt(1000),
					},
				},
				IntegratorFee: &IntegratorFee{
					Ratio:    big.NewInt(100),
					Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				},
				BankFee:            big.NewInt(200),
				ResolvingStartTime: big.NewInt(1622548800),
				CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
			},
			expectErr: false,
		},
		{
			name: "Valid Details and Order Info without Resolving Start Time",
			details: Details{
				ResolvingStartTime: nil,
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
						AllowFrom: big.NewInt(1622548800),
					},
				},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expected: &SettlementPostInteractionData{
				Whitelist: []WhitelistItem{
					{
						AddressHalf: "00000000000000000002",
						Delay:       big.NewInt(0),
					},
				},
				IntegratorFee: &IntegratorFee{
					Ratio:    big.NewInt(100),
					Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				},
				BankFee:            big.NewInt(200),
				ResolvingStartTime: big.NewInt(timeNow()), // This will be dynamically set
				CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
			},
			expectErr: false,
		},
		{
			name: "Delay too large",
			details: Details{
				ResolvingStartTime: big.NewInt(1622548800), // Example timestamp
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{
					{
						Address:   common.HexToAddress("0x0000000000000000000000000000000000000002"),
						AllowFrom: big.NewInt(1622649800),
					},
				},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expected: &SettlementPostInteractionData{
				Whitelist: []WhitelistItem{
					{
						AddressHalf: "00000000000000000002",
						Delay:       big.NewInt(1000),
					},
				},
				IntegratorFee: &IntegratorFee{
					Ratio:    big.NewInt(100),
					Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				},
				BankFee:            big.NewInt(200),
				ResolvingStartTime: big.NewInt(1622548800),
				CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
			},
			expectErr:   true,
			expectedErr: fmt.Errorf("delay too big - %d must be less than %d", 101000, uint16Max),
		},
		{
			name: "Whitelist empty",
			details: Details{
				ResolvingStartTime: nil,
				Fees: Fees{
					IntFee: IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000001"),
					},
					BankFee: big.NewInt(200),
				},
				Whitelist: []AuctionWhitelistItem{},
			},
			orderInfo: FusionOrderV4{
				Receiver: "0x0000000000000000000000000000000000000003",
			},
			expectErr:   true,
			expectedErr: errors.New("whitelist cannot be empty"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateSettlementPostInteractionData(tc.details, tc.orderInfo)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				// Setting the dynamic field to the expected result for comparison
				if tc.details.ResolvingStartTime == nil {
					tc.expected.ResolvingStartTime = result.ResolvingStartTime
				}
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateExtension(t *testing.T) {
	tests := []struct {
		name      string
		params    CreateExtensionParams
		expected  *Extension
		expectErr bool
	}{
		{
			name: "Valid Parameters with Permit",
			params: CreateExtensionParams{
				settlementAddress: "0x0000000000000000000000000000000000000001",
				postInteractionData: &SettlementPostInteractionData{
					Whitelist: []WhitelistItem{
						{
							AddressHalf: "abcdef",
							Delay:       big.NewInt(1000),
						},
					},
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee:            big.NewInt(200),
					ResolvingStartTime: big.NewInt(1622548800),
					CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
				},
				orderInfo: FusionOrderV4{
					MakerAsset: "0x0000000000000000000000000000000000000004",
					Receiver:   "0x0000000000000000000000000000000000000005",
				},
				details: Details{
					Auction: &AuctionDetails{
						StartTime: 1000,
						Duration:  2000,
					},
				},
				extraParams: ExtraParams{
					Permit: "0xabcdef",
				},
			},
			expected: &Extension{
				MakingAmountData: "0x0000000000000000000000000000000000000001" + "00000000000000000003e80007d0000000",
				TakingAmountData: "0x0000000000000000000000000000000000000001" + "00000000000000000003e80007d0000000",
				PostInteraction:  "0x0000000000000000000000000000000000000001000000c800640000000000000000000000000000000000000002000000000000000000000000000000000000000360b62140abcdef03e80f",
				MakerPermit:      "0x0000000000000000000000000000000000000004abcdef",
			},
			expectErr: false,
		},
		{
			name: "Valid Parameters without Permit",
			params: CreateExtensionParams{
				settlementAddress: "0x0000000000000000000000000000000000000001",
				postInteractionData: &SettlementPostInteractionData{
					Whitelist: []WhitelistItem{
						{
							AddressHalf: "abcdef",
							Delay:       big.NewInt(1000),
						},
					},
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee:            big.NewInt(200),
					ResolvingStartTime: big.NewInt(1622548800),
					CustomReceiver:     common.HexToAddress("0x0000000000000000000000000000000000000003"),
				},
				orderInfo: FusionOrderV4{
					MakerAsset: "0x0000000000000000000000000000000000000004",
					Receiver:   "0x0000000000000000000000000000000000000005",
				},
				details: Details{
					Auction: &AuctionDetails{
						StartTime: 1000,
						Duration:  2000,
					},
				},
				extraParams: ExtraParams{},
			},
			expected: &Extension{
				MakingAmountData: "0x0000000000000000000000000000000000000001" + "00000000000000000003e80007d0000000",
				TakingAmountData: "0x0000000000000000000000000000000000000001" + "00000000000000000003e80007d0000000",
				PostInteraction:  "0x0000000000000000000000000000000000000001000000c800640000000000000000000000000000000000000002000000000000000000000000000000000000000360b62140abcdef03e80f",
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CreateExtension(tc.params)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name       string
		staticSalt string
		params     CreateOrderDataParams
		expected   *Order
		expectErr  bool
	}{
		{
			name:       "Valid Order with Integrator Fee",
			staticSalt: "180431658011416401710119735245975317914670388782711199",
			params: CreateOrderDataParams{
				settlementAddress: "0x0000000000000000000000000000000000000001",
				postInteractionData: &SettlementPostInteractionData{
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee: big.NewInt(200),
				},
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
				orderInfo: FusionOrderV4{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Receiver:     "0x0000000000000000000000000000000000000006",
				},
				details: Details{
					Auction: &AuctionDetails{
						StartTime: 1000,
						Duration:  2000,
					},
				},
				extraParams: ExtraParams{
					Nonce: big.NewInt(1),
				},
				makerTraits: &orderbook.MakerTraits{
					AllowedSender:      "0x0000000000000000000000000000000000000007",
					Expiry:             5000,
					Nonce:              1,
					AllowPartialFills:  true,
					AllowMultipleFills: true,
				},
			},
			expected: &Order{
				FusionExtension: &Extension{
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
				Inner: orderbook.OrderData{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Salt:         "1e24059a92eed490b75f51b98fbf2e143c8e9712a419f59a92eed490b75f51b98fbf2e143c8e9712a419f",
					MakerTraits:  "0x4000000000000000000000000000000001000000138800000000000000000007",
					Receiver:     "0x0000000000000000000000000000000000000001", // Address of settlementAddress because Integrator Fee is set
					Extension:    "343845d3ef4b5505456e95d059a92eed490b75f51b98fbf2e143c8e9712a419f",
				},
				SettlementExtension: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				OrderInfo: FusionOrderV4{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Receiver:     "0x0000000000000000000000000000000000000006",
				},
				AuctionDetails: &AuctionDetails{
					StartTime: 1000,
					Duration:  2000,
				},
				PostInteractionData: &SettlementPostInteractionData{
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(100),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee: big.NewInt(200),
				},
				Extra: ExtraData{
					UnwrapWETH:           false,
					Nonce:                big.NewInt(1),
					Permit:               "",
					AllowPartialFills:    false,
					AllowMultipleFills:   false,
					OrderExpirationDelay: 0,
					EnablePermit2:        false,
					Source:               "",
				},
			},
			expectErr: false,
		},
		{
			name:       "Valid Order without Integrator Fee",
			staticSalt: "180431658011416401710119735245975317914670388782711199",
			params: CreateOrderDataParams{
				settlementAddress: "0x0000000000000000000000000000000000000001",
				postInteractionData: &SettlementPostInteractionData{
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(0),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee: big.NewInt(200),
				},
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
				orderInfo: FusionOrderV4{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Receiver:     "0x0000000000000000000000000000000000000006",
				},
				details: Details{
					Auction: &AuctionDetails{
						StartTime: 1000,
						Duration:  2000,
					},
				},
				extraParams: ExtraParams{
					Nonce: big.NewInt(1),
				},
				makerTraits: &orderbook.MakerTraits{
					AllowedSender:      "0x0000000000000000000000000000000000000007",
					Expiry:             5000,
					Nonce:              1,
					AllowPartialFills:  true,
					AllowMultipleFills: true,
				},
			},
			expected: &Order{
				FusionExtension: &Extension{
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
				Inner: orderbook.OrderData{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Salt:         "1e24059a92eed490b75f51b98fbf2e143c8e9712a419f59a92eed490b75f51b98fbf2e143c8e9712a419f",
					MakerTraits:  "0x4000000000000000000000000000000001000000138800000000000000000007",
					Receiver:     "0x0000000000000000000000000000000000000006", // Address of orderInfo.Receiver because Integrator Fee is not set
					Extension:    "343845d3ef4b5505456e95d059a92eed490b75f51b98fbf2e143c8e9712a419f",
				},
				SettlementExtension: common.HexToAddress("0x0000000000000000000000000000000000000001"),
				OrderInfo: FusionOrderV4{
					Maker:        "0x0000000000000000000000000000000000000003",
					MakerAsset:   "0x0000000000000000000000000000000000000004",
					TakerAsset:   "0x0000000000000000000000000000000000000005",
					MakingAmount: "1000",
					TakingAmount: "2000",
					Receiver:     "0x0000000000000000000000000000000000000006",
				},
				AuctionDetails: &AuctionDetails{
					StartTime: 1000,
					Duration:  2000,
				},
				PostInteractionData: &SettlementPostInteractionData{
					IntegratorFee: &IntegratorFee{
						Ratio:    big.NewInt(0),
						Receiver: common.HexToAddress("0x0000000000000000000000000000000000000002"),
					},
					BankFee: big.NewInt(200),
				},
				Extra: ExtraData{
					UnwrapWETH:           false,
					Nonce:                big.NewInt(1),
					Permit:               "",
					AllowPartialFills:    false,
					AllowMultipleFills:   false,
					OrderExpirationDelay: 0,
					EnablePermit2:        false,
					Source:               "",
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			originalRandBigIntFunc := random_number_generation.BigIntMaxFunc

			staticSalt, err := BigIntFromString(tc.staticSalt)
			require.NoError(t, err)
			random_number_generation.BigIntMaxFunc = func(b *big.Int) (*big.Int, error) {
				return staticSalt, nil
			}
			result, err := CreateOrder(tc.params)
			random_number_generation.BigIntMaxFunc = originalRandBigIntFunc
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
