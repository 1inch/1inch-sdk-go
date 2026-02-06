package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := aggregation.NewConfigurationAPI(
		constants.PolygonChainId,
		"https://api.1inch.dev",
		devPortalToken,
	)
	if err != nil {
		log.Fatalf("Failed to create configuration: %v", err)
	}
	client, err := aggregation.NewClientOnlyAPI(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	ctx := context.Background()

	response, err := client.GetLiquiditySources(ctx)
	if err != nil {
		log.Fatalf("Failed to get liquidity sources: %v", err)
	}

	for _, value := range response.Protocols {
		responseIndented, err := json.MarshalIndent(value, "", "    ")
		if err != nil {
			log.Fatalf("Failed to marshal response: %v", err)
		}
		log.Printf("Liquidity source: %s\n", responseIndented)
		break // Just print the first one as an example
	}
}
