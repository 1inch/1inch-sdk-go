# Dev Portal Go SDK

The SDK requires a minimum version of Go `1.21`.

Check out the [release notes](https://github.com/1inch/1inch-sdk//blob/main/CHANGELOG.md) for information about the latest bug fixes, updates, and features added to the SDK.

This is a Go SDK to simplify interactions with the 1inch Dev Portal APIs. When complete, it will support all endpoints tracked by our official docs [here](https://portal.1inch.dev/documentation/authentication). See the [Current Functionality](#current-functionality) section for an up-to-date view of the SDK functionality.

Beyond mirroring the Developer Portal APIs, this SDK also supports token approvals, permit signature generation, and the execution of 1inch swaps onchain for EOA wallets.

Jump To:
* [Getting Started](#getting-started)
* [Current Functionality](#current-functionality)


## Supported APIs

*Swap API*
- [Developer Portal Docs](https://portal.1inch.dev/documentation/swap)
- [SDK Example](https://github.com/1inch/1inch-sdk/blob/main/golang/client/examples/swap/get_swap/main.go)

*Orderbook API*
- [Developer Portal Docs](https://portal.1inch.dev/documentation/orderbook)
- [SDK Example](https://github.com/1inch/1inch-sdk/blob/main/golang/client/examples/orderbook/get_orders/main.go)

## Getting started

To get started working with the SDK, set up your project for Go modules and retrieve the SDK dependencies with `go get`. This example shows how you can use the SDK to make an API request using the SDK's Swap API service:

###### Initialize Project
```
mkdir ~/hello1inch
cd ~/hello1inch
go mod init hello1inch
```

###### Add SDK Dependencies
```
go get github.com/1inch/1inch-sdk/golang
```

###### Write Code
In your preferred editor add the following content to `main.go`

**Note**: The 1inch Dev Portal Token can be generated at https://portal.1inch.dev

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk/golang/client"
	"github.com/1inch/1inch-sdk/golang/client/models"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/web3providers"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
		Web3HttpProviders: []client.Web3ProviderConfig{
			{
				ChainId: chains.Polygon,
				Url:     web3providers.Polygon,
			},
		},
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the swap request
	swapParams := models.GetSwapParams{
		ChainId:      chains.Polygon,
		SkipWarnings: false,
		AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
			Src:             tokens.PolygonFrax,
			Dst:             tokens.PolygonWeth,
			From:            os.Getenv("WALLET_ADDRESS"),
			Amount:          amounts.Ten16,
			DisableEstimate: true,
			Slippage:        0.5,
		},
	}

	swapData, _, err := c.SwapApi.GetSwap(context.Background(), swapParams)
	if err != nil {
		log.Fatalf("Failed to swap tokens: %v", err)
	}

	swapDataRawIndented, err := json.MarshalIndent(swapData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal swap data: %v", err)
	}

	fmt.Printf("%s\n", string(swapDataRawIndented))
}
```

###### Compile and Execute
```sh
go run .
```

Documentation for all API calls can be found at https://portal.1inch.dev/documentation

More example programs using the SDK can be found in the [examples directory](https://github.com/1inch/1inch-sdk/blob/main/golang/client/examples)

## Getting Help

If you have questions, want to discuss the tool, or have found a bug, please open an [issue](https://github.com/1inch/1inch-sdk/issues) here on GitHub


## Development

Please see our [SDK Developer Guide](https://github.com/1inch/1inch-sdk/blob/main/golang/DEVELOPMENT.md) if you would like to contribute 