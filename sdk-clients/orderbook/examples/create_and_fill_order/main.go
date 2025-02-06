package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
	gethCommon "github.com/ethereum/go-ethereum/common"
)

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	wmatic      = "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"
	usdc        = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"
	ten16       = "10000000000000000"
	ten6        = "1000000"
	zeroAddress = "0x0000000000000000000000000000000000000000"
	chainId     = 137
)

func main() {
	ctx := context.Background()

	config, err := orderbook.NewConfiguration(orderbook.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    chainId,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatal(err)
	}
	client, err := orderbook.NewClient(config)

	expireAfter := time.Now().Add(time.Hour).Unix()

	seriesNonce, err := client.GetSeriesNonce(ctx, client.Wallet.Address())
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get series nonce: %v", err))
	}

	makerTraits, err := orderbook.NewMakerTraits(orderbook.MakerTraitsParams{
		AllowedSender:      zeroAddress,
		ShouldCheckEpoch:   false,
		UsePermit2:         false,
		UnwrapWeth:         false,
		HasExtension:       false,
		HasPreInteraction:  false,
		HasPostInteraction: false,
		AllowMultipleFills: true,
		AllowPartialFills:  true,
		Expiry:             expireAfter,
		Nonce:              seriesNonce.Int64(),
		Series:             0,
	})
	if err != nil {
		log.Fatalf("Failed to create maker traits: %v", err)
	}

	createOrderResponse, err := client.CreateOrder(ctx, orderbook.CreateOrderParams{
		Wallet:                         client.Wallet,
		SeriesNonce:                    seriesNonce,
		MakerTraits:                    makerTraits,
		ExpireAfterUnix:                time.Now().Add(time.Hour * 10).Unix(),
		Maker:                          client.Wallet.Address().Hex(),
		MakerAsset:                     wmatic,
		TakerAsset:                     usdc,
		MakingAmount:                   ten16,
		TakingAmount:                   ten6,
		Taker:                          zeroAddress,
		SkipWarnings:                   false,
		EnableOnchainApprovalsIfNeeded: false,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Failed to create order: %v\n", err))
	}
	if !createOrderResponse.Success {
		log.Fatalf("Request completed, but order creation status was a failure: %v\n", createOrderResponse)
	}

	fmt.Println("Order created! Getting order hash...")

	// Sleep to accommodate free-tier API keys
	time.Sleep(time.Second)

	ordersByCreatorResponse, err := client.GetOrdersByCreatorAddress(ctx, orderbook.GetOrdersByCreatorAddressParams{
		CreatorAddress: client.Wallet.Address().Hex(),
	})

	fmt.Printf("Order hash: %v\n", ordersByCreatorResponse[0].OrderHash)
	fmt.Println("Getting signature...")

	// Sleep to accommodate free-tier API keys
	time.Sleep(time.Second)

	orderWithSignature, err := client.GetOrderWithSignature(ctx, orderbook.GetOrderParams{
		OrderHash:               ordersByCreatorResponse[0].OrderHash,
		SleepBetweenSubrequests: true,
	})

	fmt.Println("Getting retrieved! Filling order...")

	fillOrderData, err := client.GetFillOrderCalldata(orderWithSignature, nil)
	if err != nil {
		log.Fatalf("Failed to get fill order calldata: %v", err)
	}

	aggregationRouter, err := constants.Get1inchRouterFromChainId(chainId)
	if err != nil {
		log.Fatalf("Failed to get 1inch router address: %v", err)
	}
	aggregationRouterAddress := gethCommon.HexToAddress(aggregationRouter)

	tx, err := client.TxBuilder.New().SetData(fillOrderData).SetTo(&aggregationRouterAddress).SetGas(150000).Build(ctx)
	if err != nil {
		fmt.Printf("Failed to build transaction: %v\n", err)
		return
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %v\n", err)
		return
	}

	err = client.Wallet.BroadcastTransaction(ctx, signedTx)
	if err != nil {
		fmt.Printf("Failed to broadcast transaction: %v\n", err)
		return
	}

	// Waiting for transaction, just an example of it
	fmt.Printf("Transaction has been broadcast. View it on Polygonscan here: %v\n", fmt.Sprintf("https://polygonscan.com/tx/%v", signedTx.Hash().Hex()))
	for {
		receipt, err := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
		if receipt != nil {
			fmt.Println("Transaction complete!")
			return
		}
		if err != nil {
			fmt.Println("Waiting for transaction to be mined")
		}
		select {
		case <-time.After(1000 * time.Millisecond): // check again after a delay
		case <-ctx.Done():
			fmt.Println("Context cancelled")
			return
		}
	}
}
