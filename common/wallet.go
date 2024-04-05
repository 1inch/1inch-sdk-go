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

	GetGasTipCap(ctx context.Context) (*big.Int, error)
	GetGasPrice(ctx context.Context) (*big.Int, error)

	Sign(tx *types.Transaction) (*types.Transaction, error)
	BroadcastTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, txHash gethCommon.Hash) (*types.Receipt, error)

	TokenPermit(cd ContractPermitData) (string, error)

	IsEIP1559Applicable() bool
	ChainId() int64
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
