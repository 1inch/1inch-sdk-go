package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/chains"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/client/models"
)

var SimpleEthereumConfig = models.ClientConfig{
	DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
	Web3HttpProviders: []models.Web3Provider{
		{
			ChainId: chains.Ethereum,
			Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
		},
	},
}

func TestNewConfig(t *testing.T) {
	testcases := []struct {
		description              string
		config                   models.ClientConfig
		expectedErrorDescription string
	}{
		{
			description: "Success",
			config: models.ClientConfig{
				DevPortalApiKey: "abc123",
				Web3HttpProviders: []models.Web3Provider{
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
			config: models.ClientConfig{
				DevPortalApiKey: "",
				Web3HttpProviders: []models.Web3Provider{
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
			config: models.ClientConfig{
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
