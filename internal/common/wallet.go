package common

import (
	"context"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Wallet interface {
	Nonce(ctx context.Context) (uint64, error)
	Address() gethCommon.Address
	Balance(ctx context.Context) (*big.Int, error)

	Sign(tx *types.Transaction) (*types.Transaction, error)
	BroadcastTransaction(ctx context.Context, tx *types.Transaction) error

	// will generate the data for transaction or transaction itself
	TokenPermit()
	TokenApprove()

	// view functions
	TokenBalance()
	TokenAllowance()
}
