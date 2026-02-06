package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/sdk-clients/aggregation"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	privateKey     = os.Getenv("WALLET_KEY")
	nodeUrl        = os.Getenv("NODE_URL")
	devPortalToken = os.Getenv("DEV_PORTAL_TOKEN")
)

const (
	UsdcBase   = "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"
	WethBase   = "0x4200000000000000000000000000000000000006"
	amountUsdc = "100000" // 0.1 USDC (6 decimals)
)

func main() {
	config, err := aggregation.NewConfiguration(aggregation.ConfigurationParams{
		NodeUrl:    nodeUrl,
		PrivateKey: privateKey,
		ChainId:    constants.BaseChainId,
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

	walletAddr := client.Wallet.Address().Hex()

	// Step 1: Check Allowance
	allowanceData, err := client.GetApproveAllowance(ctx, aggregation.GetAllowanceParams{
		TokenAddress:  UsdcBase,
		WalletAddress: walletAddr,
	})
	if err != nil {
		log.Fatalf("Failed to get allowance: %v", err)
	}
	allowance := new(big.Int)
	if _, ok := allowance.SetString(allowanceData.Allowance, 10); !ok {
		log.Fatalf("Failed to parse allowance: %s", allowanceData.Allowance)
	}

	amountToSwap := new(big.Int)
	if _, ok := amountToSwap.SetString(amountUsdc, 10); !ok {
		log.Fatalf("Failed to parse amount: %s", amountUsdc)
	}

	// Step 2: Approve if needed
	if allowance.Cmp(amountToSwap) < 0 {
		fmt.Println("Insufficient allowance. Approving...")
		approveData, err := client.GetApproveTransaction(ctx, aggregation.GetApproveParams{
			TokenAddress: UsdcBase,
			Amount:       amountUsdc,
		})
		if err != nil {
			log.Fatalf("Failed to get approve data: %v", err)
		}
		data, err := hexutil.Decode(approveData.Data)
		if err != nil {
			log.Fatalf("Failed to decode approve data: %v", err)
		}
		to := common.HexToAddress(approveData.To)

		tx, err := client.TxBuilder.New().SetData(data).SetTo(&to).Build(ctx)
		if err != nil {
			log.Fatalf("Failed to build approve transaction: %v", err)
		}
		signedTx, err := client.Wallet.Sign(tx)
		if err != nil {
			log.Fatalf("Failed to sign approve transaction: %v", err)
		}
		err = client.Wallet.BroadcastTransaction(ctx, signedTx)
		if err != nil {
			log.Fatalf("Failed to broadcast approve transaction: %v", err)
		}

		fmt.Printf("Approve transaction sent: https://basescan.org/tx/%s\n", signedTx.Hash().Hex())

		// Wait for approval to be mined
		for {
			receipt, _ := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
			if receipt != nil {
				fmt.Println("Approve transaction confirmed.")
				break
			}
			time.Sleep(2 * time.Second)
		}
	} else {
		fmt.Println("Sufficient allowance already present.")
	}

	// Step 3: Perform Swap
	swapData, err := client.GetSwap(ctx, aggregation.GetSwapParams{
		Src:      UsdcBase,
		Dst:      WethBase,
		Amount:   amountUsdc,
		From:     walletAddr,
		Slippage: 1, // 1% slippage
	})
	if err != nil {
		log.Fatalf("Failed to get swap data: %v", err)
	}

	tx, err := client.TxBuilder.New().
		SetData(swapData.TxNormalized.Data).
		SetTo(&swapData.TxNormalized.To).
		SetGas(swapData.TxNormalized.Gas).
		SetValue(swapData.TxNormalized.Value).
		Build(ctx)
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

	fmt.Printf("Swap transaction sent: https://basescan.org/tx/%s\n", signedTx.Hash().Hex())

	// Wait for swap transaction to be mined
	for {
		receipt, _ := client.Wallet.TransactionReceipt(ctx, signedTx.Hash())
		if receipt != nil {
			fmt.Println("Swap transaction confirmed!")
			break
		}
		time.Sleep(2 * time.Second)
	}
}
