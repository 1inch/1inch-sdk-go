package main

import (
	"context"
	"encoding/hex"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
)

/*
This examples demonstrates how to swap tokens on the PolygonChainId network using the 1inch SDK.
The only thing you need to provide is your wallet address, wallet key, and dev portal token.
This can be done through your environment, or you can directly set them in the variables below
*/

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	PolygonDai  = "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063"
	PolygonWeth = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
)

func main() {
	config, err := aggregation.NewConfiguration(nodeUrl, privateKey, constants.EthereumChainId, "https://api.1inch.dev", devPortalToken)
	if err != nil {
		return
	}
	client, err := aggregation.NewClient(config)
	if err != nil {
		panic(err)
		return
	}
	ctx := context.Background()

	amountToSwap := big.NewInt(1e18)

	allowanceData, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
		TokenAddress:  PolygonDai,
		WalletAddress: client.Wallet.Address().Hex(),
	})

	allowance := new(big.Int)
	allowance.SetString(allowanceData.Allowance, 10)

	cmp := amountToSwap.Cmp(allowance)

	if cmp > 0 {
		approveData, err := client.GetApproveTransaction(ctx, aggregation.GetApproveParams{
			TokenAddress: PolygonDai,
			Amount:       amountToSwap.String(),
		})
		if err != nil {
			panic(err)
			return
		}
		data, err := hex.DecodeString(approveData.Data[2:])
		if err != nil {
			return
		}

		to := common.HexToAddress(approveData.Data)

		tx, err := client.TxBuilder.New().SetData(data).SetTo(&to).Build(ctx)
		if err != nil {
			return
		}

		signedTx, err := client.Wallet.Sign(tx)
		if err != nil {
			return
		}

		err = client.Wallet.BroadcastTransaction(ctx, signedTx)
		if err != nil {
			return
		}
	}

}
