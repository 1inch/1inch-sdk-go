package web3_provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (w Wallet) Sign(tx *types.Transaction) (*types.Transaction, error) {
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(w.ChainId), w.privateKey)
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

func (w Wallet) BuildTransaction(ctx context.Context, to *common.Address, value *big.Int, gas uint64, data []byte) (*types.Transaction, error) {
	nonce, err := w.Nonce(ctx)
	if err != nil {
		return nil, err
	}

	gasPrice, err := w.GetGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	if w.IsEIP1559Applicable() {
		gasTipCap, err := w.GetGasTipCap(ctx)
		if err != nil {
			return nil, err
		}

		return types.NewTx(&types.DynamicFeeTx{
			ChainID:   w.ChainId,
			Nonce:     nonce,
			Gas:       gas,
			To:        to,
			Value:     value,
			Data:      data,
			GasTipCap: gasTipCap,
			GasFeeCap: gasPrice,
		}), nil
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       to,
		Value:    value,
		Data:     data,
	}), nil
}
