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

	whitelistedTokensList, err := client.GetCustomToken(ctx, tokens.CustomTokensControllerGetTokenInfoParams{
		Address: "0x111111111117dc0aa78b770fa6a738034120c302",
	})
	if err != nil {
		log.Fatalf("failed to search token: %v", err)
	}

	jsonTokens, err := json.MarshalIndent(whitelistedTokensList, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal tokens: %v", err)
	}

	fmt.Println("Tokens:", string(jsonTokens))
}
