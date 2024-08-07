package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/nft"
)

var (
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	// Initialize a new configuration using the 1inch SDK.
	config, err := nft.NewConfiguration(nft.ConfigurationParams{
		ApiUrl: "https://api.1inch.dev",
		ApiKey: devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	// Create a new client with the provided configuration.
	client, err := nft.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Create a new context for the API call.
	ctx := context.Background()

	// Define the parameters for getting NFTs by address.
	params := nft.GetNftsByAddressParams{
		ChainIds: []nft.GetNftsByAddressParamsChainIds{
			constants.EthereumChainId,
			constants.GnosisChainId,
		},
		Address: "0x083fc10cE7e97CaFBaE0fE332a9c4384c5f54E45",
	}

	// Get NFTs by address.
	nftsByAddress, err := client.GetNFTsByAddress(ctx, params)
	if err != nil {
		log.Fatalf("failed to GetNftsByAddress: %v", err)
	}

	// Marshal the response to a pretty-printed JSON format.
	nftsByAddressIndented, err := json.MarshalIndent(nftsByAddress, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal nftsByAddress: %v", err)
	}

	// Output the response.
	fmt.Printf("GetNftsByAddress: %s\n", nftsByAddressIndented)
}
