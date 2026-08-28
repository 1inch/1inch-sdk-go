package tokens

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProviderTokenDecoding pins how token API responses decode into the
// corrected generated types: tags as {provider, value} objects, chain ids as
// exact integers, and optional fields as plain values.
func TestProviderTokenDecoding(t *testing.T) {
	tests := []struct {
		name     string
		payload  string
		expected ProviderTokenDtoFixed
	}{
		{
			name: "Full token with object tags",
			payload: `{
				"address": "0x111111111117dc0aa78b770fa6a738034120c302",
				"chainId": 1,
				"decimals": 18,
				"name": "1INCH Token",
				"symbol": "1INCH",
				"providers": ["1inch"],
				"eip2612": true,
				"logoURI": "https://tokens.1inch.io/1inch.png",
				"tags": [{"provider": "1inch", "value": "tokens"}]
			}`,
			expected: ProviderTokenDtoFixed{
				Address:   "0x111111111117dc0aa78b770fa6a738034120c302",
				ChainId:   1,
				Decimals:  18,
				Name:      "1INCH Token",
				Symbol:    "1INCH",
				Providers: []string{"1inch"},
				Eip2612:   true,
				LogoURI:   "https://tokens.1inch.io/1inch.png",
				Tags:      []TagDto{{Provider: "1inch", Value: "tokens"}},
			},
		},
		{
			name: "Aurora chain id decodes exactly",
			payload: `{
				"address": "0x0000000000000000000000000000000000000000",
				"chainId": 1313161554,
				"decimals": 18,
				"name": "t",
				"symbol": "T",
				"providers": [],
				"tags": []
			}`,
			expected: ProviderTokenDtoFixed{
				Address:   "0x0000000000000000000000000000000000000000",
				ChainId:   1313161554,
				Decimals:  18,
				Name:      "t",
				Symbol:    "T",
				Providers: []string{},
				Tags:      []TagDto{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out ProviderTokenDtoFixed
			require.NoError(t, json.Unmarshal([]byte(tc.payload), &out))
			assert.Equal(t, tc.expected, out)
		})
	}
}
