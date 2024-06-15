package fusion

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	publicAddress = os.Getenv("WALLET_ADDRESS")
	privateKey    = os.Getenv("WALLET_KEY")
)

const (
	usdc         = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	wmatic       = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	amount       = 100000000
	amountString = "100000000"
	chainId      = 137
)

func TestCreateOrder(t *testing.T) {
	tests := []struct {
		name                string
		serializedQuoteData string
		data                string
	}{
		{
			name:                "Encode/Decode Interaction",
			serializedQuoteData: serializedQuoteData,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var quote GetQuoteOutputFixed
			err := json.Unmarshal([]byte(tc.serializedQuoteData), &quote)
			require.NoError(t, err)

			orderParams := OrderParams{
				FromTokenAddress: usdc,
				ToTokenAddress:   wmatic,
				Amount:           amountString,
				Receiver:         "0x0000000000000000000000000000000000000000",
			}

			fusionOrderParamsData := FusionOrderParamsData{
				NetworkId: chainId,
				Preset:    Fast, // TODO currently always choosing the fast preset
				Receiver:  orderParams.Receiver,
			}

			additionalParams := AdditionalParams{
				FromAddress: publicAddress,
			}

			//quoteMarshalledIndented, err := json.MarshalIndent(quote, "", "  ")
			//require.NoError(t, err)
			//fmt.Printf("quote: %s\n", quoteMarshalledIndented)

			preparedOrder, orderbookOrder, err := CreateOrder(orderParams, quote, fusionOrderParamsData, additionalParams, privateKey)

			preparedOrderIndented, err := json.MarshalIndent(preparedOrder, "", "  ")
			require.NoError(t, err)

			fmt.Printf("preparedOrder: %s\n\n\n", preparedOrderIndented)

			orderbookOrderIndented, err := json.MarshalIndent(orderbookOrder, "", "  ")
			require.NoError(t, err)

			fmt.Printf("orderbookOrder: %s\n", orderbookOrderIndented)

		})
	}
}

const serializedQuoteData = `{
  "feeToken": "0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270",
  "fromTokenAmount": "100000000",
  "presets": {
    "fast": {
      "allowMultipleFills": false,
      "allowPartialFills": false,
      "auctionDuration": 180,
      "auctionEndAmount": "163176110818644018694",
      "auctionStartAmount": "163996380810118260311",
      "bankFee": "0",
      "estP": 100,
      "exclusiveResolver": null,
      "gasCost": {
        "gasBumpEstimate": 0,
        "gasPriceEstimate": "0"
      },
      "initialRateBump": 50269,
      "points": [
        {
          "coefficient": 30330,
          "delay": 126
        }
      ],
      "startAuctionIn": 17,
      "tokenFee": "60252092168960658"
    },
    "medium": {
      "allowMultipleFills": true,
      "allowPartialFills": true,
      "auctionDuration": 360,
      "auctionEndAmount": "163176110818644018694",
      "auctionStartAmount": "164220699009660650244",
      "bankFee": "0",
      "estP": 100,
      "exclusiveResolver": null,
      "gasCost": {
        "gasBumpEstimate": 0,
        "gasPriceEstimate": "0"
      },
      "initialRateBump": 64016,
      "points": [
        {
          "coefficient": 53962,
          "delay": 57
        },
        {
          "coefficient": 50269,
          "delay": 6
        },
        {
          "coefficient": 11277,
          "delay": 198
        }
      ],
      "startAuctionIn": 17,
      "tokenFee": "60252092168960658"
    },
    "slow": {
      "allowMultipleFills": true,
      "allowPartialFills": true,
      "auctionDuration": 600,
      "auctionEndAmount": "163176110818644018694",
      "auctionStartAmount": "165697198048403150647",
      "bankFee": "0",
      "estP": 100,
      "exclusiveResolver": null,
      "gasCost": {
        "gasBumpEstimate": 0,
        "gasPriceEstimate": "0"
      },
      "initialRateBump": 154501,
      "points": [
        {
          "coefficient": 53962,
          "delay": 390
        },
        {
          "coefficient": 50269,
          "delay": 6
        },
        {
          "coefficient": 23684,
          "delay": 135
        }
      ],
      "startAuctionIn": 17,
      "tokenFee": "60252092168960658"
    }
  },
  "prices": {
    "usd": {
      "fromToken": "0.9997701783033165",
      "toToken": "0.60917238"
    }
  },
  "quoteId": "b7b7164b-8afa-45ff-ad7f-3c27e2a388f7",
  "recommended_preset": "fast",
  "settlementAddress": "0xfb2809a5314473e1165f6b58018e20ed8f07b840",
  "suggested": true,
  "toTokenAmount": "164056646141520582264",
  "volume": {
    "usd": {
      "fromToken": "99.977017",
      "toToken": "99.938777584847909916"
    }
  },
  "whitelist": [
    "0x46fd018b32a9315ef5b4c0866635457d36ab318d",
    "0xc1b19a08c2798c6930b8f3a44b7b0d08f4e198b8",
    "0x0000000000000000000000000000000000000000",
    "0xad3b67bca8935cb510c8d18bd45f0b94f54a968f",
    "0x0000000000000000000000000000000000000000",
    "0x0000000000000000000000000000000000000000",
    "0x62f861201db5fdc04c48c976bf098c4dba0a061d",
    "0x0000000000000000000000000000000000000000"
  ]
}`
