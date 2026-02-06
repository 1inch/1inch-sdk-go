package web3_provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (w Wallet) Nonce(ctx context.Context) (uint64, error) {
	nonce, err := w.ethClient.NonceAt(ctx, *w.address, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce: %w", err)
	}
	return nonce, nil
}

func (w Wallet) Address() common.Address {
	return *w.address
}

func (w Wallet) Balance(ctx context.Context) (*big.Int, error) {
	balance, err := w.ethClient.BalanceAt(ctx, *w.address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve balance: %w", err)
	}
	return balance, nil
}

func (w Wallet) ChainId() int64 {
	return w.chainId.Int64()
}
