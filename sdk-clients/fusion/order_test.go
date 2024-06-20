package fusion

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

var (
	publicAddress = os.Getenv("WALLET_ADDRESS")
	privateKey    = os.Getenv("WALLET_KEY")
)

const (
	usdc         = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	wmatic       = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	amountString = "1000000000000000000"
	chainId      = 137
)

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name                        string
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
			name: "Successful order creation",
			orderParams: OrderParams{
				FromTokenAddress: wmatic,
				ToTokenAddress:   usdc,
				Amount:           amountString,
				Receiver:         "0x0000000000000000000000000000000000000000",
				Preset:           "fast",
			},
			additionalParams: AdditionalParams{
				NetworkId:   chainId,
				FromAddress: publicAddress,
				PrivateKey:  privateKey,
			},
			auctionStartTime:            1718671900,
			nonce:                       big.NewInt(887174712009),
			resolverStartTime:           1718671883,
			baseSaltValue:               "35020243109857195061155306569",
			serializedQuoteData:         `{"feeToken":"0x3c499c542cef5e3811e1192ce70d8cc03d5c3359","fromTokenAmount":"1000000000000000000","presets":{"fast":{"allowMultipleFills":false,"allowPartialFills":false,"auctionDuration":180,"auctionEndAmount":"538946","auctionStartAmount":"557310","bankFee":"0","estP":100,"exclusiveResolver":null,"gasCost":{"gasBumpEstimate":0,"gasPriceEstimate":"0"},"initialRateBump":340757,"points":[],"startAuctionIn":17,"tokenFee":"18366"},"medium":{"allowMultipleFills":true,"allowPartialFills":true,"auctionDuration":360,"auctionEndAmount":"538946","auctionStartAmount":"576251","bankFee":"0","estP":100,"exclusiveResolver":null,"gasCost":{"gasBumpEstimate":0,"gasPriceEstimate":"0"},"initialRateBump":692202,"points":[{"coefficient":681533,"delay":6},{"coefficient":340757,"delay":6}],"startAuctionIn":17,"tokenFee":"18366"},"slow":{"allowMultipleFills":true,"allowPartialFills":true,"auctionDuration":600,"auctionEndAmount":"538946","auctionStartAmount":"581432","bankFee":"0","estP":100,"exclusiveResolver":null,"gasCost":{"gasBumpEstimate":0,"gasPriceEstimate":"0"},"initialRateBump":788335,"points":[{"coefficient":681533,"delay":81},{"coefficient":340757,"delay":6}],"startAuctionIn":17,"tokenFee":"18366"}},"prices":{"usd":{"fromToken":"0.57493897","toToken":"0.9995015368854032"}},"quoteId":"55c3f478-b176-448c-b968-656c19b9c04a","recommended_preset":"fast","settlementAddress":"0xfb2809a5314473e1165f6b58018e20ed8f07b840","suggested":true,"toTokenAmount":"575677","volume":{"usd":{"fromToken":"0.57493897","toToken":"0.57539"}},"whitelist":["0x46fd018b32a9315ef5b4c0866635457d36ab318d","0xc1b19a08c2798c6930b8f3a44b7b0d08f4e198b8","0x0000000000000000000000000000000000000000","0xad3b67bca8935cb510c8d18bd45f0b94f54a968f","0x0000000000000000000000000000000000000000","0x0000000000000000000000000000000000000000","0x62f861201db5fdc04c48c976bf098c4dba0a061d","0x0000000000000000000000000000000000000000"]}`,
			serializedPreparedOrderData: `{"order":{"FusionExtension":{"MakerAssetSuffix":"","TakerAssetSuffix":"","MakingAmountData":"0xfb2809A5314473E1165f6B58018E20ed8F07B840000000000000006670da1c0000b4053315","TakingAmountData":"0xfb2809A5314473E1165f6B58018E20ed8F07B840000000000000006670da1c0000b4053315","Predicate":"","MakerPermit":"","PreInteraction":"","PostInteraction":"0xfb2809A5314473E1165f6B58018E20ed8F07B8406670da0bc0866635457d36ab318d0000f3a44b7b0d08f4e198b80000000000000000000000000000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d000000000000000000000000000040","CustomData":""},"Inner":{"makerAsset":"0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270","takerAsset":"0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359","makingAmount":"1000000000000000000","takingAmount":"538946","salt":"712810ef08aca692b6d59c49fc131590b1edc52d382c2a9684cae76e49ca45bf","maker":"0x50c5df26654B5EFBdD0c54a062dfa6012933deFe","allowedSender":"","receiver":"0x0000000000000000000000000000000000000000","makerTraits":"0x8a0000000000000000000000ce8fbbcac9006670dad000000000000000000000","extension":"357969f7ed9a797c95a9da11fc131590b1edc52d382c2a9684cae76e49ca45bf"},"SettlementExtension":"0xfb2809a5314473e1165f6b58018e20ed8f07b840","OrderInfo":{"maker":"0x50c5df26654B5EFBdD0c54a062dfa6012933deFe","makerAsset":"0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270","makerTraits":"","makingAmount":"1000000000000000000","receiver":"0x0000000000000000000000000000000000000000","salt":"","takerAsset":"0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359","takingAmount":"538946"},"AuctionDetails":{"startTime":1718671900,"duration":180,"initialRateBump":340757,"points":[],"gasCost":{"gasBumpEstimate":0,"gasPriceEstimate":0}},"PostInteractionData":{"Whitelist":[{"AddressHalf":"c0866635457d36ab318d","Delay":0},{"AddressHalf":"f3a44b7b0d08f4e198b8","Delay":0},{"AddressHalf":"00000000000000000000","Delay":0},{"AddressHalf":"d18bd45f0b94f54a968f","Delay":0},{"AddressHalf":"00000000000000000000","Delay":0},{"AddressHalf":"00000000000000000000","Delay":0},{"AddressHalf":"c976bf098c4dba0a061d","Delay":0},{"AddressHalf":"00000000000000000000","Delay":0}],"IntegratorFee":{"Ratio":0,"Receiver":"0x0000000000000000000000000000000000000000"},"BankFee":0,"ResolvingStartTime":1718671883,"CustomReceiver":"0x0000000000000000000000000000000000000000"},"Extra":{"UnwrapWETH":false,"Nonce":887174712009,"Permit":"","AllowPartialFills":false,"AllowMultipleFills":false,"OrderExpirationDelay":0,"EnablePermit2":false,"Source":""}},"hash":"0xb18805dedfb8f3deedc1d40b777bfb7dde24243651fe86202bb01d1a4e50103b","quoteId":"55c3f478-b176-448c-b968-656c19b9c04a"}`,
			serializedLimitOrderData:    `{"orderHash":"0xb18805dedfb8f3deedc1d40b777bfb7dde24243651fe86202bb01d1a4e50103b","signature":"0xcd253835b8a52dc76b58901443fa975b498c0256c4e375c3d58ae53a619f5c962f3d749b927eac0dd5407a4e8624efd40094aca62dd6b48735323aa519f4d9511c","data":{"makerAsset":"0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270","takerAsset":"0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359","makingAmount":"1000000000000000000","takingAmount":"538946","salt":"0x9a042bfb67cf14b0a1a98c4ae5d6295e2c08820","maker":"0x50c5df26654B5EFBdD0c54a062dfa6012933deFe","allowedSender":"0x0000000000000000000000000000000000000000","receiver":"0x0000000000000000000000000000000000000000","makerTraits":"0x8a0000000000000000000000ce8fbbcac9006670dad000000000000000000000","extension":"0x000000c30000004a0000004a0000004a0000004a000000250000000000000000fb2809A5314473E1165f6B58018E20ed8F07B840000000000000006670da1c0000b4053315fb2809A5314473E1165f6B58018E20ed8F07B840000000000000006670da1c0000b4053315fb2809A5314473E1165f6B58018E20ed8F07B8406670da0bc0866635457d36ab318d0000f3a44b7b0d08f4e198b80000000000000000000000000000d18bd45f0b94f54a968f0000000000000000000000000000000000000000000000000000c976bf098c4dba0a061d000000000000000000000000000040"}}`,
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

			preparedOrder, orderbookOrder, err := CreateFusionOrderData(quote, tc.orderParams, tc.additionalParams)
			timeNow = originalTimeNowFunc
			CalcAuctionStartTimeFunc = originalCalcAuctionStartTimeFunc
			random_number_generation.BigIntMaxFunc = originalRandBigIntFunc

			assert.Equal(t, expectedOrdrebookOrder, *orderbookOrder)
			assert.Equal(t, expectedPreparedOrder, *preparedOrder)

		})
	}
}
