package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/sdk-clients/fusion"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	publicAddress  = os.Getenv("WALLET_ADDRESS")
	privateKey     = os.Getenv("WALLET_KEY")
)

const (
	usdc    = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	wmatic  = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	weth    = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
	amount  = "200000000000000"
	chainId = 137
)

func main() {
	config, err := fusion.NewConfiguration(fusion.ConfigurationParams{
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
		ChainId:    chainId,
		PrivateKey: privateKey,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := fusion.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	fromToken := weth
	toToken := usdc

	response, err := client.GetQuote(ctx, fusion.QuoterControllerGetQuoteParamsFixed{
		FromTokenAddress: fromToken,
		ToTokenAddress:   toToken,
		Amount:           amount,
		WalletAddress:    publicAddress,
		EnableEstimate:   true,
		Surplus:          true,
	})
	if err != nil {
		log.Fatalf("failed to request: %v", err)
	}

	output, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v\n", err)
	}
	fmt.Printf("Response: %s\n", string(output))

	orderParams := fusion.OrderParams{
		WalletAddress:    publicAddress,
		FromTokenAddress: fromToken,
		ToTokenAddress:   toToken,
		Amount:           amount,
		Receiver:         "0x0000000000000000000000000000000000000000",
		Preset:           fusion.Fast,
	}

	orderHash, err := client.PlaceOrder(ctx, *response, orderParams, client.Wallet)
	if err != nil {
		log.Fatalf("failed to place order: %v", err)
	}

	fmt.Printf("Order placed! Order hash: %s\n", orderHash)
	fmt.Println("Monitoring order until it completes...")

	for {
		select {
		case <-time.After(1 * time.Second):
			order, err := client.GetOrderStatus(ctx, orderHash)
			if err != nil {
				fmt.Printf("failed to get order from order hash: %v", err)
				return
			}

			fmt.Printf("Order status: %s\n", order.Status)
			if order.Status == "filled" {
				return
			}
		}
	}
}
