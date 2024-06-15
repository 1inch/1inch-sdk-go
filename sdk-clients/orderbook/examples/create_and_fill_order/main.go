package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

/*
This example demonstrates how to create an order on the Polygon network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

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

	ecdsaPrivateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf(fmt.Sprintf("error converting private key to ECDSA: %v", err))
	}
	publicKey := ecdsaPrivateKey.Public()
	publicAddress := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	expireAfter := time.Now().Add(time.Hour).Unix()

	seriesNonce, err := client.GetSeriesNonce(ctx, publicAddress)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get series nonce: %v", err))
	}

	makerTraits := orderbook.NewMakerTraits(orderbook.MakerTraitsParams{
		AllowedSender:      zeroAddress,
		ShouldCheckEpoch:   false,
		UsePermit2:         false,
		UnwrapWeth:         false,
		HasExtension:       false,
		HasPreInteraction:  false,
		HasPostInteraction: true,
		Expiry:             expireAfter,
		Nonce:              seriesNonce.Int64(),
		Series:             0, // TODO: Series 0 always?
	})

	createOrderResponse, err := client.CreateOrder(ctx, orderbook.CreateOrderParams{
		SeriesNonce:                    seriesNonce,
		MakerTraits:                    makerTraits,
		PrivateKey:                     privateKey,
		ExpireAfter:                    time.Now().Add(time.Hour * 10).Unix(), // TODO update the field name to have "unix" suffix
		Maker:                          publicAddress.Hex(),
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

	// Sleep to accommodate free-tier API keys
	time.Sleep(time.Second)

	getOrderResponse, err := client.GetOrdersByCreatorAddress(ctx, orderbook.GetOrdersByCreatorAddressParams{
		CreatorAddress: publicAddress.Hex(),
	})

	fmt.Printf("Order created! \nOrder hash: %v\n", getOrderResponse[0].OrderHash)

	// Sleep to accommodate free-tier API keys
	time.Sleep(time.Second)

	getOrderRresponse, err := client.GetOrder(ctx, orderbook.GetOrderParams{
		OrderHash: getOrderResponse[0].OrderHash,
	})

	fillOrderData, err := client.GetFillOrderCalldata(getOrderRresponse, nil)

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
