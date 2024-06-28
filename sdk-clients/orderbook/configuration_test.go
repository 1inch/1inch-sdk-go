package orderbook

import (
	"testing"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigurationAPI(t *testing.T) {
	configAPI, err := NewConfigurationAPI(1, "https://api.example.com", "apikey123")

	assert.NoError(t, err)
	assert.NotNil(t, configAPI)
	assert.Equal(t, uint64(1), configAPI.API.chainId)
	assert.Equal(t, "https://api.example.com", configAPI.ApiURL)
	assert.Equal(t, "apikey123", configAPI.ApiKey)
}

func TestNewConfigurationWallet(t *testing.T) {
	testCases := []struct {
		name       string
		privateKey string
		nodeURL    string
		chainId    uint64
		wantErr    bool
	}{
		{
			name:       "Valid inputs",
			privateKey: "965e092fdfc08940d2bd05c7b5c7e1c51e283e92c7f52bbf1408973ae9a9acb7",
			nodeURL:    "https://localhost:8545",
			chainId:    constants.EthereumChainId,
			wantErr:    false,
		},
		{
			name:       "Invalid private key",
			privateKey: "invalidkey",
			nodeURL:    "https://localhost:8545",
			chainId:    constants.EthereumChainId,
			wantErr:    true,
		},
		{
			name:       "Empty private key",
			privateKey: "",
			nodeURL:    "https://localhost:8545",
			chainId:    constants.EthereumChainId,
			wantErr:    true,
		},
		{
			name:       "Empty node URL",
			privateKey: "45779132284297842289675692834abcdef9876543212345678900987654321",
			nodeURL:    "",
			chainId:    constants.EthereumChainId,
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConfigurationWallet(tc.nodeURL, tc.privateKey, tc.chainId)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
