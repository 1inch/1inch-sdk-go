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

	interval, err := client.GetSyncedInterval(ctx)
	if err != nil {
		log.Fatalf("failed to GetSyncedInterval: %v", err)
	}

	intervalIndented, err := json.MarshalIndent(interval, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal interval: %v", err)
	}

	fmt.Printf("GetSyncedInterval: %s\n", intervalIndented)
}
