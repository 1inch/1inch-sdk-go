package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/sdk-clients/fusion"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	publicAddress  = os.Getenv("WALLET_ADDRESS")
	privateKey     = os.Getenv("WALLET_KEY")
)

const (
	usdc         = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	wmatic       = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	amount       = 1000000
	amountString = "1000000"
	chainId      = 137
)

func main() {
	config, err := fusion.NewConfiguration(fusion.ConfigurationParams{
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
		ChainId: chainId,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := fusion.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	response, err := client.GetQuote(ctx, fusion.QuoterControllerGetQuoteParams{
		FromTokenAddress: usdc,
		ToTokenAddress:   wmatic,
		Amount:           amount,
		WalletAddress:    publicAddress,
		EnableEstimate:   true,
	})
	if err != nil {
		log.Fatalf("failed to request: %v", err)
	}

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v\n", err)
	}
	fmt.Printf("Response: %s\n", string(output))

	fmt.Println("Placing order")
	_, err = client.PlaceOrder(ctx, *response, fusion.OrderParams{
		FromTokenAddress: usdc,
		ToTokenAddress:   wmatic,
		Amount:           amountString,
		WalletAddress:    publicAddress,
		Receiver:         "0x0000000000000000000000000000000000000000",
	}, fusion.PlaceOrderParams{
		Maker:      publicAddress,
		PrivateKey: privateKey,
	})
	if err != nil {
		log.Fatalf("failed to request: %v", err)
	}
}
