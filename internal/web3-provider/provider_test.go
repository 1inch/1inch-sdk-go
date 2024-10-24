package web3_provider

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultWalletProvider_FailureCases(t *testing.T) {
	testcases := []struct {
		description string
		privateKey  string
		nodeURL     string
		chainID     uint64
		expectError string
	}{
		{
			description: "Invalid Private Key - blank",
			privateKey:  "",
			nodeURL:     "https://mainnet.infura.io/v3/randomProjectId",
			chainID:     1,
			expectError: "failed to initialize private key",
		},
		{
			description: "Invalid Private Key - bad characters",
			privateKey:  "malformed_private_key",
			nodeURL:     "https://mainnet.infura.io/v3/randomProjectId",
			chainID:     1,
			expectError: "failed to initialize private key",
		},
		{
			description: "Invalid node URL",
			privateKey:  "85cc05822dc41dbd5253767374b12ca1d08d4d347af2c0bf7bbff8edc3dfa950",
			nodeURL:     "",
			chainID:     1,
			expectError: "failed to create eth client",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {
			_, err := DefaultWalletProvider(tc.privateKey, tc.nodeURL, tc.chainID)
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectError)
		})
	}
}
