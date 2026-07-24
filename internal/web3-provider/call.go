package web3_provider

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
)

func (w Wallet) Call(ctx context.Context, contractAddress gethCommon.Address, callData []byte) ([]byte, error) {
	if w.ethClient == nil {
		return nil, fmt.Errorf("wallet has no node connection: create it with a node URL to make on-chain calls")
	}
	nodeMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}
	resp, err := w.ethClient.CallContract(ctx, nodeMsg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	return resp, nil
}
