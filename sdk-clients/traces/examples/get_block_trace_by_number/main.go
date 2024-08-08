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

	blockTrace, err := client.GetBlockTraceByNumber(ctx, traces.GetBlockTraceByNumberParam(17378176))
	if err != nil {
		log.Fatalf("failed to GetBlockTraceByNumber: %v", err)
	}

	blockTraceIndented, err := json.MarshalIndent(blockTrace, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal blockTrace: %v", err)
	}

	fmt.Printf("GetBlockTraceByNumber: %s\n", blockTraceIndented)
}
