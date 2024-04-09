package multicall

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type CallData struct {
	To         string `json:"to"`
	Data       string `json:"data"`
	MethodName string `json:"methodName,omitempty"`
	Gas        uint64 `json:"-"`
}

type request struct {
	To   common.Address
	Data []byte
}

type response struct {
	Results          [][]byte
	LastSuccessIndex *big.Int
}
