package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/svanas/1inch-sdk/golang/helpers/consts/chains"
)

var SimpleEthereumConfig = Config{
	DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
	Web3HttpProviders: []Web3ProviderConfig{
		{
			ChainId: chains.Ethereum,
			Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
		},
	},
}

func TestNewConfig(t *testing.T) {
	testcases := []struct {
		description              string
		config                   Config
		expectedErrorDescription string
	}{
		{
			description: "Success",
			config: Config{
				DevPortalApiKey: "abc123",
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Ethereum,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
					},
				},
			},
			expectedErrorDescription: "",
		},
		{
			description: "Error - no API key",
			config: Config{
				DevPortalApiKey: "",
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Ethereum,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
					},
				},
			},
			expectedErrorDescription: "config validation error: API key is required",
		},
		{
			description: "Error - no web3 provider key",
			config: Config{
				DevPortalApiKey: "123",
			},
			expectedErrorDescription: "config validation error: at least one web3 provider URL is required",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			_, err := NewClient(tc.config)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
