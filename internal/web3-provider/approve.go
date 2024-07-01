package web3_provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func (w Wallet) GenerateApproveCallData(addressTo string, amount uint64) (string, error) {
	spenderAddress := common.HexToAddress(addressTo)

	callData, err := w.erc20ABI.Pack("approve", spenderAddress, big.NewInt(int64(amount)))
	if err != nil {
		return "", fmt.Errorf("failed to pack ABI data: %v", err)
	}

	return "0x" + common.Bytes2Hex(callData), nil
}

func (w Wallet) TokenAllowance(ctx context.Context, tokenAddress string, spenderAddress string) (*big.Int, error) {
	token := common.HexToAddress(tokenAddress)
	spender := common.HexToAddress(spenderAddress)

	callData, err := w.erc20ABI.Pack("allowance", w.address, spender)
	if err != nil {
		return nil, fmt.Errorf("failed to pack ABI data: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &token,
		Data: callData,
	}

	result, err := w.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	var allowance *big.Int
	err = w.erc20ABI.UnpackIntoInterface(&allowance, "allowance", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %v", err)
	}

	return allowance, nil
}
