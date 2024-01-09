# Dev Portal Go SDK

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

## Development

### Type generation

Type generation is done using the `generate_types.sh` script. To add a new swagger file or update an existing one, place the swagger file in `swagger-static` and run the script. It will generate the types file and place it in the appropriately-named folder in a sub-folder inside the `client` directory

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
