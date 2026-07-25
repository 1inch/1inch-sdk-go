package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation"
)

/*
This example fetches the address tokens must be approved to before swapping
(the 1inch Aggregation Router).

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	if devPortalToken == "" {
		log.Fatal("set DEV_PORTAL_TOKEN to run this example")
	}

	config, err := aggregation.NewConfigurationAPI(
		constants.PolygonChainId,
		"https://api.1inch.com",
		devPortalToken,
	)
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := aggregation.NewClientOnlyAPI(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	response, err := client.GetApproveSpender(ctx)
	if err != nil {
		log.Fatalf("failed to get approve spender: %v", err)
	}

	fmt.Printf("Approve spender (1inch router): %s\n", response.Address)
}
