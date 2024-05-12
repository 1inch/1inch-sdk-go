package main

import (
	"context"
	"fmt"
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
	config, err := nft.NewConfiguration("https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := nft.NewClient(config)
	if err != nil {
		return
	}
	ctx := context.Background()

	chains, err := client.GetSupportedChains(ctx)
	if err != nil {
		fmt.Println("failed to GetSupportedChains: %w", err)
		return
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
		fmt.Println("failed to GetNftsByAddressParams: %w", err)
		return
	}

	fmt.Println("GetNftsByAddressParams:", nftsByAddress)
	time.Sleep(time.Second)
}
