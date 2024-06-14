package web3_provider

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (w Wallet) GenerateApproveCallData(addressTo string, amount uint64) (string, error) {
	spenderAddress := common.HexToAddress(addressTo)

	callData, err := w.erc20ABI.Pack("approve", spenderAddress, big.NewInt(int64(amount)))
	if err != nil {
		return "", fmt.Errorf("failed to pack ABI data: %v", err)
	}

	return common.Bytes2Hex(callData), nil
}
