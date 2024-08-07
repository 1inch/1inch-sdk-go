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

	params := traces.GetTxTraceByNumberAndHashParams{
		BlockNumber:     17378177,
		TransactionHash: "0x16897e492b2e023d8f07be9e925f2c15a91000ef11a01fc71e70f75050f1e03c",
	}

	txTrace, err := client.GetTxTraceByNumberAndHash(ctx, params)
	if err != nil {
		log.Fatalf("failed to GetTxTraceByNumberAndHash: %v", err)
	}

	txTraceIndented, err := json.MarshalIndent(txTrace, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal txTrace: %v", err)
	}

	fmt.Printf("GetTxTraceByNumberAndHash: %s\n", txTraceIndented)
}
