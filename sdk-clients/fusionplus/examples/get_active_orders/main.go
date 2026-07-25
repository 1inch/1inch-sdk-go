package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusionplus"
)

/*
This example fetches cross-chain orders that are currently open for filling.

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
  - WALLET_KEY:       private key (64 hex chars, no 0x prefix)
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
)

func main() {
	if devPortalToken == "" || privateKey == "" {
		log.Fatal("set DEV_PORTAL_TOKEN and WALLET_KEY to run this example")
	}

	config, err := fusionplus.NewConfiguration(fusionplus.ConfigurationParams{
		ApiUrl:     "https://api.1inch.com",
		ApiKey:     devPortalToken,
		PrivateKey: privateKey,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := fusionplus.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	response, err := client.GetActiveOrders(ctx, fusionplus.OrderApiControllerGetActiveOrdersParams{
		Page:  1,
		Limit: 2,
	})
	if err != nil {
		log.Fatalf("failed to get active orders: %v", err)
	}

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}
	fmt.Printf("Active orders: %s\n", output)
}
