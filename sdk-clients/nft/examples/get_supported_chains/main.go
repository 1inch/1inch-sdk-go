package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/sdk-clients/nft"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := nft.NewConfiguration(nft.ConfigurationParams{
		ApiUrl: "https://api.1inch.dev",
		ApiKey: devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	client, err := nft.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	chains, err := client.GetSupportedChains(ctx)
	if err != nil {
		log.Fatalf("failed to GetSupportedChains: %v", err)
	}

	chainsIndented, err := json.MarshalIndent(chains, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal chains: %v", err)
	}

	fmt.Printf("GetSupportedChains: %s\n", chainsIndented)
}
