package web3_provider

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TransactionExecutor interface {
}

type Wallet struct {
	ethClient  ethclient.Client
	address    common.Address
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
}

func DefaultWalletProvider(pk string, nodeURL string, chainID *big.Int) *Wallet {
	return &Wallet{
		ethClient:  ethclient.Client{},
		address:    common.Address{},
		privateKey: nil,
		chainID:    nil,
	}
}
