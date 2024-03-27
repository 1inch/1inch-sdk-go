package web3_provider

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type TransactionExecutor interface {
}

type Wallet struct {
	ethClient  ethclient.Client
	address    common.Address
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
}
