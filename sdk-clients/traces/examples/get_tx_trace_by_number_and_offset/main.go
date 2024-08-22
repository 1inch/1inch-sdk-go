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
	config, err := traces.NewConfiguration(constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	client, err := traces.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	params := traces.GetTxTraceByNumberAndOffsetParams{
		BlockNumber: 17378177,
		Offset:      1,
	}

	txTraceOffset, err := client.GetTxTraceByNumberAndOffset(ctx, params)
	if err != nil {
		log.Fatalf("failed to GetTxTraceByNumberAndOffset: %v", err)
	}

	txTraceOffsetIndented, err := json.MarshalIndent(txTraceOffset, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal txTraceOffset: %v", err)
	}

	fmt.Printf("GetTxTraceByNumberAndOffset: %s\n", txTraceOffsetIndented)
}
