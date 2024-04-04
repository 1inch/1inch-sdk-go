package aggregation

import (
	"encoding/hex"
	"fmt"
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/1inch/1inch-sdk-go/aggregation/models"
	"github.com/1inch/1inch-sdk-go/constants"
)

func (c Client) BuildSwapTransaction(d *models.SwapResponse, nonce uint64, gasPrice, gasTipCap *big.Int) (*types.Transaction, error) {
	to := gethCommon.HexToAddress(d.Tx.To)

	value, ok := new(big.Int).SetString(d.Tx.Value, 10)
	if !ok {
		return nil, fmt.Errorf("failed to convert d.Tx.Value to big.Int")
	}

	data, err := hex.DecodeString(d.Tx.Data[2:])
	if err != nil {
		return nil, err
	}

	isDynamicFeeApplicable := gasTipCap != nil && !(c.chainId == constants.BscChainId || c.chainId == constants.AuroraChainId || c.chainId == constants.ZkSyncEraChainId || c.chainId == constants.FantomChainId)

	if isDynamicFeeApplicable {
		return types.NewTx(&types.DynamicFeeTx{
			ChainID:   big.NewInt(int64(c.chainId)),
			Nonce:     nonce,
			Gas:       uint64(d.Tx.Gas),
			To:        &to,
			Value:     value,
			Data:      data,
			GasTipCap: gasTipCap,
			GasFeeCap: gasPrice,
		}), nil
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      uint64(d.Tx.Gas),
		To:       &to,
		Value:    value,
		Data:     data,
	}), nil
}
