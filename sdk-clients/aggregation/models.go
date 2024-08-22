package aggregation

import (
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
)

type SwapResponseExtended struct {
	SwapResponse

	TxNormalized NormalizedTransactionData `json:"txNormalized"`
}

type NormalizedTransactionData struct {
	Data     []byte             `json:"data"`
	Gas      uint64             `json:"gas"`
	GasPrice *big.Int           `json:"gasPrice"`
	To       gethCommon.Address `json:"to"`
	Value    *big.Int           `json:"value"`
}

type ApproveCallDataResponseExtended struct {
	ApproveCallDataResponse

	TxNormalized NormalizedTransactionData
}
