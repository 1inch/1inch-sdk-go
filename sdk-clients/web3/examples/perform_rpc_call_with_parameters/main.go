package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/web3"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	config, err := web3.NewConfiguration(web3.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  "https://api.1inch.dev",
		ApiKey:  devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	client, err := web3.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	response, err := client.PerformRpcCall(ctx, web3.PerformRpcCallParams{
		PostChainIdJSONBody: web3.PostChainIdJSONBody{
			Jsonrpc: "2.0",
			Method:  "eth_getBalance",
			Params:  []string{"0xd8da6bf26964af9d7eed9e03e53415d37aa96045", "latest"},
			Id:      "1",
		},
	})
	if err != nil {
		log.Fatalf("failed to perform rpc call: %v", err)
	}

	responseMarshaled, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal response: %v", err)
	}

	fmt.Printf("Response: %s\n", responseMarshaled)
}
