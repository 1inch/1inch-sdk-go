package main

import (
	"context"
	"encoding/json"
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
	spender           = "0x58b6a8a3302369daec383334672404ee733ab239"
)

func main() {
	config, err := balances.NewConfiguration(balances.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}

	client, err := balances.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	ctx := context.Background()

	response, err := client.GetAllowancesByWalletAddress(ctx, balances.AllowancesByWalletAddressParams{
		Wallet:  mainWalletAddress,
		Spender: spender,
	})
	if err != nil {
		log.Fatalf("Failed to get allowances by wallet address: %v\n", err)
	}

	responseIndented, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	fmt.Printf("GetAllowancesByWalletAddress: %s\n", responseIndented)
}
