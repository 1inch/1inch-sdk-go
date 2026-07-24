package transaction_builder

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/1inch/1inch-sdk-go/v4/common"
)

type TransactionBuilder struct {
	wallet    common.Wallet
	nonce     *uint64
	gasPrice  *big.Int
	gas       *uint64
	to        *gethCommon.Address `rlp:"nil"`
	value     *big.Int
	data      []byte
	gasTipCap *big.Int
	gasFeeCap *big.Int
}

func (t *TransactionBuilder) SetData(d []byte) common.TransactionBuilder {
	if d == nil {
		return t
	}
	t.data = d
	return t
}

func (t *TransactionBuilder) SetNonce(n uint64) common.TransactionBuilder {
	t.nonce = &n
	return t
}

func (t *TransactionBuilder) SetGasPrice(g *big.Int) common.TransactionBuilder {
	if g == nil {
		return t
	}
	t.gasPrice = g
	return t
}

func (t *TransactionBuilder) SetGas(g uint64) common.TransactionBuilder {
	t.gas = &g
	return t
}

func (t *TransactionBuilder) SetValue(v *big.Int) common.TransactionBuilder {
	if v == nil {
		return t
	}
	t.value = v
	return t
}

func (t *TransactionBuilder) SetTo(address *gethCommon.Address) common.TransactionBuilder {
	if address == nil {
		return t
	}
	t.to = address
	return t
}

func (t *TransactionBuilder) SetGasTipCap(g *big.Int) common.TransactionBuilder {
	if g == nil {
		return t
	}
	t.gasTipCap = g
	return t
}

func (t *TransactionBuilder) SetGasFeeCap(g *big.Int) common.TransactionBuilder {
	if g == nil {
		return t
	}
	t.gasFeeCap = g
	return t
}

func (t *TransactionBuilder) BuildLegacyTx(ctx context.Context) (*types.Transaction, error) {
	if t.to == nil && t.data == nil {
		return nil, fmt.Errorf("transaction requires data or to address")
	}

	if t.nonce == nil {
		nonce, err := t.wallet.Nonce(ctx)
		if err != nil {
			return nil, err
		}
		t.nonce = &nonce
	}

	if t.gasPrice == nil {
		gasPrice, err := t.wallet.GetGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		t.gasPrice = gasPrice
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    *t.nonce,
		GasPrice: t.gasPrice,
		Gas:      *t.gas,
		To:       t.to,
		Value:    t.value,
		Data:     t.data,
	}), nil
}

func (t *TransactionBuilder) BuildDynamicTx(ctx context.Context) (*types.Transaction, error) {
	if !t.wallet.IsEIP1559Applicable() {
		return nil, fmt.Errorf("unsupported: dynamic transactions on this chain")
	}

	if t.to == nil && t.data == nil {
		return nil, fmt.Errorf("transaction requires data or to address")
	}

	if t.nonce == nil {
		nonce, err := t.wallet.Nonce(ctx)
		if err != nil {
			return nil, err
		}
		t.nonce = &nonce
	}

	if t.gasTipCap == nil {
		gasTipCap, err := t.wallet.GetGasTipCap(ctx)
		if err != nil {
			return nil, err
		}
		t.gasTipCap = gasTipCap
	}

	if t.gasFeeCap == nil {
		gasPrice, err := t.wallet.GetGasPrice(ctx)
		if err != nil {
			return nil, err
		}
		// The node's suggested gas price tracks the current base fee, so using it
		// directly as the fee cap makes the transaction invalid as soon as the base
		// fee rises before inclusion. Doubling the suggestion keeps the cap at or
		// above twice the base fee plus the tip; the network only charges base fee
		// plus tip, so the headroom adds no cost.
		t.gasFeeCap = new(big.Int).Mul(gasPrice, big.NewInt(2))
		if t.gasFeeCap.Cmp(t.gasTipCap) < 0 {
			t.gasFeeCap = new(big.Int).Set(t.gasTipCap)
		}
	}

	if t.gas == nil {
		// Fee fields are omitted from the estimate: gas usage does not depend on
		// them, and including a price makes the node reject the estimate whenever
		// the base fee moves or the account cannot prepay at the capped price
		gas, err := t.wallet.GetGasEstimate(ctx, ethereum.CallMsg{
			From:  t.wallet.Address(),
			To:    t.to,
			Value: t.value,
			Data:  t.data,
		})
		if err != nil {
			return nil, err
		}
		t.gas = &gas
	}

	return types.NewTx(&types.DynamicFeeTx{
		ChainID:   big.NewInt(t.wallet.ChainId()),
		Nonce:     *t.nonce,
		GasTipCap: t.gasTipCap,
		GasFeeCap: t.gasFeeCap,
		Gas:       *t.gas,
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
