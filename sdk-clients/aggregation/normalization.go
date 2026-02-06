package aggregation

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/1inch/1inch-sdk-go/constants"
)

func normalizeSwapResponse(resp SwapResponse) (*SwapResponseExtended, error) {
	if !common.IsHexAddress(resp.Tx.To) {
		return nil, fmt.Errorf("invalid to address: %s", resp.Tx.To)
	}
	toAddress := common.HexToAddress(resp.Tx.To)

	gas := uint64(resp.Tx.Gas)

	gasPrice := big.NewInt(0)
	_, ok := gasPrice.SetString(resp.Tx.GasPrice, 10)
	if !ok {
		return nil, fmt.Errorf("invalid gas price: %s", resp.Tx.GasPrice)
	}

	value := big.NewInt(0)
	_, ok = value.SetString(resp.Tx.Value, 10)
	if !ok {
		return nil, fmt.Errorf("invalid tx value: %s", resp.Tx.Value)
	}

	data, err := hexutil.Decode(resp.Tx.Data)
	if err != nil {
		return nil, fmt.Errorf("invalid tx data: %w", err)
	}

	normalizedTx := NormalizedTransactionData{
		Data:     data,
		Gas:      gas,
		GasPrice: gasPrice,
		To:       toAddress,
		Value:    value,
	}

	extendedResp := SwapResponseExtended{
		SwapResponse: resp,
		TxNormalized: normalizedTx,
	}

	return &extendedResp, nil
}

func normalizeApproveCallDataResponse(resp ApproveCallDataResponse) (*ApproveCallDataResponseExtended, error) {
	if !common.IsHexAddress(resp.To) {
		return nil, fmt.Errorf("invalid to address: %s", resp.To)
	}
	toAddress := common.HexToAddress(resp.To)

	gasPrice := big.NewInt(0)
	_, ok := gasPrice.SetString(resp.GasPrice, 10)
	if !ok {
		return nil, fmt.Errorf("invalid gas price: %s", resp.GasPrice)
	}

	value := big.NewInt(0)
	_, ok = value.SetString(resp.Value, 10)
	if !ok {
		return nil, fmt.Errorf("invalid value: %s", resp.Value)
	}

	data, err := hexutil.Decode(resp.Data)
	if err != nil {
		return nil, fmt.Errorf("invalid data: %w", err)
	}

	normalizedTx := NormalizedTransactionData{
		Data:     data,
		Gas:      constants.Erc20ApproveGas,
		GasPrice: gasPrice,
		To:       toAddress,
		Value:    value,
	}

	extendedResp := ApproveCallDataResponseExtended{
		ApproveCallDataResponse: resp,
		TxNormalized:            normalizedTx,
	}

	return &extendedResp, nil
}
