package onchain

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetDynamicFeeTx(client *ethclient.Client, chainID *big.Int, fromAddress common.Address, to string, data []byte) *types.Transaction {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}

	fmt.Printf("Current nonce: %v\n", nonce)

	gasTipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas tip cap: %v", err)
	}

	gasFeeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Failed to suggest gas fee cap: %v", err)
	}

	toAddress := common.HexToAddress(to)
	value := big.NewInt(0)     // in wei (0 eth)
	gasLimit := uint64(210000) // TODO make sure this value is always correct

	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasFeeCap,
		GasTipCap: gasTipCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      data,
	})
}

func WaitForTransaction(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if receipt != nil {
			return receipt, nil
		}
		if err != nil {
			fmt.Println("Transaction not yet mined...")
		}
		select {
		case <-time.After(1 * time.Second): // check again after a delay
		case <-context.Background().Done():
			return nil, context.Background().Err()
		}
	}
}
