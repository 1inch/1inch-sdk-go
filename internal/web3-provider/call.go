package web3_provider

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
)

func (w Wallet) Call(ctx context.Context, contractAddress gethCommon.Address, callData []byte) ([]byte, error) {
	nodeMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}
	resp, err := w.ethClient.CallContract(ctx, nodeMsg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: error: %s", err)
	}

	return resp, nil
}
