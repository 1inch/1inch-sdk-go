package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/v4/constants"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/fusionplus"
	"github.com/1inch/1inch-sdk-go/v4/sdk-clients/orderbook"
)

/*
This example places a cross-chain fusion order bridging USDC from Arbitrum to
Base with a signed EIP-2612 permit embedded in the order, so no prior router
allowance is needed and the maker sends no transactions at all.

The permit grants the 1inch Aggregation Router exactly the trade amount and is
executed on-chain by the protocol during the fill. Signing it requires reading
the token's permit nonce, so an RPC endpoint for the source chain is required.

Requires the following environment variables:
  - DEV_PORTAL_TOKEN: 1inch Developer Portal API key
  - WALLET_KEY:       private key (64 hex chars, no 0x prefix)
  - NODE_URL:         RPC endpoint for the source chain (Arbitrum)
*/

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
)

const (
	srcChain = 42161 // Arbitrum
	dstChain = 8453  // Base

	arbitrumUsdc = "0xaf88d065e77c8cC2239327C5EDb3A432268e5831"
	baseUsdc     = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"

	amount = "1500000" // 1.5 USDC (6 decimals)
)

func main() {
	if devPortalToken == "" || privateKey == "" || nodeUrl == "" {
		log.Fatal("set DEV_PORTAL_TOKEN, WALLET_KEY, and NODE_URL to run this example")
	}
	ctx := context.Background()

	// The orderbook client is RPC-connected and used to sign the permit, which
	// requires the token's current permit nonce from the source chain
	orderbookConfig, err := orderbook.NewConfiguration(orderbook.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    srcChain,
		ApiUrl:     "https://api.1inch.com",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create orderbook configuration: %v", err)
	}
	orderbookClient, err := orderbook.NewClient(orderbookConfig)
	if err != nil {
		log.Fatalf("failed to create orderbook client: %v", err)
	}

	plusConfig, err := fusionplus.NewConfiguration(fusionplus.ConfigurationParams{
		ApiUrl:     "https://api.1inch.com",
		ApiKey:     devPortalToken,
		PrivateKey: privateKey,
	})
	if err != nil {
		log.Fatalf("failed to create fusionplus configuration: %v", err)
	}
	client, err := fusionplus.NewClient(plusConfig)
	if err != nil {
		log.Fatalf("failed to create fusionplus client: %v", err)
	}

	// The maker address must match the signing key, so it is derived from the wallet
	owner := orderbookClient.Wallet.Address()
	makingAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		log.Fatalf("invalid amount: %s", amount)
	}

	// Sign an EIP-2612 permit granting the router exactly the trade amount
	permitData, err := orderbookClient.Wallet.GetContractDetailsForPermit(
		ctx,
		gethCommon.HexToAddress(arbitrumUsdc),
		gethCommon.HexToAddress(constants.AggregationRouterV6),
		makingAmount,
		time.Now().Add(30*time.Minute).Unix(),
	)
	if err != nil {
		log.Fatalf("failed to get permit details: %v", err)
	}
	permit, err := orderbookClient.Wallet.TokenPermit(*permitData)
	if err != nil {
		log.Fatalf("failed to sign permit: %v", err)
	}
	fmt.Println("Permit signed")

	// The permit is supplied to both the quote request and the order
	quoteParams := fusionplus.QuoterControllerGetQuoteParamsFixed{
		SrcChain:        srcChain,
		DstChain:        dstChain,
		SrcTokenAddress: arbitrumUsdc,
		DstTokenAddress: baseUsdc,
		Amount:          amount,
		WalletAddress:   owner.Hex(),
		EnableEstimate:  true,
		Permit:          permit,
	}
	quote, err := client.GetQuote(ctx, quoteParams)
	if err != nil {
		log.Fatalf("failed to get quote: %v", err)
	}

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
		Permit:       permit,
	}, client.Wallet)
	if err != nil {
		log.Fatalf("failed to place order: %v", err)
	}

	fmt.Printf("Order placed: %s\n", orderHash)
	fmt.Println("Monitoring the order and submitting secrets as escrows deploy...")

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
