package main

import (
	"context"
	"encoding/hex"
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
This example demonstrates how to swap tokens on the PolygonChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	PolygonFRAX = "0x45c32fa6df82ead1e2ef74d17b76547eddfaff89"
	PolygonWeth = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
)

const (
	UniswapPermit2Polygon = "0x000000000022D473030F116dDEE9F6B43aC78BA3"
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
		log.Fatalf("Failed to create configuration: %v\n", err)
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	ctx := context.Background()

	amountToSwap := big.NewInt(1e17)

	//allowanceData, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
	//	TokenAddress:  PolygonFRAX,
	//	WalletAddress: client.Wallet.Address().Hex(),
	//})
	//if err != nil {
	//	log.Fatalf("Failed to get allowance: %v\n", err)
	//}
	//
	//allowance := new(big.Int)
	//allowance.SetString(allowanceData.Allowance, 10)

	//cmp := amountToSwap.Cmp(allowance)
	//
	//var permit string

	perimit2allowance, err := client.Wallet.TokenAllowance(ctx, PolygonFRAX, UniswapPermit2Polygon)
	if err != nil {
		log.Fatalf("Failed to get allowance: %v\n", err)
	}

	var permit string

	cmp2 := amountToSwap.Cmp(perimit2allowance)
	if cmp2 > 0 {
		approveDataString, err := client.Wallet.GenerateApproveCallData(UniswapPermit2Polygon, amountToSwap.Uint64())
		if err != nil {
			log.Fatalf("Failed to generate call data: %v\n", err)
		}

		approveData, err := hex.DecodeString(approveDataString[2:])
		if err != nil {
			log.Fatalf("Failed to decode approveDataString: %v\n", err)
		}

		toAddressFrax := common.HexToAddress(PolygonFRAX)

		builder := client.TxBuilder.New()

		tx, err := builder.SetData(approveData).SetTo(&toAddressFrax).SetGas(21_000 * 5).Build(ctx)
		if err != nil {
			log.Fatalf("Failed to build approval  transaction: %v\n", err)
		}
		signedTx, err := client.Wallet.Sign(tx)
		if err != nil {
			log.Fatalf("Failed to sign transaction: %v\n", err)
		}

		err = client.Wallet.BroadcastTransaction(ctx, signedTx)
		if err != nil {
			log.Fatalf("Failed to broadcast transaction: %v\n", err)
		}

		spender, err := client.GetApproveSpender(ctx)
		if err != nil {
			panic(err)
		}
		now := time.Now()
		twoDaysLater := now.Add(time.Hour * 24 * 2)

		permitData, err := client.Wallet.GetContractDetailsForPermit(ctx, common.HexToAddress(UniswapPermit2Polygon), common.HexToAddress(spender.Address), amountToSwap, twoDaysLater.Unix())
		if err != nil {
			log.Fatalf("Failed to get permit data: %v\n", err)
		}

		permit, err = client.Wallet.TokenPermit(*permitData)
		if err != nil {
			log.Fatalf("Failed to sign permit: %v\n", err)
		}
	}

	if cmp2 == 0 {
		spender, err := client.GetApproveSpender(ctx)
		if err != nil {
			panic(err)
		}
		now := time.Now()
		twoDaysLater := now.Add(time.Hour * 24 * 2)

		permitData, err := client.Wallet.GetContractDetailsForPermit(ctx, common.HexToAddress(UniswapPermit2Polygon), common.HexToAddress(spender.Address), amountToSwap, twoDaysLater.Unix())
		if err != nil {
			log.Fatalf("Failed to get permit data: %v\n", err)
		}

		permit, err = client.Wallet.TokenPermit(*permitData)
		if err != nil {
			log.Fatalf("Failed to sign permit: %v\n", err)
		}
	}

	swapParams := aggregation.GetSwapParams{
		Src:        PolygonFRAX,
		Dst:        PolygonWeth,
		Amount:     amountToSwap.String(),
		From:       client.Wallet.Address().Hex(),
		Slippage:   1,
		UsePermit2: true,
	}
	if permit != "" {
		swapParams.Permit = permit
	}
	swapData, err := client.GetSwap(ctx, swapParams)
	if err != nil {
		log.Fatalf("Failed to get swap data: %v\n", err)
	}

	builder := client.TxBuilder.New()

	tx, err := builder.SetData(swapData.TxNormalized.Data).SetTo(&swapData.TxNormalized.To).SetGas(swapData.TxNormalized.Gas).SetValue(swapData.TxNormalized.Value).Build(ctx)
	if err != nil {
		log.Fatalf("Failed to build transaction: %v\n", err)
	}
	signedTx, err := client.Wallet.Sign(tx)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v\n", err)
	}

	err = client.Wallet.BroadcastTransaction(ctx, signedTx)
	if err != nil {
		log.Fatalf("Failed to broadcast transaction: %v\n", err)
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
