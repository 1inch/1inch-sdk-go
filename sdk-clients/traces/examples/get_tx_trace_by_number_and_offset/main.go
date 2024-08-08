package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/traces"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	// Initialize a new configuration using the 1inch SDK for the Ethereum chain.
	config, err := traces.NewConfiguration(constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	// Create a new client with the provided configuration.
	client, err := traces.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Create a new context for the API call.
	ctx := context.Background()

	// Define the parameters for getting the transaction trace by number and offset.
	params := traces.GetTxTraceByNumberAndOffsetParams{
		BlockNumber: 17378177,
		Offset:      1,
	}

	// Get the transaction trace by number and offset.
	txTraceOffset, err := client.GetTxTraceByNumberAndOffset(ctx, params)
	if err != nil {
		log.Fatalf("failed to GetTxTraceByNumberAndOffset: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	txTraceOffsetIndented, err := json.MarshalIndent(txTraceOffset, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal txTraceOffset: %v", err)
	}

	// Output the response.
	fmt.Printf("GetTxTraceByNumberAndOffset: %s\n", txTraceOffsetIndented)
}
