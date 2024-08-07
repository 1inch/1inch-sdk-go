package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/balances"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	mainWalletAddress = "0x1C17622cfa9B6fD2043A76DfC39A5B5a109aa708"
)

func main() {
	// Initialize a new configuration using the 1inch SDK.
	config, err := balances.NewConfiguration(balances.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}

	// Create a new client with the provided configuration.
	client, err := balances.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	// Create a new context for the API call.
	ctx := context.Background()

	// Get balances by wallet address.
	response, err := client.GetBalancesByWalletAddress(ctx, balances.BalancesByWalletAddressParams{Wallet: mainWalletAddress})
	if err != nil {
		log.Fatalf("Failed to get balances by wallet address: %v\n", err)
	}

	// Output the response.
	fmt.Println("GetBalancesByWalletAddress:", response)
}
