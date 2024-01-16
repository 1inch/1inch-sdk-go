package onchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type GetTxConfig struct {
	ChainId     *big.Int
	FromAddress common.Address
	Value       *big.Int
	To          string
	Data        []byte
}
