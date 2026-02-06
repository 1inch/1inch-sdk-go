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
	secondaryWalletAddress = "0x28C6c06298d514Db089934071355E5743bf21d60"

	spenderInch = "0x111111125421ca6dc452d289314280a0f8842a65"
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

	response, err := client.GetBalancesAndAllowancesByWalletAddressList(ctx, balances.BalancesAndAllowancesByWalletAddressListParams{
		Wallet:  secondaryWalletAddress,
		Spender: spenderInch,
	})
	if err != nil {
		log.Fatalf("Failed to get balances and allowances by wallet address list: %v", err)
	}

	responseIndented, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	fmt.Printf("GetBalancesAndAllowancesByWalletAddressList: %s\n", responseIndented)
}
