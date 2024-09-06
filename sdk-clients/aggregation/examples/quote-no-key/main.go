package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := aggregation.NewConfigurationAPI( // A light-weight API that is not capable of signing transactions
		constants.PolygonChainId,
		"https://api.1inch.dev",
		devPortalToken,
	)
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClientOnlyAPI(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:             "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359",
		Dst:             "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
		Amount:          "1000",
		From:            "0x0000000000000000000000000000000000000000", // Change this to any wallet address
		Slippage:        1,
		DisableEstimate: true, // This stops the 1inch API from failing if the wallet is not able to make the swap
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
