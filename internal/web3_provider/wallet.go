package web3_provider

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func (w Wallet) Nonce(ctx context.Context) (uint64, error) {
	nonce, err := w.ethClient.NonceAt(ctx, w.address, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce: %v", err)
	}
	return nonce, nil
}

func (w Wallet) Address() common.Address {
	publicKey := w.privateKey.Public()

	return crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))
}

func (w Wallet) Balance(ctx context.Context) (*big.Int, error) {
	balance, err := w.ethClient.BalanceAt(ctx, w.address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve balance: %v", err)
	}

	return balance, nil
}

func (w Wallet) Sign(tx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(w.chainID), w.privateKey)
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
