package onchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TxConfig struct {
	Description   string
	PublicAddress common.Address
	PrivateKey    string
	ChainId       *big.Int
	Value         *big.Int
	To            string
	Data          []byte
}

type Erc20ApprovalConfig struct {
	ChainId        int
	Key            string
	Erc20Address   common.Address
	PublicAddress  common.Address
	SpenderAddress common.Address
}

type Erc20RevokeConfig struct {
	ChainId                 int
	Key                     string
	Erc20Address            common.Address
	PublicAddress           common.Address
	SpenderAddress          common.Address
	AllowanceDecreaseAmount *big.Int
}
