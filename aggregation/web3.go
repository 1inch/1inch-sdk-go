package aggregation

import (
	"encoding/hex"
	"errors"
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/1inch/1inch-sdk-go/aggregation/models"
)

func (c Client) BuildSwapTransaction(d *models.SwapResponse, nonce uint64, gasTipCap *big.Int, gasFeeCap *big.Int) (*types.Transaction, error) {
	to := gethCommon.HexToAddress(d.Tx.To)

	value := new(big.Int)
	value, ok := value.SetString(d.Tx.Value, 10)
	if !ok {
		return nil, errors.New("failed to convert d.Tx.Value to big.Int")
	}

	data, err := hex.DecodeString(d.Tx.Data[2:])
	if err != nil {
		return nil, err
	}

	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(int64(c.chainId)),
		Nonce:     nonce,
		Gas:       uint64(d.Tx.Gas),
		To:        &to,
		Value:     value,
		Data:      data,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
	}), nil
}
