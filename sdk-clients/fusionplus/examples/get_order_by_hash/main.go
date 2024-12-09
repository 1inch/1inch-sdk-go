package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/sdk-clients/fusionplus"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
)

func main() {
	config, err := fusionplus.NewConfiguration(fusionplus.ConfigurationParams{
		ApiUrl:     "https://api.1inch.dev",
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

	response, err := client.GetReadyToAcceptFills(ctx, fusionplus.GetOrderByOrderHashParams{
		Hash: "0x97729858044d3838c82f2ea5ca4764bd20bfdf1f99d3af05786e4a358b16fa91",
	})
	if err != nil {
		log.Fatalf("failed to request: %v", err)
	}

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v\n", err)
	}
	fmt.Printf("Response: %s\n", string(output))
}
