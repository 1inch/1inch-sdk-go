package common

import (
	"context"
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Wallet interface {
	Call(ctx context.Context, contractAddress gethCommon.Address, callData []byte) ([]byte, error)

	Nonce(ctx context.Context) (uint64, error)
	Address() gethCommon.Address
	Balance(ctx context.Context) (*big.Int, error)

	GetGasTipCap(ctx context.Context) (*big.Int, error)
	GetGasPrice(ctx context.Context) (*big.Int, error)

	Sign(tx *types.Transaction) (*types.Transaction, error)
	BroadcastTransaction(ctx context.Context, tx *types.Transaction) error
	TransactionReceipt(ctx context.Context, txHash gethCommon.Hash) (*types.Receipt, error)

	GetContractDetailsForPermit(ctx context.Context, token gethCommon.Address, spender gethCommon.Address, amount *big.Int, deadline int64) (*ContractPermitData, error)
	GetContractDetailsForPermitDaiLike(ctx context.Context, token gethCommon.Address, spender gethCommon.Address, deadline int64) (*ContractPermitDataDaiLike, error)
	TokenPermit(cd ContractPermitData) (string, error)
	TokenPermitDaiLike(cd ContractPermitDataDaiLike) (string, error)

	GetSeriesNonce(ctx context.Context, token gethCommon.Address, publicAddress gethCommon.Address) (*big.Int, error) // TODO this should not be built into the wallet

	IsEIP1559Applicable() bool
	ChainId() int64
	//TokenApprove()

	// view functions
	//TokenBalance()
	//TokenAllowance()

}

type ContractPermitData struct {
	FromToken              string
	Spender                string
	Name                   string
	Version                string
	PublicAddress          string
	ChainId                int
	Nonce                  int64
	Deadline               int64
	Amount                 *big.Int
	IsDomainWithoutVersion bool
}

type ContractPermitDataDaiLike struct {
	FromToken              string
	Spender                string
	Name                   string
	Version                string
	Holder                 string
	ChainId                int
	Nonce                  int64
	Expiry                 int64
	Allowed                bool
	IsDomainWithoutVersion bool
}
