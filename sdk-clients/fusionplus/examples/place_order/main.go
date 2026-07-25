package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusionplus"
)

/*
This example places a cross-chain fusion order bridging USDC from Arbitrum to
Base, submits the hashlock secrets as resolvers deploy escrows, and monitors the
order until it reaches a terminal status.

The maker must already have granted the 1inch Aggregation Router an allowance for
the source-chain token (see the aggregation approve example), or the order can
carry a signed permit instead (see the place_order_permit example).

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
  - WALLET_KEY:       private key (64 hex chars, no 0x prefix)
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
)

const (
	srcChain = 42161 // Arbitrum
	dstChain = 8453  // Base

	arbitrumUsdc = "0xaf88d065e77c8cC2239327C5EDb3A432268e5831"
	baseUsdc     = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"

	amount = "1500000" // 1.5 USDC (6 decimals)
)

func main() {
	if devPortalToken == "" || privateKey == "" {
		log.Fatal("set DEV_PORTAL_TOKEN and WALLET_KEY to run this example")
	}

	config, err := fusionplus.NewConfiguration(fusionplus.ConfigurationParams{
		ApiUrl:     "https://api.1inch.com",
		ApiKey:     devPortalToken,
		PrivateKey: privateKey,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := fusionplus.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	// The maker address must match the signing key, so it is derived from the wallet
	owner := client.Wallet.Address().Hex()

	quoteParams := fusionplus.QuoterControllerGetQuoteParamsFixed{
		SrcChain:        srcChain,
		DstChain:        dstChain,
		SrcTokenAddress: arbitrumUsdc,
		DstTokenAddress: baseUsdc,
		Amount:          amount,
		WalletAddress:   owner,
		EnableEstimate:  true,
	}
	quote, err := client.GetQuote(ctx, quoteParams)
	if err != nil {
		log.Fatalf("failed to get quote: %v", err)
	}

	// Each fill of the order requires revealing a secret; generate one per fill
	// and lock the order to their hashes
	preset, err := fusionplus.GetPreset(quote.Presets, quote.RecommendedPreset)
	if err != nil {
		log.Fatalf("failed to get preset: %v", err)
	}

	secrets := make([]string, int(preset.SecretsCount))
	secretHashes := make([]string, int(preset.SecretsCount))
	for i := range secrets {
		if secrets[i], err = fusionplus.GetRandomBytes32(); err != nil {
			log.Fatalf("failed to generate secret: %v", err)
		}
		if secretHashes[i], err = fusionplus.HashSecret(secrets[i]); err != nil {
			log.Fatalf("failed to hash secret: %v", err)
		}
	}

	var hashLock *fusionplus.HashLock
	if len(secrets) == 1 {
		hashLock, err = fusionplus.ForSingleFill(secrets[0])
	} else {
		hashLock, err = fusionplus.ForMultipleFills(secrets)
	}
	if err != nil {
		log.Fatalf("failed to create hashlock: %v", err)
	}

	orderHash, err := client.PlaceOrder(ctx, quoteParams, quote, fusionplus.OrderParams{
		HashLock:     hashLock,
		SecretHashes: secretHashes,
		Receiver:     constants.ZeroAddress,
		Preset:       quote.RecommendedPreset,
	}, client.Wallet)
	if err != nil {
		log.Fatalf("failed to place order: %v", err)
	}

	fmt.Printf("Order placed: %s\n", orderHash)
	fmt.Println("Monitoring the order and submitting secrets as escrows deploy...")

	monitorFusionPlusOrder(ctx, client, orderHash, secrets)
}

// monitorFusionPlusOrder polls the order status, submits a secret for each fill
// whose escrows are deployed, and returns when the order reaches a terminal status
func monitorFusionPlusOrder(ctx context.Context, client *fusionplus.Client, orderHash string, secrets []string) {
	submitted := 0
	deadline := time.Now().Add(15 * time.Minute)
	for time.Now().Before(deadline) {
		time.Sleep(5 * time.Second)

		order, err := client.GetOrderByOrderHash(ctx, fusionplus.GetOrderByOrderHashParams{Hash: orderHash})
		if err != nil {
			fmt.Printf("status poll failed, retrying: %v\n", err)
			continue
		}

		fmt.Printf("Order status: %s\n", order.Status)
		switch string(order.Status) {
		case "executed":
			fmt.Println("Order executed; funds arrive on the destination chain shortly")
			return
		case "refunded", "cancelled", "expired":
			log.Fatalf("order ended without executing (status %s)", order.Status)
		}

		fills, err := client.GetReadyToAcceptFills(ctx, fusionplus.GetReadyToAcceptFillsParams{Hash: orderHash})
		if err != nil {
			fmt.Printf("fills poll failed, retrying: %v\n", err)
			continue
		}
		for ; submitted < len(fills.Fills) && submitted < len(secrets); submitted++ {
			if err := client.SubmitSecret(ctx, fusionplus.SecretInput{
				OrderHash: orderHash,
				Secret:    secrets[submitted],
			}); err != nil {
				log.Fatalf("failed to submit secret %d: %v", submitted, err)
			}
			fmt.Printf("Submitted secret %d\n", submitted)
		}
	}
	log.Fatalf("order %s did not reach a terminal status within 15 minutes", orderHash)
}
