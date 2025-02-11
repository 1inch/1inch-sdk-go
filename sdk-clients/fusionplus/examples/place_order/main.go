package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/sdk-clients/fusionplus"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
	publicAddress  = os.Getenv("WALLET_ADDRESS")
	privateKey     = os.Getenv("WALLET_KEY")
)

func main() {
	config, err := fusionplus.NewConfiguration(fusionplus.ConfigurationParams{
		ApiUrl:     "https://api.1inch.dev",
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

	srcChain := 42161
	dstChain := 8453

	srcToken := "0xaf88d065e77c8cC2239327C5EDb3A432268e5831"
	dstToken := "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"

	invert := true
	if invert {
		srcChain, dstChain = dstChain, srcChain
		srcToken, dstToken = dstToken, srcToken
	}

	quoteParams := fusionplus.QuoterControllerGetQuoteParamsFixed{
		SrcChain:        float32(srcChain),
		DstChain:        float32(dstChain),
		SrcTokenAddress: srcToken,
		DstTokenAddress: dstToken,
		Amount:          "1500000",
		WalletAddress:   publicAddress,
		EnableEstimate:  true,
	}
	quote, err := client.GetQuote(ctx, quoteParams)
	if err != nil {
		log.Fatalf("failed to get quote: %v", err)
	}

	preset, err := fusionplus.GetPreset(quote.Presets, quote.RecommendedPreset)
	if err != nil {
		log.Fatalf("Failed to get preset: %v", err)
	}
	secretsCount := preset.SecretsCount

	secrets := make([]string, int(secretsCount))
	for i := 0; i < int(secretsCount); i++ {
		randomBytes, err := fusionplus.GetRandomBytes32()
		if err != nil {
			log.Fatalf("Failed to get random bytes: %v", err)
		}
		secrets[i] = randomBytes
	}
	var secretHashes []string
	for _, secret := range secrets {
		secretHash, err := fusionplus.HashSecret(secret)
		if err != nil {
			log.Fatalf("Failed to hash secret: %v", err)
		}
		secretHashes = append(secretHashes, secretHash)
	}

	var hashLock *fusionplus.HashLock

	if secretsCount == 1 {
		hashLock, err = fusionplus.ForSingleFill(secrets[0])
		if err != nil {
			log.Fatalf("Failed to create hashlock: %v", err)
		}
	} else {
		hashLock, err = fusionplus.ForMultipleFills(secrets)
		if err != nil {
			log.Fatalf("Failed to create hashlock: %v", err)
		}
	}

	orderParams := fusionplus.OrderParams{
		HashLock:     hashLock,
		SecretHashes: secretHashes,
		Receiver:     "0x0000000000000000000000000000000000000000",
		Preset:       quote.RecommendedPreset,
	}

	orderHash, err := client.PlaceOrder(ctx, quoteParams, quote, orderParams, client.Wallet)
	if err != nil {
		log.Fatalf("Failed to create order data: %v", err)
	}

	// Get order by hash
	order, err := client.GetOrderByOrderHash(ctx, fusionplus.GetOrderByOrderHashParams{
		Hash: orderHash,
	})
	if err != nil {
		log.Fatalf("Failed to get order by hash: %v", err)
	}

	orderQuickLookIndented, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v\n", err)
	}
	fmt.Printf("Order: %s\n", string(orderQuickLookIndented))

	// Define loop parameters
	delay := 1 * time.Second // Delay between retries
	retryCount := 0          // Current retry count
	orderStatus := ""        // Current order status

	// Loop until order status is "executed" or max retries reached
	for {
		// Get order by hash
		order, err = client.GetOrderByOrderHash(ctx, fusionplus.GetOrderByOrderHashParams{
			Hash: orderHash,
		})
		if err != nil {
			log.Printf("Failed to get order by hash: %v", err)
		} else {
			// Assuming order.Status is a string. Adjust the field access as per actual response structure.
			orderStatus = string(order.Status)
			fmt.Printf("Attempt %d: Order Status: %s\n", retryCount+1, orderStatus)

			// Check if status is "executed"
			if orderStatus == "executed" {
				fmt.Println("Order has been executed.")
				break
			}

			// Check if status is "executed"
			if orderStatus == "refunded" {
				fmt.Println("Order has been refunded.")
				break
			}
		}

		// TODO fix params on this
		fills, err := client.GetReadyToAcceptFills(ctx, fusionplus.GetOrderByOrderHashParams{
			Hash: orderHash,
		})
		if err != nil {
			log.Fatalf("failed to request: %v", err)
		}

		if len(fills.Fills) > 0 {
			// TODO the secret index needs to match the index of the fill object, but I can ignore it for single-secre orders
			err = client.SubmitSecret(ctx, fusionplus.SecretInput{
				OrderHash: orderHash,
				Secret:    secrets[0],
			})
			if err != nil {
				log.Fatalf("failed to submit secret: %v", err)
			} else {
				fmt.Println("Secret submitted!")
			}
		}

		fmt.Printf("Fills: %v\n", fills)

		// Increment retry count
		retryCount++

		// Wait before next retry
		time.Sleep(delay)
	}

	orderIndented, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal response: %v\n", err)
	}
	fmt.Printf("Order: %s\n", string(orderIndented))
}
