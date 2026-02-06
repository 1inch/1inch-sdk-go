package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/history"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := history.NewConfiguration(history.ConfigurationParams{
		ApiUrl: "https://api.1inch.dev",
		ApiKey: devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}
	client, err := history.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	ctx := context.Background()

	historyEvents, err := client.GetHistoryEventsByAddress(ctx, history.EventsByAddressParams{
		Address: "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
		ChainId: constants.EthereumChainId,
	})
	if err != nil {
		log.Fatalf("failed to GetHistoryEventsByAddress: %v", err)
	}

	historyEventsIndented, err := json.MarshalIndent(historyEvents, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshal historyEvents: %v", err)
	}

	fmt.Printf("GetHistoryEventsByAddress: %s\n", historyEventsIndented)

}
