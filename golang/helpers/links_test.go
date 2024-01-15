package helpers

import (
	"testing"

	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/stretchr/testify/assert"
)

func TestGetBlockExplorerTxLinkInfo(t *testing.T) {
	testCases := []struct {
		description string
		chainId     int
		txHash      string
		expected    string
	}{
		{
			description: "Ethereum mainnet transaction",
			chainId:     chains.Ethereum,
			txHash:      "0x123",
			expected:    "View it Etherscan here: https://etherscan.io/tx/0x123\n",
		},
		{
			description: "Polygon network transaction",
			chainId:     chains.Polygon,
			txHash:      "0x456",
			expected:    "View it PolygonScan here: https://polygonscan.com/tx/0x456\n",
		},
		{
			description: "Unknown network",
			chainId:     111,
			txHash:      "0x456",
			expected:    "Tx Id: 0x456\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Call the function
			output := GetBlockExplorerTxLinkInfo(tc.chainId, tc.txHash)

			// Check the output
			assert.Equal(t, tc.expected, output)
		})
	}
}
