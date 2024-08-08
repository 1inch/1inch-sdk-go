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
	config, err := nft.NewConfiguration(nft.ConfigurationParams{
		ApiUrl: "https://api.1inch.dev",
		ApiKey: devPortalToken,
	})
	if err != nil {
		log.Fatalf("failed to create configuration: %v", err)
	}

	client, err := nft.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	params := nft.GetNftsByAddressParams{
		ChainIds: []nft.GetNftsByAddressParamsChainIds{
			constants.EthereumChainId,
			constants.GnosisChainId,
		},
		Address: "0x083fc10cE7e97CaFBaE0fE332a9c4384c5f54E45",
	}

	nftsByAddress, err := client.GetNFTsByAddress(ctx, params)
	if err != nil {
		log.Fatalf("failed to GetNftsByAddress: %v", err)
	}

	nftsByAddressIndented, err := json.MarshalIndent(nftsByAddress, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal nftsByAddress: %v", err)
	}

	fmt.Printf("GetNftsByAddress: %s\n", nftsByAddressIndented)
}
