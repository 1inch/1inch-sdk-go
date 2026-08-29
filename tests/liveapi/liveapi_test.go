//go:build liveapi

// Package liveapi smoke-tests the SDK's decode paths against the production
// 1inch Dev Portal API. These tests are read-only (quotes and lookups; no
// orders, no transactions) and validate that live API responses still decode
// into the SDK's generated types — catching upstream behavior drift that spec
// comparison alone cannot see.
//
// Run with a Dev Portal token:
//
//	DEV_PORTAL_TOKEN=... go test -tags liveapi ./tests/liveapi/
//
// The spec-drift workflow runs these weekly.
package liveapi

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v5/constants"
	"github.com/1inch/1inch-sdk-go/v5/sdk-clients/fusionplus"
	"github.com/1inch/1inch-sdk-go/v5/sdk-clients/tokens"
)

const apiURL = "https://api.1inch.dev"

// throwawayKey is a well-known publicly-burned private key (anvil/hardhat dev
// account #0). It is used only to satisfy client constructors for read-only
// quote requests; nothing is signed or broadcast.
const throwawayKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

func devPortalToken(t *testing.T) string {
	t.Helper()
	token := os.Getenv("DEV_PORTAL_TOKEN")
	if token == "" {
		t.Skip("DEV_PORTAL_TOKEN not set")
	}
	return token
}

func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// TestTokensSearchDecodes validates that the live token API's responses decode
// into the corrected generated types (object tags, integer chain ids,
// value-typed optional fields).
func TestTokensSearchDecodes(t *testing.T) {
	config, err := tokens.NewConfiguration(tokens.ConfigurationParams{
		ChainId: constants.EthereumChainId,
		ApiUrl:  apiURL,
		ApiKey:  devPortalToken(t),
	})
	require.NoError(t, err)
	client, err := tokens.NewClient(config)
	require.NoError(t, err)

	results, err := client.SearchTokenAllChains(testContext(t), tokens.SearchAllChainsParams{
		Query: "1inch",
		Limit: 5,
	})
	require.NoError(t, err)
	require.NotEmpty(t, results, "search returned no tokens")

	for _, token := range results {
		assert.NotEmpty(t, token.Address)
		assert.Positivef(t, token.ChainId, "chain id must decode as a positive integer, got %d for %s", token.ChainId, token.Address)
	}
}

// TestFusionPlusQuoteDecodes validates that a live cross-chain quote decodes
// into the corrected generated types (string quoteId, string amounts and
// safety deposits, populated presets).
func TestFusionPlusQuoteDecodes(t *testing.T) {
	config, err := fusionplus.NewConfiguration(fusionplus.ConfigurationParams{
		ApiUrl:     apiURL,
		ApiKey:     devPortalToken(t),
		PrivateKey: throwawayKey,
	})
	require.NoError(t, err)
	client, err := fusionplus.NewClient(config)
	require.NoError(t, err)

	quote, err := client.GetQuote(testContext(t), fusionplus.QuoteParams{
		SrcChain:        constants.ArbitrumChainId,
		DstChain:        constants.BaseChainId,
		SrcTokenAddress: "0x82af49447d8a07e3bd95bd0d56f35241523fbab1", // WETH (Arbitrum)
		DstTokenAddress: "0x4200000000000000000000000000000000000006", // WETH (Base)
		Amount:          "100000000000000000",                         // 0.1 WETH
		WalletAddress:   "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", // throwawayKey's address
		EnableEstimate:  true,
	})
	require.NoError(t, err)

	assert.NotEmpty(t, quote.QuoteId, "quoteId must decode as a non-empty string")
	assert.NotEmpty(t, quote.SrcTokenAmount)
	assert.NotEmpty(t, quote.DstTokenAmount)
	assert.NotEmpty(t, quote.SrcSafetyDeposit)
	assert.NotEmpty(t, quote.RecommendedPreset)

	preset, err := fusionplus.GetPreset(quote.Presets, quote.RecommendedPreset)
	require.NoError(t, err)
	assert.NotEmpty(t, preset.AuctionStartAmount)
	assert.NotEmpty(t, preset.SecretsCount)
}
