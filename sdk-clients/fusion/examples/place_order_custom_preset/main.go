package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusion"
)

/*
This example places a gasless fusion order selling USDC for WETH on Polygon using
a custom auction preset instead of the API's fast/medium/slow presets, then
monitors it until it reaches a terminal status.

The maker must already have granted the 1inch Aggregation Router an allowance for
the sell token (see the aggregation approve example).

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
  - WALLET_KEY:       private key (64 hex chars, no 0x prefix)
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
)

const (
	usdc       = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	weth       = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
	amountUsdc = "500000" // 0.5 USDC (6 decimals)
	chainId    = 137
)

func main() {
	if devPortalToken == "" || privateKey == "" {
		log.Fatal("set DEV_PORTAL_TOKEN and WALLET_KEY to run this example")
	}

	config, err := fusion.NewConfiguration(fusion.ConfigurationParams{
		ApiUrl:     "https://api.1inch.com",
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

	// The maker address must match the signing key, so it is derived from the wallet
	owner := client.Wallet.Address().Hex()

	customPreset := fusion.CustomPreset{
		AuctionDuration:    30,
		AuctionStartAmount: "240000000000000",
		AuctionEndAmount:   "150000000000000",
		Points: []fusion.CustomPresetPoint{
			{ToTokenAmount: "240000000000000", Delay: 10},
			{ToTokenAmount: "150000000000000", Delay: 20},
		},
	}

	quote, err := client.GetQuoteWithCustomPreset(ctx,
		fusion.QuoterControllerGetQuoteWithCustomPresetsParamsFixed{
			FromTokenAddress: usdc,
			ToTokenAddress:   weth,
			Amount:           amountUsdc,
			WalletAddress:    owner,
			EnableEstimate:   true,
			Surplus:          true,
		},
		customPreset,
	)
	if err != nil {
		log.Fatalf("failed to get quote: %v", err)
	}

	quoteIndented, err := json.MarshalIndent(quote, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal quote: %v", err)
	}
	fmt.Printf("Quote: %s\n", quoteIndented)

	orderHash, err := client.PlaceOrder(ctx, *quote, fusion.OrderParams{
		WalletAddress:    owner,
		FromTokenAddress: usdc,
		ToTokenAddress:   weth,
		Amount:           amountUsdc,
		Receiver:         constants.ZeroAddress,
		Preset:           fusion.Custom,
		CustomPreset:     &customPreset,
	}, client.Wallet)
	if err != nil {
		log.Fatalf("failed to place order: %v", err)
	}

	fmt.Printf("Order placed: %s\n", orderHash)
	fmt.Println("Monitoring the order until it completes...")

	deadline := time.Now().Add(5 * time.Minute)
	for time.Now().Before(deadline) {
		time.Sleep(3 * time.Second)

		order, err := client.GetOrderStatus(ctx, orderHash)
		if err != nil {
			fmt.Printf("status poll failed, retrying: %v\n", err)
			continue
		}

		fmt.Printf("Order status: %s\n", order.Status)
		switch order.Status {
		case "filled":
			fmt.Println("Order filled")
			return
		case "expired", "cancelled", "refunded", "false-predicate", "not-enough-balance-or-allowance", "wrong-permit":
			log.Fatalf("order ended without filling (status %s)", order.Status)
		}
	}
	log.Fatalf("order %s did not reach a terminal status within 5 minutes", orderHash)
}
