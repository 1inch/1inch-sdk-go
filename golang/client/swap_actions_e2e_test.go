package client

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/client/tenderly"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/web3providers"
	"github.com/stretchr/testify/require"
)

func TestSwapTokensTenderlyE2E(t *testing.T) {

	testcases := []struct {
		description         string
		tenderlyDescription string
		config              Config
		swapParams          swap.SwapTokensParams
		approvalType        swap.ApprovalType
		expectedOutput      string
	}{
		{
			description:         "Polygon - Swap 0.01 DAI for USDC - Approval - Does not support traditional permit interface",
			tenderlyDescription: "DP-DAI->USDC-Approval",
			config: Config{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Polygon,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
					},
				},
				TenderlyKey: os.Getenv("TENDERLY_API_KEY"),
			},
			swapParams: swap.SwapTokensParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.PolygonDai,
					Dst:      tokens.PolygonUsdc,
					Amount:   "10000000000000000",
					From:     os.Getenv("WALLET_ADDRESS"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Polygon,
				WalletKey:     os.Getenv("WALLET_KEY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS"),
				ApprovalType:  swap.PermitIfPossible,
			},
		},
		{
			description:         "Polygon - Swap 0.01 FRAX for USDC - Approval - Forced",
			tenderlyDescription: "DP-FRAX->USDC-Approval",
			config: Config{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Polygon,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
					},
				},
				TenderlyKey: os.Getenv("TENDERLY_API_KEY"),
			},
			swapParams: swap.SwapTokensParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.PolygonFrax,
					Dst:      tokens.PolygonUsdc,
					Amount:   "10000000000000000",
					From:     os.Getenv("WALLET_ADDRESS"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Polygon,
				WalletKey:     os.Getenv("WALLET_KEY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS"),
				ApprovalType:  swap.ApprovalAlways,
			},
		},
		{
			description:         "Polygon - Swap 0.01 FRAX for USDC - Permit",
			tenderlyDescription: "DP-FRAX->USDC-Permit",
			config: Config{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Polygon,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
					},
				},
				TenderlyKey: os.Getenv("TENDERLY_API_KEY"),
			},
			swapParams: swap.SwapTokensParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.PolygonFrax,
					Dst:      tokens.PolygonUsdc,
					Amount:   "10000000000000000",
					From:     os.Getenv("WALLET_ADDRESS"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Polygon,
				WalletKey:     os.Getenv("WALLET_KEY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS"),
				ApprovalType:  swap.PermitIfPossible,
			},
		},
		{
			description:         "Arbitrum - Swap 0.01 USDC for DAI - Approve - Arbitrum unsupported right now",
			tenderlyDescription: "DP-USDC->DAI-Approve",
			config: Config{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Arbitrum,
						Url:     web3providers.Arbitrum,
					},
				},
				TenderlyKey: os.Getenv("TENDERLY_API_KEY"),
			},
			swapParams: swap.SwapTokensParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.ArbitrumUsdc,
					Dst:      tokens.ArbitrumDai,
					Amount:   "10000",
					From:     os.Getenv("WALLET_ADDRESS"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Arbitrum,
				WalletKey:     os.Getenv("WALLET_KEY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS"),
				ApprovalType:  swap.PermitIfPossible,
			},
		},
		{
			description:         "Arbitrum - Swap $0.01 worth of ETH for USDC - Native token for ERC20",
			tenderlyDescription: "DP-ETH->USDC",
			config: Config{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Arbitrum,
						Url:     web3providers.Arbitrum,
					},
				},
				TenderlyKey: os.Getenv("TENDERLY_API_KEY"),
			},
			swapParams: swap.SwapTokensParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.NativeToken,
					Dst:      tokens.ArbitrumUsdc,
					Amount:   "3974301376798",
					From:     os.Getenv("WALLET_ADDRESS"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Arbitrum,
				WalletKey:     os.Getenv("WALLET_KEY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS"),
				ApprovalType:  swap.PermitIfPossible,
			},
		},
		{
			description:         "Ethereum - Swap $0.01 worth of 1inch for ETH - Force Permit1",
			tenderlyDescription: "DP-1inch->ETH-Permit1",
			config: Config{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []Web3ProviderConfig{
					{
						ChainId: chains.Ethereum,
						Url:     web3providers.Ethereum,
					},
				},
				TenderlyKey: os.Getenv("TENDERLY_API_KEY"),
			},
			swapParams: swap.SwapTokensParams{
				AggregationControllerGetSwapParams: swap.AggregationControllerGetSwapParams{
					Src:      tokens.Ethereum1inch,
					Dst:      tokens.NativeToken,
					Amount:   "20000000000000000",
					From:     os.Getenv("WALLET_ADDRESS"),
					Slippage: 0.5,
				},
				ChainId:       chains.Ethereum,
				SkipWarnings:  true,
				WalletKey:     os.Getenv("WALLET_KEY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS"),
				ApprovalType:  swap.PermitAlways,
			},
		},
	}

	err := cleanupForksFromPreviousTests(os.Getenv("TENDERLY_API_KEY"))
	require.NoError(t, err, fmt.Errorf("failed to delete forks from previous test runs: %v", err))

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			// Create the 1inch client
			c, err := NewClient(tc.config)
			require.NoError(t, err)

			// Swap tokens
			err = c.Actions.SwapTokens(tc.swapParams)
			if err != nil {
				log.Fatalf("Failed to swap tokens: %v", err)
			}
			require.NoError(t, err)
		})
	}
}

func cleanupForksFromPreviousTests(tenderlyApiKey string) error {

	forksResponse, err := tenderly.GetTenderlyForks(tenderlyApiKey)
	if err != nil {
		return fmt.Errorf("failed to get tenderly forks: %v", err)
	}

	for _, fork := range forksResponse.Forks {
		if strings.Contains(fork.Alias, "Dev Portal") {
			err := tenderly.DeleteTenderlyFork(tenderlyApiKey, fork.ID)
			if err != nil {
				return fmt.Errorf("failed to delete tenderly fork: %v", err)
			}
		}
	}
	return nil
}
