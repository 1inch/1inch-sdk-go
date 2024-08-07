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
	secondaryWalletAddress = "0x28C6c06298d514Db089934071355E5743bf21d60"

	spenderInch = "0x111111125421ca6dc452d289314280a0f8842a65"
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

	// Get balances and allowances by wallet address list.
	response, err := client.GetBalancesAndAllowancesByWalletAddressList(ctx, balances.BalancesAndAllowancesByWalletAddressListParams{
		Wallet:  secondaryWalletAddress,
		Spender: spenderInch,
	})
	if err != nil {
		log.Fatalf("Failed to get balances and allowances by wallet address list: %v\n", err)
	}

	// Output the response.
	fmt.Println("GetBalancesAndAllowancesByWalletAddressList:", response)
}
