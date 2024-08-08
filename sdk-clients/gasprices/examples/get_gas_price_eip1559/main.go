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
	config, err := gasprices.NewConfiguration(gasprices.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	client, err := gasprices.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	gasPriceEIP1559, err := client.GetGasPriceEIP1559(ctx)
	if err != nil {
		log.Fatalf("failed to GetGasPriceEIP1559: %v", err)
	}

	gasPriceEIP1559Indented, err := json.MarshalIndent(gasPriceEIP1559, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal gasPriceEIP1559: %v", err)
	}

	fmt.Printf("GetGasPriceEIP1559: %s\n", gasPriceEIP1559Indented)
}
