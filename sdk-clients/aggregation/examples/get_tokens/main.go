package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/aggregation"
)

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

	response, err := client.GetTokens(ctx)
	if err != nil {
		log.Fatalf("failed to get tokens: %v", err)
	}

	for _, value := range response.Tokens {
		responseIndented, err := json.MarshalIndent(value, "", "    ")
		if err != nil {
			log.Fatalf("failed to marshal response: %v", err)
		}
		log.Printf("Token: %s\n", responseIndented)
		break // Just print the first one as an example
	}
}
