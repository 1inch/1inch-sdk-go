package aggregation

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/constants"
)

func normalizeSwapResponse(resp SwapResponse) (*SwapResponseExtended, error) {
	toAddress := common.HexToAddress(resp.Tx.To)
	if !common.IsHexAddress(resp.Tx.To) {
		return nil, errors.New("invalid 'To' address")
	}

	gas := uint64(resp.Tx.Gas)

	gasPrice := big.NewInt(0)
	_, ok := gasPrice.SetString(resp.Tx.GasPrice, 10)
	if !ok {
		return nil, errors.New("invalid 'GasPrice' value")
	}

	value := big.NewInt(0)
	_, ok = value.SetString(resp.Tx.Value, 10)
	if !ok {
		return nil, errors.New("invalid 'Value' value")
	}

	data, err := hex.DecodeString(resp.Tx.Data[2:])
	if err != nil {
		return nil, errors.New("invalid 'Data' value")
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
	toAddress := common.HexToAddress(resp.To)
	if !common.IsHexAddress(resp.To) {
		return nil, errors.New("invalid 'To' address")
	}

	gasPrice := big.NewInt(0)
	_, ok := gasPrice.SetString(resp.GasPrice, 10)
	if !ok {
		return nil, errors.New("invalid 'GasPrice' value")
	}

	value := big.NewInt(0)
	_, ok = value.SetString(resp.Value, 10)
	if !ok {
		return nil, errors.New("invalid 'Value' value")
	}

	data, err := hex.DecodeString(resp.Data[2:])
	if err != nil {
		return nil, errors.New("invalid 'Data' value")
	}

	normalizedTx := NormalizedTransactionData{
		Data:     data,
		Gas:      constants.ERC20_APPROVE_GAS,
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
