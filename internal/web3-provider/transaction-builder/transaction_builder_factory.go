package transaction_builder

import (
	"github.com/1inch/1inch-sdk-go/common"
)

type TransactionBuilderFactory struct {
	wallet common.Wallet
}

func NewFactory(w common.Wallet) *TransactionBuilderFactory {
	return &TransactionBuilderFactory{
		wallet: w,
	}
}

func (f TransactionBuilderFactory) New() common.TransactionBuilder {
	return &TransactionBuilder{
		wallet:    f.wallet,
		nonce:     0,
		gasPrice:  nil,
		gas:       0,
		to:        nil,
		value:     nil,
		data:      nil,
		gasTipCap: nil,
		gasFeeCap: nil,
	}
}
