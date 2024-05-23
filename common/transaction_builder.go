package common

import (
	"context"
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TransactionBuilder interface {
	SetData(d []byte) TransactionBuilder
	SetNonce(uint64) TransactionBuilder
	SetGasPrice(*big.Int) TransactionBuilder
	SetGas(uint64) TransactionBuilder
	SetValue(*big.Int) TransactionBuilder
	SetTo(*gethCommon.Address) TransactionBuilder
	SetGasTipCap(*big.Int) TransactionBuilder
	SetGasFeeCap(*big.Int) TransactionBuilder

	BuildLegacyTx(context.Context) (*types.Transaction, error)
	BuildDynamicTx(context.Context) (*types.Transaction, error)
	Build(context.Context) (*types.Transaction, error)
}

type TransactionBuilderFactory interface {
	New() TransactionBuilder
}
