package web3_provider

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (w Wallet) Sign(tx *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSignerForChainID(w.chainId)
	signedTx, err := types.SignTx(tx, signer, w.privateKey)
	if err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (w Wallet) SignBytes(bytes []byte) ([]byte, error) {
	signature, err := crypto.Sign(bytes, w.privateKey)
	if err != nil {
		return nil, fmt.Errorf("error signing challenge hash: %v", err)
	}
	return signature, nil
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
