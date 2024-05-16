package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/tokens"
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
	config, err := tokens.NewConfiguration(tokens.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := tokens.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	whitelistedTokens, err := client.WhitelistedTokens(ctx, tokens.TokenListControllerTokensParams{})
	if err != nil {
		log.Fatalf("failed to search token: %v", err)
	}

	jsonTokens, err := json.MarshalIndent(whitelistedTokens, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal tokens: %v", err)
	}

	fmt.Println("Tokens:", string(jsonTokens))
}
