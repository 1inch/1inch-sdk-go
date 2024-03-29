package web3_provider

import (
	"context"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/abis"
)

type ContractPermitData struct {
	FromToken     string
	Name          string
	Version       string
	PublicAddress string
	ChainId       int
	Key           string
	Nonce         int64
	Deadline      int64
}

func (w Wallet) TokenPermit(cd ContractPermitData) (string, error) {
	return "", nil
}

func (w Wallet) GetContractDetailsForPermit(ctx context.Context, token common.Address, deadline int64) (*ContractPermitData, error) {
	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
	if err != nil {
		return nil, err
	}

	contractName, err := callAndUnpackContractMethod(ctx, token, parsedABI, &w.ethClient, "name")
	if err != nil {
		return nil, err
	}

	contractVersion, err := callAndUnpackContractMethod(ctx, token, parsedABI, &w.ethClient, "version")
	if err != nil {
		return nil, err
	}

	contractNonceStr, err := callAndUnpackContractMethod(ctx, token, parsedABI, &w.ethClient, "nonce", []common.Address{token})
	if err != nil {
		return nil, err
	}

	contractNonce, err := strconv.ParseInt(contractNonceStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &ContractPermitData{
		FromToken:     token.Hex(),
		PublicAddress: w.address.Hex(),
		ChainId:       int(w.chainID.Int64()),
		Name:          contractName,
		Version:       contractVersion,
		Nonce:         contractNonce,
		Deadline:      deadline,
	}, nil
}

func callAndUnpackContractMethod(ctx context.Context, token common.Address, parsedABI abi.ABI, client *ethclient.Client, methodName string, methodArgs ...interface{}) (string, error) {
	data, err := parsedABI.Pack(methodName, methodArgs...)
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		To:   &token,
		Data: data,
	}

	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return "", err
	}

	var returnValue string
	err = parsedABI.UnpackIntoInterface(&returnValue, methodName, result)
	if err != nil {
		return "", err
	}

	return returnValue, nil
}
