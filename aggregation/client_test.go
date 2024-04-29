package aggregation

import (
	"context"
	"math/big"
	"testing"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"
	transaction_builder "github.com/1inch/1inch-sdk-go/internal/transaction-builder"
)

func TestNewClient(t *testing.T) {
	mockAPI := api{
		chainId:      1,
		httpExecutor: nil,
	}

	mockWallet := NewMyWallet(gethCommon.HexToAddress("0x0000000000000000000000000000000000000000"), big.NewInt(constants.EthereumChainId))
	mockTxBuilder := transaction_builder.NewFactory(mockWallet)

	cfg := &Configuration{
		APIConfiguration: &ConfigurationAPI{
			API: mockAPI,
		},
		WalletConfiguration: &ConfigurationWallet{
			Wallet:    mockWallet,
			TxBuilder: mockTxBuilder,
		},
	}

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, mockAPI, client.api)
	assert.Equal(t, mockWallet, client.Wallet)
	assert.Equal(t, mockTxBuilder, client.TxBuilder)

	cfgOnlyAPI := &Configuration{
		APIConfiguration: &ConfigurationAPI{
			API: mockAPI,
		},
	}

	clientOnlyAPI, err := NewClient(cfgOnlyAPI)
	assert.NoError(t, err)
	assert.NotNil(t, clientOnlyAPI)
	assert.Equal(t, mockAPI, clientOnlyAPI.api)
	assert.Nil(t, clientOnlyAPI.Wallet)
	assert.Nil(t, clientOnlyAPI.TxBuilder)
}

type MyWallet struct {
	address gethCommon.Address
	chainID *big.Int

	WantedGasPrice  *big.Int
	WantedCall      []byte
	WantedNonce     uint64
	WantedBalance   *big.Int
	WantedGasTipCap *big.Int
}

func (w *MyWallet) GetContractDetailsForPermit(ctx context.Context, token gethCommon.Address, spender gethCommon.Address, amount *big.Int, deadline int64) (*common.ContractPermitData, error) {
	return nil, nil
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
