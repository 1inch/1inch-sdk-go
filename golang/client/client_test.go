package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var SimpleEthereumConfig = Config{
	TargetEnvironment: EnvironmentProduction,
	ChainId:           chains.Ethereum,
	DevPortalApiKey:   os.Getenv("DEV_PORTAL_TOKEN"),
	Web3HttpProvider:  os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
	WalletKey:         os.Getenv("WALLET_KEY"),
}

func TestNewConfig(t *testing.T) {
	testcases := []struct {
		description              string
		config                   Config
		expectedEnvironment      string
		expectedErrorDescription string
	}{
		{
			description: "Production",
			config: Config{
				TargetEnvironment: EnvironmentProduction,
				DevPortalApiKey:   "abc123",
				Web3HttpProvider:  os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
				ChainId:           chains.Ethereum,
			},
			expectedEnvironment:      baseUrlProduction.Host,
			expectedErrorDescription: "",
		},
		{
			description: "Production (excluded entry)",
			config: Config{
				DevPortalApiKey:  "abc123",
				Web3HttpProvider: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
				ChainId:          chains.Ethereum,
			},
			expectedEnvironment:      baseUrlProduction.Host,
			expectedErrorDescription: "",
		},
		{
			description: "Staging",
			config: Config{
				TargetEnvironment: EnvironmentStaging,
				DevPortalApiKey:   "abc123",
				Web3HttpProvider:  os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
				ChainId:           chains.Ethereum,
			},
			expectedEnvironment:      baseUrlStaging.Host,
			expectedErrorDescription: "",
		},
		{
			description: "Error - unrecognized environment",
			config: Config{
				TargetEnvironment: Environment("invalid"),
				DevPortalApiKey:   "abc123",
				Web3HttpProvider:  os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
				ChainId:           chains.Ethereum,
			},
			expectedEnvironment:      baseUrlStaging.Host,
			expectedErrorDescription: "unrecognized environment: invalid",
		},
		{
			description: "Error - no API key",
			config: Config{
				DevPortalApiKey:  "",
				Web3HttpProvider: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
				ChainId:          chains.Ethereum,
			},
			expectedErrorDescription: "config validation error: API key is required",
		},
		{
			description: "Error - no web3 provider key",
			config: Config{
				DevPortalApiKey: "123",
				ChainId:         chains.Ethereum,
			},
			expectedErrorDescription: "config validation error: web3 provider URL is required",
		},
		{
			description: "Error - no chain ID",
			config: Config{
				DevPortalApiKey:  "123",
				Web3HttpProvider: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
			},
			expectedErrorDescription: "config validation error: chain ID is required",
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			c, err := NewClient(tc.config)
			if tc.expectedErrorDescription != "" {
				if err == nil {
					assert.FailNow(t, "Expected error message, but error was nil")
				}
				assert.Equal(t, tc.expectedErrorDescription, err.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expectedEnvironment, c.BaseURL.Host)
		})
	}
}
