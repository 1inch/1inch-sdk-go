package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/gasprices"
)

/*
This example demonstrates how to swap tokens on the EthereumChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

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
		log.Fatal("failed to create configuration: %w", err)
	}

	configLegacyChain, err := gasprices.NewConfiguration(gasprices.ConfigurationParams{
		ChainId: constants.AuroraChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatal("failed to create configuration for legacy chain: %w", err)
	}

	client, err := gasprices.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	clientLegacyChain, err := gasprices.NewClient(configLegacyChain)
	if err != nil {
		log.Fatalf("failed to create legacy client: %v", err)
	}
	ctx := context.Background()

	gasPriceEIP15559, err := client.GetGasPriceEIP1559(ctx)
	if err != nil {
		log.Fatalf("failed to GetGasPriceEIP1559: %v", err)
	}

	fmt.Println("GetGasPriceEIP1559:", gasPriceEIP15559)
	time.Sleep(time.Second)

	gasPriceLegacy, err := clientLegacyChain.GetGasPriceLegacy(ctx)
	if err != nil {
		log.Fatalf("failed to GetGasPriceLegacy: %v", err)
	}

	fmt.Println("GetGasPriceLegacy:", gasPriceLegacy)
	time.Sleep(time.Second)

}
