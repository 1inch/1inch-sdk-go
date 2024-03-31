package web3_provider

import (
	"crypto/ecdsa"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/abis"
)

type TransactionExecutor interface {
}

type Wallet struct {
	ethClient  ethclient.Client
	address    common.Address
	privateKey *ecdsa.PrivateKey
	chainID    *big.Int
	erc20ABI   *abi.ABI
}

func DefaultWalletProvider(pk string, nodeURL string, chainID *big.Int) *Wallet {
	erc20ABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
	if err != nil {
		return nil
	}
	return &Wallet{
		ethClient:  ethclient.Client{},
		address:    common.Address{},
		privateKey: nil,
		chainID:    nil,
		erc20ABI:   &erc20ABI,
	}
}
