# Dev Portal Go SDK

The SDK requires a minimum version of Go `1.21`.

Check out the [release notes]() for information about the latest bug fixes, updates, and features added to the SDK.

Jump To:
* [Getting Started](#getting-started)

## Overview

This is a Go SDK to simplify interactions with the 1inch Dev Portal APIs. When complete, it will support all endpoints tracked by our official docs [here](https://portal.1inch.dev/documentation/authentication). See the `Current Functionality` section for an up-to-date view of the SDK functionality.

Beyond mirroring the Developer Portal APIs, this SDK also supports token approvals, permit signature generation, and the execution of 1inch swaps onchain for EOA wallets. 

## Current Functionality

**Supported APIs**

*Swap API*
- All endpoints supported
- Ethereum, Polygon, and Arbitrum tested (but should support all 1inch-supported chains)
- Swaps can be executed onchain from within the SDK using `Permit1` when supported and `Approve` in all other cases

*Orderbook API*
- Most endpoints supported
- Posting orders to Ethereum and Polygon is working. Other chains likely will not work at the moment

## Versioning

This library is currently in the developer preview phase (versions 0.x.x). There will be significant changes to the design of this library leading up to a 1.0.0 release. You can expect the API calls, library structure, etc. to break between each release. Once the library version reaches 1.0.0 and beyond, it will follow traditional semver conventions. 

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

## Project structure

This SDK is powered by a [client struct](https://github.com/1inch/1inch-sdk/blob/main/golang/client/client.go) that contains instances of all Services used to talk to the 1inch APIs

Each Service maps 1-to-1 with the underlying Dev Portal REST API. See [SwapService](https://github.com/1inch/1inch-sdk/blob/main/golang/client/swap.go) as an example. Under each function, you will find the matching REST API path)

Each Service uses various types and functions to do its job that are kept separate from the main service file. These can be found in the accompanying folder within the client directory (see the [swap](https://github.com/1inch/1inch-sdk/tree/main/golang/client/swap) package) 

## Issues/Suggestions

For any problems you have with the SDK or suggestions for improvements, please create an 

## Development

Please see our [SDK Developer Guide]() if you would like to contribute 