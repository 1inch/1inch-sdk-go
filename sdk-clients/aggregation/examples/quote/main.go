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

/*
This example demonstrates how to swap tokens on the PolygonChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    constants.PolygonChainId,
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
		Src:      "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359",
		Dst:      "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270",
		Amount:   "1000",
		From:     client.Wallet.Address().Hex(),
		Slippage: 1,
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
