package aggregation

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"

	"github.com/1inch/1inch-sdk-go/common"
	"github.com/1inch/1inch-sdk-go/constants"
	transaction_builder "github.com/1inch/1inch-sdk-go/internal/transaction-builder"
)

func TestNewClient(t *testing.T) {
	mockAPI := api{
		chainId: 1,
	}

	mockWallet := NewMyWallet(gethCommon.HexToAddress("0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914"), big.NewInt(constants.EthereumChainId))
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

func TestClient(t *testing.T) {
	mockAPI := api{
		chainId: 1,
		httpExecutor: &mockHttpExecutor{
			ResponseObj: mockedSwapHttpApiResp,
		},
	}

	mockWallet := NewMyWallet(gethCommon.HexToAddress("0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914"), big.NewInt(constants.EthereumChainId))
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

	swapData, err := client.GetSwap(context.Background(), GetSwapParams{
		Src:               "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
		Dst:               "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Amount:            "10000",
		Slippage:          1,
		From:              "0x083fc10ce7e97cafbae0fe332a9c4384c5f54e45",
		IncludeTokensInfo: true,
		IncludeGas:        true,
		IncludeProtocols:  true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, swapData)
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

func (w MyWallet) GetGasEstimate(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return 123, nil
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

type mockHttpExecutor struct {
	Called      bool
	ExecuteErr  error
	ResponseObj any
}

func (m *mockHttpExecutor) ExecuteRequest(ctx context.Context, payload common.RequestPayload, v any) error {
	m.Called = true
	if m.ExecuteErr != nil {
		return m.ExecuteErr
	}

	// Copy the mock response object to v
	if m.ResponseObj != nil && v != nil {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return fmt.Errorf("v must be a non-nil pointer")
		}
		reflect.Indirect(rv).Set(reflect.ValueOf(m.ResponseObj))
	}
	return nil
}

var mockedSwapHttpApiResp = SwapResponse{
	SrcToken: &TokenInfo{
		Address:  "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
		Symbol:   "LDO",
		Name:     "Lido DAO Token",
		Decimals: 18,
		LogoURI:  "https://tokens.1inch.io/0x5a98fcbea516cf06857215779fd812ca3bef1b32.png",
		Tags: []string{
			"tokens",
		},
	},
	DstAmount: "6",
	DstToken: &TokenInfo{
		Address:  "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Symbol:   "WETH",
		Name:     "Wrapped Ether",
		Decimals: 18,
		LogoURI:  "https://tokens.1inch.io/0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2.png",
		Tags: []string{
			"PEG:ETH",
			"tokens",
		},
	},
	Protocols: [][][]SelectedProtocol{
		{
			{
				{
					FromTokenAddress: "0x5a98fcbea516cf06857215779fd812ca3bef1b32",
					Name:             "SUSHI",
					Part:             100,
					ToTokenAddress:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				},
			},
		},
	},
	Tx: TransactionData{
		Data:     "0x0502b1c50000000000000000000000005a98fcbea516cf06857215779fd812ca3bef1b32000000000000000000000000000000000000000000000000000000000000271000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000100000000000000003b6d0340c558f600b34a5f69dd2f0d06cb8a88d829b7420ade8bb62d",
		From:     "0x2c9b2dbdba8a9c969ac24153f5c1c23cb0e63914",
		Gas:      257615,
		GasPrice: "22800337026",
		To:       "0x1111111254eeb25477b68fb85ed929f73a960582",
		Value:    "0",
	},
}
