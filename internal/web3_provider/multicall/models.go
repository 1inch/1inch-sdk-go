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

type MulticallParams struct {
	Client   Client
	ChainId  int
	Calldata []CallData
}

type CallParams struct {
	Client          Client
	Data            []byte
	ContractAddress string
	Block           *big.Int
}

type request struct {
	To   common.Address
	Data []byte
}

type response struct {
	Results          [][]byte
	LastSuccessIndex *big.Int
}
