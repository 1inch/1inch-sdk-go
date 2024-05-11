package main

import (
	"context"
	"fmt"
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
	config, err := gasprices.NewConfiguration(constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}

	configLegacyChain, err := gasprices.NewConfiguration(constants.AuroraChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}

	client, err := gasprices.NewClient(config)
	if err != nil {
		return
	}

	clientLegacyChain, err := gasprices.NewClient(configLegacyChain)
	if err != nil {
		return
	}
	ctx := context.Background()

	gasPriceEIP15559, err := client.GetGasPriceEIP1559(ctx)
	if err != nil {
		fmt.Println("failed to GetGasPriceEIP1559: %w", err)
		return
	}

	fmt.Println("GetGasPriceEIP1559:", gasPriceEIP15559)
	time.Sleep(time.Second)

	gasPriceLegacy, err := clientLegacyChain.GetGasPriceLegacy(ctx)
	if err != nil {
		fmt.Println("failed to GetGasPriceLegacy: %w", err)
		return
	}

	fmt.Println("GetGasPriceLegacy:", gasPriceLegacy)
	time.Sleep(time.Second)

}
