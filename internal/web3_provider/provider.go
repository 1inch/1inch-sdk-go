package web3_provider

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TransactionExecutor interface {
}

type Wallet struct {
	ethClient     ethclient.Client
	address       common.Address
	privateKeyHex string
}
