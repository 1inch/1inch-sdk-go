package web3_provider

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (w Wallet) Sign(tx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(w.chainId), w.privateKey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (w Wallet) BroadcastTransaction(ctx context.Context, tx *types.Transaction) error {
	err := w.ethClient.SendTransaction(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %v", err)
	}
	return nil
}

func (w Wallet) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return w.ethClient.TransactionReceipt(ctx, txHash)
}
