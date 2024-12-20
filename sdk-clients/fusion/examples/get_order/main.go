package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/fusion"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
	chainId        = constants.PolygonChainId
	orderhash      = "insert_order_hash_here"
)

func main() {
	config, err := fusion.NewConfiguration(fusion.ConfigurationParams{
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
		ChainId:    uint64(chainId),
		PrivateKey: privateKey,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := fusion.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	response, err := client.GetOrderStatus(ctx, orderhash)
	if err != nil {
		log.Fatalf("failed to request: %v", err)
	}

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v\n", err)
	}
	fmt.Printf("Response: %s\n", string(output))
}
