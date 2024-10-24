package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
)

/*
This example demonstrates how to swap tokens on the PolygonChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

func main() {
	SwapWithCustomConnectorTokens(
		constants.EthereumChainId,
		"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
		"0x6b175474e89094c44da98b954eedeac495271d0f,0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48,0xdac17f958d2ee523a2206206994597c13d831ec7",
		"1290000000000000000000",
	)
	SwapWithCustomSlippage(
		constants.EthereumChainId,
		"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
		"1290000000000000000000",
		15,
	)

	SwapWithCustomProtocols(
		constants.EthereumChainId,
		"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
		"1290000000000000000000",
		"UNISWAP_V1,UNISWAP_V2",
	)

	SwapWithCustomReceiver(
		constants.EthereumChainId,
		"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
		"1290000000000000000000",
		"0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
	)

	SwapWithCustomFeeAndReferrer(
		constants.EthereumChainId,
		"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
		"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
		"1290000000000000000000",
		2,
		"0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
	)
}

func SwapWithCustomConnectorTokens(chain uint64, src string, dst string, connectors string, amount string) {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chain,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:             src,
		Dst:             dst,
		ConnectorTokens: connectors,
		Amount:          amount,
		From:            client.Wallet.Address().Hex(),
		Slippage:        1,
		DisableEstimate: true,
	})
	if err != nil {
		log.Fatalf("Failed to get swap data: %v\n", err)
	}

	tx, err := client.TxBuilder.New().SetData(swapData.TxNormalized.Data).SetTo(&swapData.TxNormalized.To).SetGas(swapData.TxNormalized.Gas).SetValue(swapData.TxNormalized.Value).Build(ctx)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v\n", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v\n", err)
	}
	if signedTx != nil {
		fmt.Println("Signed transaction: ", signedTx.Hash().Hex())
	}
}

func SwapWithCustomSlippage(chain uint64, src string, dst string, amount string, slippage float32) {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chain,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:             src,
		Dst:             dst,
		Amount:          amount,
		From:            client.Wallet.Address().Hex(),
		Slippage:        slippage,
		DisableEstimate: true,
	})
	if err != nil {
		log.Fatalf("Failed to get swap data: %v\n", err)
	}

	tx, err := client.TxBuilder.New().SetData(swapData.TxNormalized.Data).SetTo(&swapData.TxNormalized.To).SetGas(swapData.TxNormalized.Gas).SetValue(swapData.TxNormalized.Value).Build(ctx)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v\n", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v\n", err)
	}
	if signedTx != nil {
		fmt.Println("Signed transaction: ", signedTx.Hash().Hex())
	}
}

func SwapWithCustomProtocols(chain uint64, src string, dst string, amount string, protocols string) {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chain,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	ctx := context.Background()
	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:             src,
		Dst:             dst,
		Amount:          amount,
		From:            client.Wallet.Address().Hex(),
		Slippage:        1,
		Protocols:       protocols,
		DisableEstimate: true,
	})
	if err != nil {
		log.Fatalf("Failed to get swap data: %v\n", err)
	}

	tx, err := client.TxBuilder.New().SetData(swapData.TxNormalized.Data).SetTo(&swapData.TxNormalized.To).SetGas(swapData.TxNormalized.Gas).SetValue(swapData.TxNormalized.Value).Build(ctx)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v\n", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v\n", err)
	}
	if signedTx != nil {
		fmt.Println("Signed transaction: ", signedTx.Hash().Hex())
	}
}

func SwapWithCustomReceiver(chain uint64, src string, dst string, amount string, receiver string) {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chain,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	ctx := context.Background()
	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:             src,
		Dst:             dst,
		Amount:          amount,
		Slippage:        1,
		From:            client.Wallet.Address().Hex(),
		Receiver:        receiver,
		DisableEstimate: true,
	})
	if err != nil {
		log.Fatalf("Failed to get swap data: %v\n", err)
	}

	tx, err := client.TxBuilder.New().SetData(swapData.TxNormalized.Data).SetTo(&swapData.TxNormalized.To).SetGas(swapData.TxNormalized.Gas).SetValue(swapData.TxNormalized.Value).Build(ctx)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v\n", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v\n", err)
	}
	if signedTx != nil {
		fmt.Println("Signed transaction: ", signedTx.Hash().Hex())
	}
}

func SwapWithCustomFeeAndReferrer(chain uint64, src string, dst string, amount string, fee float32, referrer string) {
	// nodeUrl, privateKey, chain, "https://api.1inch.dev", devPortalToken
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chain,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:             src,
		Dst:             dst,
		Amount:          amount,
		Slippage:        1,
		From:            client.Wallet.Address().Hex(),
		Fee:             fee,
		Referrer:        referrer,
		DisableEstimate: true,
	})
	if err != nil {
		log.Fatalf("Failed to get swap data: %v\n", err)
	}

	tx, err := client.TxBuilder.New().SetData(swapData.TxNormalized.Data).SetTo(&swapData.TxNormalized.To).SetGas(swapData.TxNormalized.Gas).SetValue(swapData.TxNormalized.Value).Build(ctx)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v\n", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v\n", err)
	}
	if signedTx != nil {
		fmt.Println("Signed transaction: ", signedTx.Hash().Hex())
	}
}
