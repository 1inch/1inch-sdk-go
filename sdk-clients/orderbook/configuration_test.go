package orderbook

import (
	"testing"

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
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConfigurationWallet(tc.privateKey)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
