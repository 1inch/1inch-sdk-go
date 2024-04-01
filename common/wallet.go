package common

import (
	"context"
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Wallet interface {
	Nonce(ctx context.Context) (uint64, error)
	Address() gethCommon.Address
	Balance(ctx context.Context) (*big.Int, error)

	Sign(tx *types.Transaction) (*types.Transaction, error)
	BroadcastTransaction(ctx context.Context, tx *types.Transaction) error

	// will generate the data for transaction or transaction itself
	TokenPermit(cd ContractPermitData) (string, error)
	//TokenApprove()

	// view functions
	//TokenBalance()
	//TokenAllowance()
}

type ContractPermitData struct {
	FromToken     string
	Spender       string
	Name          string
	Version       string
	PublicAddress string
	ChainId       int
	Nonce         int64
	Deadline      int64
	Amount        string
}
