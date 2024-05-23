package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/nft"
)

/*
This example demonstrates how to swap tokens on the EthereumChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

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

	chains, err := client.GetSupportedChains(ctx)
	if err != nil {
		log.Fatalf("failed to GetSupportedChains: %v", err)
	}

	fmt.Println("GetSupportedChains:", chains)
	time.Sleep(time.Second)

	nftsByAddress, err := client.GetNFTsByAddress(ctx, nft.GetNftsByAddressParams{
		ChainIds: []nft.GetNftsByAddressParamsChainIds{
			constants.EthereumChainId,
			constants.GnosisChainId,
		},
		Address: "0x083fc10cE7e97CaFBaE0fE332a9c4384c5f54E45",
	})
	if err != nil {
		log.Fatalf("failed to GetNftsByAddress: %v", err)
	}

	fmt.Println("GetNftsByAddressParams:", nftsByAddress)
	time.Sleep(time.Second)
}
