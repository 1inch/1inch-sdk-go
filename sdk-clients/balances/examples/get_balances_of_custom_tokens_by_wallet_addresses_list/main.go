package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/balances"
)

/*
This example fetches the balances of a specific list of tokens for multiple
wallet addresses in one request.

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	binanceHotWallet  = "0x28C6c06298d514Db089934071355E5743bf21d60"
	binanceColdWallet = "0xF977814e90dA44bFA03b6295A0616a897441aceC"

	mainnetUsdc = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
	mainnetDai  = "0x6B175474E89094C44Da98b954EedeAC495271d0F"
)

func main() {
	if devPortalToken == "" {
		log.Fatal("set DEV_PORTAL_TOKEN to run this example")
	}

	config, err := balances.NewConfiguration(balances.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.com",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := balances.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	response, err := client.GetBalancesOfCustomTokensByWalletAddressesList(ctx, balances.BalancesOfCustomTokensByWalletAddressesListParams{
		Wallets: []string{binanceHotWallet, binanceColdWallet},
		Tokens:  []string{mainnetUsdc, mainnetDai},
	})
	if err != nil {
		log.Fatalf("failed to get balances of custom tokens by wallet list: %v", err)
	}

	responseIndented, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}
	fmt.Printf("Balances by wallet: %s\n", responseIndented)
}
