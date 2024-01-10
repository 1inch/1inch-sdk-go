# Dev Portal Go SDK

First and foremost, it is important to note that when using the SDK libraries for the 1inch aggregator or limit order protocols, you will be creating *real* transaction data that can and will be executed onchain. Always be  deliberate when calling [CreateOrder](https://github.com/1inch/1inch-sdk/blob/9703d3bf4a6d94ad64badd0e346a16d60a0c3509/golang/client/orderbook.go#L18-L18), [SwapTokens](https://github.com/1inch/1inch-sdk/blob/9703d3bf4a6d94ad64badd0e346a16d60a0c3509/golang/actions/swap.go#L21-L21), or [GetSwapData](https://github.com/1inch/1inch-sdk/blob/9703d3bf4a6d94ad64badd0e346a16d60a0c3509/golang/client/swap.go#L128-L128). When filling out the parameters for these functions, make sure you understand concepts like [slippage](https://medium.com/onomy-protocol/what-is-slippage-in-defi-62a0d068feb3) and [MEV](https://chain.link/education-hub/maximal-extractable-value-mev), as well as the difference between [USDC](https://etherscan.io/token/0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48) having 6 digits of precision on Ethereum and [DAI](https://etherscan.io/token/0x6b175474e89094c44da98b954eedeac495271d0f) having 18! 

## Overview

This is a Go SDK to simplify interactions with the 1inch Dev Portal APIs. It will support all endpoints tracked by our official docs [here](https://portal.1inch.dev/documentation/authentication).

Additionally, this SDK also supports executing 1inch swaps onchain for your wallet. 

## Using the SDK in your project

The SDK can be used by first creating a config object, calling the constructor, then accessing the service for the API of interest. Here is a simple program using the SDK that will generate swap data using the 1inch Aggregator:

**Note**: A 1inch Dev Portal Token can be generated at [portal.1inch.dev](https://portal.1inch.dev)  

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/1inch/1inch-sdk/golang/client"
	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
)

func main() {

	// Build the config for the client
	config := client.Config{
		DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
		ChainId:         chains.Polygon,
	}

	// Create the 1inch client
	c, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Build the config for the swap request
	swapParams := swap.AggregationControllerGetSwapParams{
		Src:             tokens.PolygonFrax,
		Dst:             tokens.PolygonWeth,
		From:            os.Getenv("WALLET_ADDRESS"),
		Amount:          amounts.Ten16,
		DisableEstimate: helpers.GetPtr(true),
	}

	swapData, _, err := c.Swap.GetSwapData(context.Background(), swapParams)
	if err != nil {
		log.Fatalf("Failed to get swap data: %v", err)
	}

	helpers.PrettyPrintStruct(swapData)
}
```

More example programs using the SDK can be found in the [examples directory]()

## Project structure

This SDK is powered by a [client struct](https://github.com/1inch/1inch-sdk/blob/9703d3bf4a6d94ad64badd0e346a16d60a0c3509/golang/client/client.go#L60-L60) that contains instances of all Services used to talk to the 1inch APIs

Each Service is simply a struct that contains all endpoints from a given 1inch API (see [SwapService](https://github.com/1inch/1inch-sdk/blob/9703d3bf4a6d94ad64badd0e346a16d60a0c3509/golang/client/swap.go#L11-L11))

Each Service uses various types and functions to do its job that are kept separate from the main service file. These can be found in the accompanying folder within the client directory (see the [swap](https://github.com/1inch/1inch-sdk/tree/9703d3bf4a6d94ad64badd0e346a16d60a0c3509/golang/client/swap) package) 

## Issues/Suggestions

For any problems you have with the SDK or suggestions for improvements, please create an [issue](https://github.com/1inch/1inch-sdk/issues) here on GitHub

## Development

### Type generation

Type generation is done using the `generate_types.sh` script. To add a new swagger file or update an existing one, place the swagger file in `swagger-static` and run the script. It will generate the types file and place it in the appropriately-named sub-folder inside the `client` directory

### Swagger file formatting
For consistency, Swagger files should be formatted with `prettier`

This can be installed globally using npm:

`npm install -g prettier`

If using GoLand, you can setup this action to run automatically using File Watchers:

1. Go to Settings or Preferences > Tools > File Watchers.
2. Click the + button to add a new watcher.
3. For `File type`, choose JSON.
4. For `Scope`, choose Project Files.
5. For `Program`, provide the path to the `prettier`. This can be gotten by running `which prettier`.
6. For `Arguments`, use `--write $FilePath$`.
7. For `Output paths to refresh`, use `$FilePath$`.
8. Ensure the Auto-save edited files to trigger the watcher option is checked
