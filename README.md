# Dev Portal Go SDK

The SDK requires a minimum version of Go `1.21`.

Check out the [release notes](https://github.com/1inch/1inch-sdk-go/blob/globally-refactored-main/CHANGELOG.md) for information about
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

*Swap API* - [[Docs](https://portal.1inch.dev/documentation/swap/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/globally-refactored-main/sdk-clients/aggregation/examples/quote/main.go)]

*Orderbook API* - [[Docs](https://portal.1inch.dev/documentation/orderbook/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/globally-refactored-main/sdk-clients/orderbook/examples/create_order/main.go)]

*Balances API* - [[Docs](https://portal.1inch.dev/documentation/balances/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/globally-refactored-main/sdk-clients/balances/examples/main.go)]

*Gas Price API* - [[Docs](https://portal.1inch.dev/documentation/gas-price/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/globally-refactored-main/sdk-clients/gasprices/examples/main.go)]

*NFT API* - [[Docs](https://portal.1inch.dev/documentation/nft/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/globally-refactored-main/sdk-clients/nft/examples/main.go)]

*Transaction Gateway API* - [[Docs](https://portal.1inch.dev/documentation/transaction/introduction) | [SDK Example](https://github.com/1inch/1inch-sdk-go/blob/globally-refactored-main/sdk-clients/txbroadcast/examples/main.go)]


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
go get github.com/1inch/1inch-sdk-go/sdk-clients/aggregation@globally-refactored-main
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

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
)

var (
	devPortalToken = "insert_your_dev_portal_token_here" // After initial testing, update this to read from your local environment using a function like os.GetEnv()
)

func main() {
	rpcUrl := "https://eth.llamarpc.com"
	randomPrivateKey := "e8f32e723decf4051aefac8e6c1a25ad146334449d2792c2b8b15d0b59c2a35f"
	
	config, err := aggregation.NewConfiguration(rpcUrl, randomPrivateKey, constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		fmt.Printf("Failed to create configuration: %v\n", err)
		return
	}
	client, err := aggregation.NewClient(config)

	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:             "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		Dst:             "0x6b175474e89094c44da98b954eedeac495271d0f",
		Amount:          "1000000",
		From:            client.Wallet.Address().Hex(),
		Slippage:        1,
		DisableEstimate: true,
	})
	if err != nil {
		fmt.Printf("Failed to get swap data: %v\n", err)
		return
	}

	output, err := json.MarshalIndent(swapData, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal swap data: %v\n", err)
		return
	}
	fmt.Printf("%s\n", string(output))
}
```

###### Compile and Execute

```sh
go run .
```

Documentation for all API calls can be found at https://portal.1inch.dev/documentation

Each folder inside the [sdk-clients directory](https://github.com/1inch/1inch-sdk-go/blob/globally-refactored-main/sdk-clients) 
will contain an SDK for one of the 1inch APIs and will also include dedicated examples.

## Getting Help

If you have questions, want to discuss the tool, or have found a bug, please open
an [issue](https://github.com/1inch/1inch-sdk/issues) here on GitHub

## Development

Please see our [SDK Developer Guide](https://github.com/1inch/1inch-sdk/blob/main/golang/DEVELOPMENT.md) if you would
like to contribute 