package aggregation

import (
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
)

type SwapResponseExtended struct {
	SwapResponse

	TxNormalized NormalizedTransactionData
}

type NormalizedTransactionData struct {
	Data     []byte
	Gas      uint64
	GasPrice *big.Int
	To       gethCommon.Address
	Value    *big.Int
}

type ApproveCallDataResponseExtended struct {
	ApproveCallDataResponse

	TxNormalized NormalizedTransactionData
}
