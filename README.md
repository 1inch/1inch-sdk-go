# Dev Portal Go SDK

The SDK requires a minimum version of Go `1.21`.

Check out the [release notes](https://github.com/1inch/1inch-sdk-go/blob/main/CHANGELOG.md) for information about
the latest bug fixes, updates, and features added to the SDK.

This is a Go SDK to simplify interactions with the 1inch Dev Portal APIs. When complete, it will support all endpoints
tracked by our official docs [here](https://portal.1inch.dev/documentation/authentication).

Beyond mirroring the Developer Portal APIs, this SDK also supports token approvals, permit signature generation, and the
execution of 1inch swaps onchain for EOA wallets.

Jump To:

* [Supported APIs](#supported-apis)
* [Getting Started](#getting-started)
* [Getting Help](#getting-help)
* [Development](#development)

## Supported APIs

### Token Swaps
*Swap API* - [[Docs](https://portal.1inch.dev/documentation/apis/swap/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/aggregation/examples/quote/main.go)]

*Fusion API* - [~~Docs~~ | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/fusion/examples/place_order/main.go)] (Fusion does not have a docs page at this time)

*Orderbook API* - [[Docs](https://portal.1inch.dev/documentation/apis/orderbook/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/orderbook/examples/create_order/main.go)]

### Infrastructure
*Balance API* - [[Docs](https://portal.1inch.dev/documentation/apis/balance/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/balances/examples/get_allowances_of_custom_tokens/main.go)]

*Gas Price API* - [[Docs](https://portal.1inch.dev/documentation/apis/gas-price/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/gasprices/examples/get_gas_price_eip1559/main.go)]

*History API* [[Docs](https://portal.1inch.dev/documentation/apis/history/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/history/examples/get_history_events_by_address/main.go)]

*NFT API* - [[Docs](https://portal.1inch.dev/documentation/apis/nft/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/nft/examples/main.go)]

*Portfolio API* - [[Docs](https://portal.1inch.dev/documentation/apis/portfolio/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/portfolio/examples/get_current_protocols_value/main.go)]

*Spot Price API* - [[Docs](https://portal.1inch.dev/documentation/apis/spot-price/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/spotprices/examples/get_prices_for_requested_tokens/main.go)]

*Token API* - [[Docs](https://portal.1inch.dev/documentation/apis/tokens/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/tokens/examples/get_custom_token/main.go)]

*Traces API* - [[Docs](https://portal.1inch.dev/documentation/apis/traces/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/traces/examples/get_tx_trace_by_number_and_hash/main.go)]

*Transaction Gateway API* - [[Docs](https://portal.1inch.dev/documentation/apis/transaction/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients/txbroadcast/examples/broadcast_public_transaction/main.go)]


## Getting started

To get started working with the SDK, set up your project for Go modules and retrieve the SDK dependencies with `go get`.

This example shows how you can use the SDK in a new project to request a quote to swap 1 USDC for DAI on Ethereum:

###### Initialize Project

```
mkdir ~/hello1inch
cd ~/hello1inch
go mod init hello1inch
```

###### Add SDK Dependencies

```
go get github.com/1inch/1inch-sdk-go/sdk-clients/aggregation
```

###### Write Code

In your preferred editor, add the following content to `main.go` and update the `devPortalToken` variable to use your own Dev Portal Token.

**Note**: The 1inch Dev Portal Token can be generated at https://portal.1inch.dev

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
)

var (
	devPortalToken = "insert_your_dev_portal_token_here" // After initial testing, update this to read from your local environment using a function like os.GetEnv()
)

func main() {
	rpcUrl := "https://eth.llamarpc.com"
	randomPrivateKey := "e8f32e723decf4051aefac8e6c1a25ad146334449d2792c2b8b15d0b59c2a35f"

	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    rpcUrl,
		PrivateKey: randomPrivateKey,
		ChainId:    constants.EthereumChainId,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:             "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", // USDC
		Dst:             "0x111111111117dc0aa78b770fa6a738034120c302", // 1INCH
		Amount:          "100000000",
		From:            client.Wallet.Address().Hex(),
		Slippage:        1,
		DisableEstimate: true,
	})
	if err != nil {
		log.Fatalf("Failed to get swap data: %v\n", err)
	}

	output, err := json.MarshalIndent(swapData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal swap data: %v\n", err)
	}
	fmt.Printf("%s\n", string(output))
}
```

###### Compile and Execute

```sh
go run .
```

Documentation for all API calls can be found at https://portal.1inch.dev/documentation

Each folder inside the [sdk-clients directory](https://github.com/1inch/1inch-sdk-go/blob/main/sdk-clients) 
will contain an SDK for one of the 1inch APIs and will also include dedicated examples.

## Getting Help

If you have questions, want to discuss the tool, or have found a bug, please open
an [issue](https://github.com/1inch/1inch-sdk/issues) here on GitHub

## Development

Please see our [SDK Developer Guide](https://github.com/1inch/1inch-sdk-go/blob/main/DEVELOPMENT.md) if you would
like to contribute 