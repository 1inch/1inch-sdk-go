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
	"github.com/1inch/1inch-sdk-go/internal/web3-provider/multicall"
)

type Wallet struct {
	multicall             *multicall.Client
	ethClient             *ethclient.Client
	address               *common.Address
	privateKey            *ecdsa.PrivateKey
	chainId               *big.Int
	erc20ABI              *abi.ABI
	seriesNonceManagerABI *abi.ABI
}

func DefaultWalletProvider(pk string, nodeURL string, chainId uint64) (*Wallet, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(constants.Erc20ABI))
	if err != nil {
		return nil, err
	}
	seriesNonceManagerABI, err := abi.JSON(strings.NewReader(constants.SeriesNonceManagerABI))
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize private key: %w", err)
	}
	ethClient, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client: %w", err)
	}

	publicKey := privateKey.Public()
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	m, err := multicall.NewMulticall(ethClient, chainId)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		multicall:             m,
		ethClient:             ethClient,
		address:               &address,
		privateKey:            privateKey,
		chainId:               big.NewInt(int64(chainId)),
		erc20ABI:              &erc20ABI,
		seriesNonceManagerABI: &seriesNonceManagerABI,
	}, nil
}

func DefaultWalletOnlyProvider(pk string, chainId uint64) (*Wallet, error) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey))

	return &Wallet{
		address:    &address,
		privateKey: privateKey,
		chainId:    big.NewInt(int64(chainId)),
	}, nil
}
