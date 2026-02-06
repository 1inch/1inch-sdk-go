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

	tokenAddress1 = "0x0d8775f648430679a709e98d2b0cb6250d2887ef"
	tokenAddress2 = "0x58b6a8a3302369daec383334672404ee733ab239"

	spender = "0x58b6a8a3302369daec383334672404ee733ab239"
)

func main() {
	config, err := balances.NewConfiguration(balances.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v", err)
	}

	client, err := balances.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	response, err := client.GetAllowancesOfCustomTokensByWalletAddress(ctx, balances.AllowancesOfCustomTokensByWalletAddressParams{
		Spender: spender,
		Wallet:  mainWalletAddress,
		Tokens:  []string{tokenAddress1, tokenAddress2},
	})
	if err != nil {
		log.Fatalf("Failed to get allowances of custom tokens by wallet address: %v", err)
	}

	responseIndented, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	fmt.Printf("GetAllowancesOfCustomTokensByWalletAddress: %s\n", responseIndented)
}
