package transaction_builder

import (
	"context"
	"math/big"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/common"
)

func TestTransactionBuilder_Build(t *testing.T) {
	w := NewMyWallet(gethCommon.HexToAddress("0x0000000000000000000000000000000000000000"), big.NewInt(1))
	factory := NewFactory(w)
	ctx := context.Background()

	testSpecialNonSetValue := uint64(999)
	to := gethCommon.HexToAddress("0x000000000000000000000000000000000000dead")
	var tests = []struct {
		name           string
		nonce          uint64
		to             *gethCommon.Address
		gas            uint64
		value          *big.Int
		gasFee         *big.Int
		gasTip         *big.Int
		expectError    bool
		expectedValues map[string]interface{}
	}{
		{
			name:        "All fields set",
			nonce:       uint64(1),
			to:          &to,
			gas:         10_000,
			value:       big.NewInt(100000),
			gasFee:      big.NewInt(25),
			gasTip:      big.NewInt(10000),
			expectError: false,
			expectedValues: map[string]interface{}{
				"Nonce":     uint64(1),
				"To":        gethCommon.HexToAddress("0x000000000000000000000000000000000000dead"),
				"Value":     big.NewInt(100000),
				"GasFeeCap": big.NewInt(25),
				"GasTipCap": big.NewInt(10000),
				"Gas":       uint64(10_000),
				"ChainId":   big.NewInt(1),
			},
		},
		{
			name:        "Missing gasFee and gasTip",
			nonce:       uint64(2),
			gas:         20_000,
			value:       big.NewInt(200000),
			to:          &to,
			gasFee:      nil,
			gasTip:      nil,
			expectError: false,
			expectedValues: map[string]interface{}{
				"Nonce":     uint64(2),
				"To":        to,
				"Value":     big.NewInt(200000),
				"Gas":       uint64(20_000),
				"GasTipCap": big.NewInt(23),
				"GasFeeCap": big.NewInt(23),
				"ChainId":   big.NewInt(1),
			},
		},
		{
			name:        "Missing gasFee and gasTip and Nonce",
			nonce:       testSpecialNonSetValue,
			gas:         20_000,
			value:       big.NewInt(200000),
			gasFee:      nil,
			gasTip:      nil,
			to:          &to,
			expectError: false,
			expectedValues: map[string]interface{}{
				"Nonce":     uint64(44),
				"To":        to,
				"Value":     big.NewInt(200000),
				"ChainId":   big.NewInt(1),
				"Gas":       uint64(20_000),
				"GasTipCap": big.NewInt(23),
				"GasFeeCap": big.NewInt(23),
			},
		},
		{
			name:        "Missing to and data",
			nonce:       testSpecialNonSetValue,
			gas:         20_000,
			value:       big.NewInt(200000),
			gasFee:      nil,
			gasTip:      nil,
			to:          nil,
			expectError: true,
			expectedValues: map[string]interface{}{
				"To": nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			builder := factory.New()

			if tc.to != nil {
				builder = builder.SetTo(tc.to)
			}

			if tc.gasTip != nil {
				builder = builder.SetGasTipCap(tc.gasTip)
			}

			if tc.gasFee != nil {
				builder = builder.SetGasFeeCap(tc.gasFee)
			}

			if tc.gas != testSpecialNonSetValue {
				builder = builder.SetGas(tc.gas)
			}

			if tc.value != nil {
				builder = builder.SetValue(tc.value)
			}

			if tc.nonce != testSpecialNonSetValue {
				builder = builder.SetNonce(tc.nonce)
			}
			tx, err := builder.Build(ctx)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Check the expected values
				require.Equal(t, tc.expectedValues["Nonce"], tx.Nonce())
				require.Equal(t, tc.expectedValues["To"], *tx.To())
				require.Equal(t, tc.expectedValues["Value"], tx.Value())

				require.Equal(t, tc.expectedValues["GasFeeCap"], tx.GasFeeCap())

				require.Equal(t, tc.expectedValues["GasTipCap"], tx.GasTipCap())

				require.Equal(t, tc.expectedValues["Gas"], tx.Gas())
				require.Equal(t, tc.expectedValues["ChainId"], tx.ChainId())
			}
		})
	}
}

type MyWallet struct {
	address gethCommon.Address
	chainID *big.Int
}

func (w *MyWallet) GetContractDetailsForPermit(ctx context.Context, token gethCommon.Address, spender gethCommon.Address, amount *big.Int, deadline int64) (*common.ContractPermitData, error) {
	//TODO implement me
	panic("implement me")
}

func NewMyWallet(address gethCommon.Address, chainID *big.Int) *MyWallet {
	return &MyWallet{
		address: address,
		chainID: chainID,
	}
}

func (w *MyWallet) Call(ctx context.Context, contractAddress gethCommon.Address, callData []byte) ([]byte, error) {
	return nil, nil
}

func (w *MyWallet) Nonce(ctx context.Context) (uint64, error) {
	return 44, nil
}

func (w *MyWallet) Address() gethCommon.Address {
	return w.address
}

func (w *MyWallet) Balance(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}

func (w *MyWallet) GetGasTipCap(ctx context.Context) (*big.Int, error) {
	// For EIP-1559 transactions
	return big.NewInt(23), nil
}

func (w *MyWallet) GetGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(23), nil
}

func (w *MyWallet) Sign(tx *types.Transaction) (*types.Transaction, error) {
	return tx, nil
}

func (w *MyWallet) SignBytes(data []byte) ([]byte, error) {
	return nil, nil
}

func (w *MyWallet) BroadcastTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
}

func (w *MyWallet) EstimateGas(ctx context.Context, contractAddress gethCommon.Address, callData []byte) (uint64, error) {
	return 0, nil
}

func (w *MyWallet) TransactionReceipt(ctx context.Context, txHash gethCommon.Hash) (*types.Receipt, error) {
	return nil, nil
}

func (w *MyWallet) GetContractDetailsForPermitDaiLike(ctx context.Context, token gethCommon.Address, spender gethCommon.Address, deadline int64) (*common.ContractPermitDataDaiLike, error) {
	return nil, nil
}

func (w *MyWallet) TokenPermitDaiLike(cd common.ContractPermitDataDaiLike) (string, error) {
	return "", nil
}

func (w *MyWallet) TokenPermit(cd common.ContractPermitData) (string, error) {
	return "", nil
}

func (w *MyWallet) IsEIP1559Applicable() bool {
	return true
}

func (w *MyWallet) ChainId() int64 {
	return w.chainID.Int64()
}
