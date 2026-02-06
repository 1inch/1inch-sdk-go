package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
)

/*
This example demonstrates how to swap tokens on the Polygon network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	PolygonFrax = "0x45c32fa6df82ead1e2ef74d17b76547eddfaff89"
	PolygonWeth = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
)

func main() {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    constants.PolygonChainId,
		ApiUrl:     "https://api.1inch.dev",
		ApiKey:     devPortalToken,
	})
	if err != nil {
		log.Fatalf("Failed to create configuration: %v", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	ctx := context.Background()

	amountToSwap := big.NewInt(1e16)

	allowanceData, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
		TokenAddress:  PolygonFrax,
		WalletAddress: client.Wallet.Address().Hex(),
	})
	if err != nil {
		log.Fatalf("Failed to get allowance: %v", err)
	}

	allowance := new(big.Int)
	if _, ok := allowance.SetString(allowanceData.Allowance, 10); !ok {
		log.Fatalf("Failed to parse allowance: %s", allowanceData.Allowance)
	}

	cmp := amountToSwap.Cmp(allowance)

	var permit string

	if cmp > 0 {
		spender, err := client.GetApproveSpender(ctx)
		if err != nil {
			log.Fatalf("Failed to get approve spender: %v", err)
		}
		now := time.Now()
		twoDaysLater := now.Add(time.Hour * 24 * 2)

		permitData, err := client.Wallet.GetContractDetailsForPermit(ctx, common.HexToAddress(PolygonFrax), common.HexToAddress(spender.Address), amountToSwap, twoDaysLater.Unix())
		if err != nil {
			log.Fatalf("Failed to get permit data: %v", err)
		}

		permit, err = client.Wallet.TokenPermit(*permitData)
		if err != nil {
			log.Fatalf("Failed to sign permit: %v", err)
		}
	}

	swapParams := aggregation.GetSwapParams{
		Src:      PolygonFrax,
		Dst:      PolygonWeth,
		Amount:   amountToSwap.String(),
		From:     client.Wallet.Address().Hex(),
		Slippage: 1,
	}
	if permit != "" {
		swapParams.Permit = permit
	}
	swapData, err := client.GetSwap(ctx, swapParams)
	if err != nil {
		log.Fatalf("Failed to get swap data: %v", err)
	}

	builder := client.TxBuilder.New()

	tx, err := builder.SetData(swapData.TxNormalized.Data).SetTo(&swapData.TxNormalized.To).SetGas(swapData.TxNormalized.Gas).SetValue(swapData.TxNormalized.Value).Build(ctx)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
	}

	err = client.Wallet.BroadcastTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to broadcast transaction: %v", err)
	}

	// Waiting for transaction, just an example of it
	fmt.Printf("Transaction has been broadcast. View it on Polygonscan here: https://polygonscan.com/tx/%s\n", signedTx.Hash().Hex())
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
