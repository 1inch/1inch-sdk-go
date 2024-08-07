package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/gasprices"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	configLegacyChain, err := gasprices.NewConfiguration(gasprices.ConfigurationParams{
		ChainId: constants.AuroraChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration for legacy chain: %v", err)
	}

	clientLegacyChain, err := gasprices.NewClient(configLegacyChain)
	if err != nil {
		log.Fatalf("failed to create legacy client: %v", err)
	}

	ctx := context.Background()

	gasPriceLegacy, err := clientLegacyChain.GetGasPriceLegacy(ctx)
	if err != nil {
		log.Fatalf("failed to GetGasPriceLegacy: %v", err)
	}

	gasPriceLegacyIndented, err := json.MarshalIndent(gasPriceLegacy, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal gasPriceLegacy: %v", err)
	}

	fmt.Printf("GetGasPriceLegacy: %s\n", gasPriceLegacyIndented)
}
