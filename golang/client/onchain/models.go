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

type PermitSignatureConfig struct {
	FromToken     string
	Name          string
	PublicAddress string
	ChainId       int
	Key           string
	Nonce         int64
	Deadline      int64
}

type PermitParamsConfig struct {
	Owner     string
	Spender   string
	Value     *big.Int
	Deadline  int64
	Signature string
}
