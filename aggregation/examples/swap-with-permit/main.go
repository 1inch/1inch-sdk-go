package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/aggregation"
	"github.com/1inch/1inch-sdk-go/aggregation/models"
	"github.com/1inch/1inch-sdk-go/constants"
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

func main() {
	config, err := aggregation.NewDefaultConfiguration(nodeUrl, privateKey, constants.PolygonChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := aggregation.NewClient(config)

	ctx := context.Background()

	amountToSwap := big.NewInt(1e17)

	allowanceData, err := client.GetApproveAllowance(ctx, models.ApproveAllowanceParams{
		ChainId: constants.PolygonChainId,
		ApproveControllerGetAllowanceParams: models.ApproveControllerGetAllowanceParams{
			TokenAddress:  PolygonFRAX,
			WalletAddress: client.Wallet.Address().Hex(),
		},
	})

	allowance := new(big.Int)
	allowance.SetString(allowanceData.Allowance, 10)

	cmp := amountToSwap.Cmp(allowance)

	var permit string

	if cmp > 0 {
		spender, err := client.GetApproveSpender(ctx, models.ApproveSpenderParams{ChainId: constants.EthereumChainId})
		if err != nil {
			panic(err)
		}
		now := time.Now()
		twoDaysLater := now.Add(time.Hour * 24 * 2)

		permitData, err := client.Wallet.GetContractDetailsForPermit(ctx, common.HexToAddress(PolygonFRAX), common.HexToAddress(spender.Address), amountToSwap, twoDaysLater.Unix())
		if err != nil {
			panic(err)
		}

		permit, err = client.Wallet.TokenPermit(*permitData)
		if err != nil {
			panic(err)
		}
	}

	swapParams := models.AggregationControllerGetSwapParams{
		Src:      PolygonFRAX,
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
		fmt.Printf("Failed to get swap data: %v\n", err)
		return
	}

	data, err := hex.DecodeString(swapData.Tx.Data[2:])
	if err != nil {
		return
	}
	value, ok := new(big.Int).SetString(swapData.Tx.Value, 10)
	if !ok {
		return
	}
	to := common.HexToAddress(swapData.Tx.To)

	builder := client.TxBuilder.New()

	tx, err := builder.SetData(data).SetTo(&to).SetGas(uint64(swapData.Tx.Gas)).SetValue(value).Build(ctx)
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
