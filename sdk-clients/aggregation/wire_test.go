package aggregation

import (
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetSwapParamsQueryEncoding pins the exact query string produced for the
// classic-swap params — the same go-querystring encoding the HTTP executor uses
// (internal/http-executor/http.go). Aggregation is the most-used funds-critical
// swap path, so this guards its wire format: a changed or dropped `url` tag, or
// a type change that alters encoding, would fail here rather than silently
// change what the SDK sends to the API.
func TestGetSwapParamsQueryEncoding(t *testing.T) {
	tests := []struct {
		name     string
		params   GetSwapParams
		expected string
	}{
		{
			name: "full field set",
			params: GetSwapParams{
				Src:               "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				Dst:               "0x4200000000000000000000000000000000000006",
				Amount:            "1000000",
				From:              "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
				Origin:            "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
				Slippage:          1,
				Fee:               1.5,
				Receiver:          "0x1111111111111111111111111111111111111111",
				Permit:            "0xabcdef",
				IncludeTokensInfo: true,
			},
			expected: "amount=1000000" +
				"&dst=0x4200000000000000000000000000000000000006" +
				"&fee=1.5" +
				"&from=0x266E77cE9034a023056ea2845CB6A20517F6FDB7" +
				"&includeTokensInfo=true" +
				"&origin=0x266E77cE9034a023056ea2845CB6A20517F6FDB7" +
				"&permit=0xabcdef" +
				"&receiver=0x1111111111111111111111111111111111111111" +
				"&slippage=1" +
				"&src=0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		},
		{
			// The required fields encode; every omitempty field is absent. Guards
			// against an omitempty being accidentally dropped (which would send
			// zero-valued junk like fee=0 to the API).
			name: "required fields only omit the rest",
			params: GetSwapParams{
				Src:      "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				Dst:      "0x4200000000000000000000000000000000000006",
				Amount:   "1000000",
				From:     "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
				Origin:   "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
				Slippage: 1,
			},
			expected: "amount=1000000" +
				"&dst=0x4200000000000000000000000000000000000006" +
				"&from=0x266E77cE9034a023056ea2845CB6A20517F6FDB7" +
				"&origin=0x266E77cE9034a023056ea2845CB6A20517F6FDB7" +
				"&slippage=1" +
				"&src=0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			values, err := query.Values(tc.params)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, values.Encode())
		})
	}
}
