package web3_provider

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/1inch/1inch-sdk-go/constants"
)

type Wallet struct {
	ethClient  *ethclient.Client
	address    *common.Address
	privateKey *ecdsa.PrivateKey
	chainId    *big.Int
	erc20ABI   *abi.ABI
}

func DefaultWalletProvider(pk string, nodeURL string, ChainId uint64) (*Wallet, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(constants.Erc20ABI)) // Make a generic version of this ABI
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}
	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client: %v", err)
	}

	publicKey := privateKey.Public()
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	return &Wallet{
		ethClient:  ethClient,
		address:    &address,
		privateKey: privateKey,
		chainId:    big.NewInt(int64(ChainId)),
		erc20ABI:   &erc20ABI,
	}, nil
}
