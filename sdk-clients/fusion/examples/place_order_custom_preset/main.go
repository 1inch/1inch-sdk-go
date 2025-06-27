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
	usdc       = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	wmatic     = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	weth       = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
	amountEth  = "200000000000000" // ~50 cents of ETH
	amountUSDC = "500000"          // 50 cents of USDC
	chainId    = 137
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

	fromToken := usdc
	toToken := weth
	amount := amountUSDC

	response, err := client.GetQuoteWithCustomPreset(ctx,
		fusion.QuoterControllerGetQuoteWithCustomPresetsParamsFixed{
			FromTokenAddress: fromToken,
			ToTokenAddress:   toToken,
			Amount:           amount,
			WalletAddress:    publicAddress,
			EnableEstimate:   true,
			Surplus:          true,
		},
		fusion.CustomPreset{
			AuctionDuration:    30,
			AuctionStartAmount: "240000000000000",
			AuctionEndAmount:   "150000000000000",
			Points: []fusion.CustomPresetPoint{
				{ToTokenAmount: "240000000000000", Delay: 10},
				{ToTokenAmount: "150000000000000", Delay: 20},
			},
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
		Preset:           fusion.Custom,
	}

	orderHash, err := client.PlaceOrder(ctx, *response, orderParams, client.Wallet)
	if err != nil {
		log.Fatalf("failed to place order: %v", err)
	}

	fmt.Printf("Order placed! Order hash: %s\n", orderHash)
	fmt.Println("Monitoring order until it completes...")

	auctionSecondsPassed := 0
	for {
		select {
		case <-time.After(1 * time.Second):
			order, err := client.GetOrderStatus(ctx, orderHash)
			if err != nil {
				fmt.Printf("failed to get order from order hash: %v", err)
				return
			}

			auctionSecondsPassed++
			fmt.Printf("Seconds passed: %v\n", auctionSecondsPassed)
			fmt.Printf("Order status: %s\n", order.Status)
			if order.Status == "filled" {
				return
			}
			if order.Status == "expired" {
				return
			}
		}
	}
}
