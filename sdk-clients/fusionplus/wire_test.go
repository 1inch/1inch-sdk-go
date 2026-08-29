package fusionplus

import (
	"encoding/json"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v5/constants"
)

// TestQuoteParamsQueryEncoding pins the exact query string produced for the
// quoter params — the same go-querystring encoding the HTTP executor uses.
// This guards the wire format across type changes: it would have caught both
// the Aurora chain-id float32 rounding (srcChain=1313161600) and the *big.Int
// fee being silently dropped from the query entirely.
func TestQuoteParamsQueryEncoding(t *testing.T) {
	tests := []struct {
		name     string
		params   QuoteParams
		expected string
	}{
		{
			name: "All fields including Aurora chain id and fee",
			params: QuoteParams{
				SrcChain:        constants.AuroraChainId,
				DstChain:        constants.BaseChainId,
				SrcTokenAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
				DstTokenAddress: "0x4200000000000000000000000000000000000006",
				Amount:          "12345678901234567890",
				WalletAddress:   "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
				EnableEstimate:  true,
				Fee:             100,
				IsPermit2:       true,
			},
			expected: "amount=12345678901234567890" +
				"&dstChain=8453" +
				"&dstTokenAddress=0x4200000000000000000000000000000000000006" +
				"&enableEstimate=true" +
				"&fee=100" +
				"&isPermit2=true" +
				"&srcChain=1313161554" +
				"&srcTokenAddress=0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" +
				"&walletAddress=0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
		},
		{
			name: "Optional fields omitted when zero",
			params: QuoteParams{
				SrcChain:        constants.EthereumChainId,
				DstChain:        constants.ArbitrumChainId,
				SrcTokenAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
				DstTokenAddress: "0x82af49447d8a07e3bd95bd0d56f35241523fbab1",
				Amount:          "1000000",
				WalletAddress:   "0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
			},
			expected: "amount=1000000" +
				"&dstChain=42161" +
				"&dstTokenAddress=0x82af49447d8a07e3bd95bd0d56f35241523fbab1" +
				"&enableEstimate=false" +
				"&srcChain=1" +
				"&srcTokenAddress=0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" +
				"&walletAddress=0x266E77cE9034a023056ea2845CB6A20517F6FDB7",
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

// TestSignedOrderInputJSONEncoding pins the JSON sent to the relayer,
// including exact integer chain ids above float32 precision.
func TestSignedOrderInputJSONEncoding(t *testing.T) {
	tests := []struct {
		name     string
		input    SignedOrderInput
		contains []string
	}{
		{
			name: "Aurora chain id survives exactly",
			input: SignedOrderInput{
				SrcChainId: constants.AuroraChainId,
				QuoteId:    "abc-123",
				Signature:  "0xsig",
			},
			contains: []string{
				`"srcChainId":1313161554`,
				`"quoteId":"abc-123"`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := json.Marshal(tc.input)
			require.NoError(t, err)
			for _, want := range tc.contains {
				assert.Contains(t, string(raw), want)
			}
		})
	}
}

// TestResponseDecoding pins how API responses decode into the corrected
// types: quoteId as a string, USD prices as strings, points as an array, and
// chain ids as exact integers.
func TestResponseDecoding(t *testing.T) {
	t.Run("Quote", func(t *testing.T) {
		payload := `{
			"quoteId": "9a43c86d-f3d7-45b9-8cb6-803d2bdb7a6b",
			"srcTokenAmount": "12345678901234567890",
			"dstTokenAmount": "9876543210",
			"srcEscrowFactory": "0x1111111111111111111111111111111111111111",
			"dstEscrowFactory": "0x2222222222222222222222222222222222222222",
			"srcSafetyDeposit": "1000000000000000",
			"dstSafetyDeposit": "2000000000000000",
			"recommendedPreset": "fast",
			"whitelist": ["0x3333333333333333333333333333333333333333"]
		}`
		var out Quote
		require.NoError(t, json.Unmarshal([]byte(payload), &out))
		assert.Equal(t, "9a43c86d-f3d7-45b9-8cb6-803d2bdb7a6b", out.QuoteId)
		assert.Equal(t, "12345678901234567890", out.SrcTokenAmount)
		assert.Equal(t, "1000000000000000", out.SrcSafetyDeposit)
		assert.Equal(t, GetQuoteOutputRecommendedPreset("fast"), out.RecommendedPreset)
	})

	t.Run("GetOrderFillsByHashOutput", func(t *testing.T) {
		payload := `{
			"orderHash": "0xdeadbeef",
			"srcChainId": 1313161554,
			"dstChainId": 8453,
			"srcTokenPriceUsd": "1.0002",
			"dstTokenPriceUsd": "4321.99",
			"points": [{"coefficient": 100, "delay": 12}],
			"status": "executed",
			"validation": "valid"
		}`
		var out GetOrderFillsByHashOutput
		require.NoError(t, json.Unmarshal([]byte(payload), &out))
		assert.Equal(t, constants.AuroraChainId, out.SrcChainId)
		assert.Equal(t, constants.BaseChainId, out.DstChainId)
		assert.Equal(t, "1.0002", out.SrcTokenPriceUsd)
		assert.Equal(t, "4321.99", out.DstTokenPriceUsd)
		require.Len(t, out.Points, 1)
		assert.Equal(t, float32(100), out.Points[0].Coefficient)
	})
}
