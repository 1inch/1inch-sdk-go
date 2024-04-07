package transaction_builder

import (
	"context"
	"fmt"
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/1inch/1inch-sdk-go/common"
)

type TransactionBuilder struct {
	wallet common.Wallet

	nonce    uint64
	nonceSet bool

	gasPrice    *big.Int
	gasPriceSet bool

	gas    uint64
	gasSet bool

	to    *gethCommon.Address `rlp:"nil"`
	toSet bool

	value    *big.Int
	valueSet bool

	data    []byte
	dataSet bool

	gasTipCap    *big.Int
	gasTipCapSet bool

	gasFeeCap    *big.Int
	gasFeeCapSet bool
}

func (t *TransactionBuilder) SetData(d []byte) common.TransactionBuilder {
	if d == nil {
		return t
	}
	t.data = d
	t.dataSet = true
	return t
}

func (t *TransactionBuilder) SetNonce(n uint64) common.TransactionBuilder {
	t.nonce = n
	t.nonceSet = true
	return t
}

func (t *TransactionBuilder) SetGasPrice(g *big.Int) common.TransactionBuilder {
	if g == nil {
		return t
	}
	t.gasPrice = g
	t.gasPriceSet = true
	return t
}

func (t *TransactionBuilder) SetGas(g uint64) common.TransactionBuilder {
	t.gas = g
	t.gasSet = true
	return t
}

func (t *TransactionBuilder) SetValue(v *big.Int) common.TransactionBuilder {
	if v == nil {
		return t
	}
	t.value = v
	t.valueSet = true
	return t
}

func (t *TransactionBuilder) SetTo(address *gethCommon.Address) common.TransactionBuilder {
	if address == nil {
		return t
	}
	t.to = address
	t.toSet = true
	return t
}

func (t *TransactionBuilder) SetGasTipCap(g *big.Int) common.TransactionBuilder {
	if g == nil {
		return t
	}
	t.gasTipCap = g
	t.gasTipCapSet = true
	return t
}

func (t *TransactionBuilder) SetGasFeeCap(g *big.Int) common.TransactionBuilder {
	if g == nil {
		return t
	}
	t.gasFeeCap = g
	t.gasFeeCapSet = true
	return t
}

func (t *TransactionBuilder) BuildLegacyTx(ctx context.Context) (*types.Transaction, error) {
	if !t.toSet && !t.dataSet {
		return nil, fmt.Errorf("transaction without data and to params is invalid, specify the params")
	}

	if !t.nonceSet {
		nonce, err := t.wallet.Nonce(ctx)
		if err != nil {
			return nil, err
		}
		t.nonce = nonce
		t.nonceSet = true
	}

	if !t.gasPriceSet {
		gasPrice, err := t.wallet.GetGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		t.gasPrice = gasPrice
		t.gasPriceSet = true
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    t.nonce,
		GasPrice: t.gasPrice,
		Gas:      t.gas,
		To:       t.to,
		Value:    t.value,
		Data:     t.data,
	}), nil
}

func (t *TransactionBuilder) BuildDynamicTx(ctx context.Context) (*types.Transaction, error) {
	if !t.wallet.IsEIP1559Applicable() {
		return nil, fmt.Errorf("current chainId is not supported for dynamic tx")
	}

	if !t.toSet && !t.dataSet {
		return nil, fmt.Errorf("transaction without data and to params is invalid, specify the params")
	}

	if !t.nonceSet {
		nonce, err := t.wallet.Nonce(ctx)
		if err != nil {
			return nil, err
		}
		t.nonce = nonce
		t.nonceSet = true
	}

	if !t.gasTipCapSet {
		gasTipCap, err := t.wallet.GetGasTipCap(ctx)
		if err != nil {
			return nil, err
		}
		t.gasTipCap = gasTipCap
		t.gasTipCapSet = true
	}

	if !t.gasFeeCapSet {
		gasPrice, err := t.wallet.GetGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		t.gasFeeCap = gasPrice
		t.gasFeeCapSet = true
	}

	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(t.wallet.ChainId()),
		Nonce:     t.nonce,
		GasTipCap: t.gasTipCap,
		GasFeeCap: t.gasFeeCap,
		Gas:       t.gas,
		To:        t.to,
		Value:     t.value,
		Data:      t.data,
	}), nil
}

func (t *TransactionBuilder) Build(ctx context.Context) (*types.Transaction, error) {
	if t.wallet.IsEIP1559Applicable() {
		return t.BuildDynamicTx(ctx)
	}
	return t.BuildLegacyTx(ctx)
}
