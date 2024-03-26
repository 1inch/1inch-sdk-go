package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/client/models"
	"github.com/1inch/1inch-sdk-go/helpers"
	"github.com/1inch/1inch-sdk-go/helpers/consts/amounts"
	"github.com/1inch/1inch-sdk-go/helpers/consts/chains"
	"github.com/1inch/1inch-sdk-go/helpers/consts/tokens"
	"github.com/1inch/1inch-sdk-go/helpers/consts/web3providers"
	"github.com/1inch/1inch-sdk-go/internal/onchain"
	"github.com/1inch/1inch-sdk-go/internal/tenderly"
)

const ten18Hex = "0x0000000000000000000000000000000000000000000000000DE0B6B3A7640000"

func TestSwapTokensTenderlyE2E(t *testing.T) {

	testcases := []struct {
		description         string
		tenderlyDescription string
		config              models.ClientConfig
		swapParams          models.SwapTokensParams
		stateOverrides      map[string]tenderly.StateObject
		approvalType        onchain.ApprovalType
		expectedOutput      string
	}{
		{
			description:         "Polygon - Swap 0.01 DAI for USDC - Approval - Does not support traditional permit interface",
			tenderlyDescription: "DP-DAI->USDC-Approval",
			config: models.ClientConfig{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []models.Web3Provider{
					{
						ChainId: chains.Polygon,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
					},
				},
			},
			swapParams: models.SwapTokensParams{
				AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
					Src:      tokens.PolygonDai,
					Dst:      tokens.PolygonUsdc,
					Amount:   "10000000000000000",
					From:     os.Getenv("WALLET_ADDRESS_EMPTY"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Polygon,
				WalletKey:     os.Getenv("WALLET_KEY_EMPTY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS_EMPTY"),
				ApprovalType:  onchain.PermitIfPossible,
			},
			stateOverrides: map[string]tenderly.StateObject{
				os.Getenv("WALLET_ADDRESS_EMPTY"): {
					Balance: amounts.Ten18,
				},
				tokens.PolygonDai: {
					Storage: map[string]string{
						tenderly.GetStorageSlotHash(os.Getenv("WALLET_ADDRESS_EMPTY"), 0): ten18Hex,
					},
				},
			},
		},
		{
			description:         "Polygon - Swap 0.01 USDC for DAI - Permit - Contract has a version value of 2",
			tenderlyDescription: "DP-USDC->DAI-Approval",
			config: models.ClientConfig{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []models.Web3Provider{
					{
						ChainId: chains.Polygon,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
					},
				},
			},
			swapParams: models.SwapTokensParams{
				AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
					Src:      tokens.PolygonUsdc,
					Dst:      tokens.PolygonDai,
					Amount:   amounts.Ten6,
					From:     os.Getenv("WALLET_ADDRESS_EMPTY"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Polygon,
				WalletKey:     os.Getenv("WALLET_KEY_EMPTY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS_EMPTY"),
				ApprovalType:  onchain.PermitIfPossible,
			},
			stateOverrides: map[string]tenderly.StateObject{
				os.Getenv("WALLET_ADDRESS_EMPTY"): {
					Balance: amounts.Ten18,
				},
				tokens.PolygonUsdc: {
					Storage: map[string]string{
						tenderly.GetStorageSlotHash(os.Getenv("WALLET_ADDRESS_EMPTY"), 9): ten18Hex,
					},
				},
			},
		},
		{
			description:         "Polygon - Swap 0.01 FRAX for USDC - Approval - Forced",
			tenderlyDescription: "DP-FRAX->USDC-Approval",
			config: models.ClientConfig{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []models.Web3Provider{
					{
						ChainId: chains.Polygon,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
					},
				},
			},
			swapParams: models.SwapTokensParams{
				AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
					Src:      tokens.PolygonFrax,
					Dst:      tokens.PolygonUsdc,
					Amount:   "10000000000000000",
					From:     os.Getenv("WALLET_ADDRESS_EMPTY"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Polygon,
				WalletKey:     os.Getenv("WALLET_KEY_EMPTY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS_EMPTY"),
				ApprovalType:  onchain.ApprovalAlways,
			},
			stateOverrides: map[string]tenderly.StateObject{
				os.Getenv("WALLET_ADDRESS_EMPTY"): {
					Balance: amounts.Ten18,
				},
				tokens.PolygonFrax: {
					Storage: map[string]string{
						tenderly.GetStorageSlotHash(os.Getenv("WALLET_ADDRESS_EMPTY"), 0): ten18Hex,
					},
				},
			},
		},
		{
			description:         "Polygon - Swap 0.01 FRAX for USDC - Permit",
			tenderlyDescription: "DP-FRAX->USDC-Permit",
			config: models.ClientConfig{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []models.Web3Provider{
					{
						ChainId: chains.Polygon,
						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
					},
				},
			},
			swapParams: models.SwapTokensParams{
				AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
					Src:      tokens.PolygonFrax,
					Dst:      tokens.PolygonUsdc,
					Amount:   "10000000000000000",
					From:     os.Getenv("WALLET_ADDRESS_EMPTY"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Polygon,
				WalletKey:     os.Getenv("WALLET_KEY_EMPTY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS_EMPTY"),
				ApprovalType:  onchain.PermitIfPossible,
			},
			stateOverrides: map[string]tenderly.StateObject{
				os.Getenv("WALLET_ADDRESS_EMPTY"): {
					Balance: amounts.Ten18,
				},
				tokens.PolygonFrax: {
					Storage: map[string]string{
						tenderly.GetStorageSlotHash(os.Getenv("WALLET_ADDRESS_EMPTY"), 0): ten18Hex,
					},
				},
			},
		},
		{
			description:         "Arbitrum - Swap 0.01 USDC for DAI - Approve - Arbitrum unsupported right now",
			tenderlyDescription: "DP-USDC->DAI-Approve",
			config: models.ClientConfig{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []models.Web3Provider{
					{
						ChainId: chains.Arbitrum,
						Url:     web3providers.Arbitrum,
					},
				},
			},
			swapParams: models.SwapTokensParams{
				AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
					Src:      tokens.ArbitrumUsdc,
					Dst:      tokens.ArbitrumDai,
					Amount:   "10000",
					From:     os.Getenv("WALLET_ADDRESS_EMPTY"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Arbitrum,
				WalletKey:     os.Getenv("WALLET_KEY_EMPTY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS_EMPTY"),
				ApprovalType:  onchain.PermitIfPossible,
			},
			stateOverrides: map[string]tenderly.StateObject{
				os.Getenv("WALLET_ADDRESS_EMPTY"): {
					Balance: amounts.Ten18,
				},
				tokens.ArbitrumUsdc: {
					Storage: map[string]string{
						tenderly.GetStorageSlotHash(os.Getenv("WALLET_ADDRESS_EMPTY"), 9): ten18Hex,
					},
				},
			},
		},
		{
			description:         "Arbitrum - Swap $0.01 worth of ETH for USDC - Native token for ERC20",
			tenderlyDescription: "DP-ETH->USDC",
			config: models.ClientConfig{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []models.Web3Provider{
					{
						ChainId: chains.Arbitrum,
						Url:     web3providers.Arbitrum,
					},
				},
			},
			swapParams: models.SwapTokensParams{
				AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
					Src:      tokens.NativeToken,
					Dst:      tokens.ArbitrumUsdc,
					Amount:   "3974301376798",
					From:     os.Getenv("WALLET_ADDRESS_EMPTY"),
					Slippage: 0.5,
				},
				SkipWarnings:  true,
				ChainId:       chains.Arbitrum,
				WalletKey:     os.Getenv("WALLET_KEY_EMPTY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS_EMPTY"),
				ApprovalType:  onchain.PermitIfPossible,
			},
			stateOverrides: map[string]tenderly.StateObject{
				os.Getenv("WALLET_ADDRESS_EMPTY"): {
					Balance: amounts.Ten18,
				},
				tokens.NativeToken: {
					Storage: map[string]string{
						tenderly.GetStorageSlotHash(os.Getenv("WALLET_ADDRESS_EMPTY"), 0): ten18Hex,
					},
				},
			},
		},
		{
			description:         "Ethereum - Swap $0.01 worth of 1inch for ETH - Force Permit1",
			tenderlyDescription: "DP-1inch->ETH-Permit1",
			config: models.ClientConfig{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []models.Web3Provider{
					{
						ChainId: chains.Ethereum,
						Url:     web3providers.Ethereum,
					},
				},
			},
			swapParams: models.SwapTokensParams{
				AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
					Src:      tokens.Ethereum1inch,
					Dst:      tokens.NativeToken,
					Amount:   "20000000000000000",
					From:     os.Getenv("WALLET_ADDRESS_EMPTY"),
					Slippage: 0.5,
				},
				ChainId:       chains.Ethereum,
				SkipWarnings:  true,
				WalletKey:     os.Getenv("WALLET_KEY_EMPTY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS_EMPTY"),
				ApprovalType:  onchain.PermitAlways,
			},
			stateOverrides: map[string]tenderly.StateObject{
				os.Getenv("WALLET_ADDRESS_EMPTY"): {
					Balance: amounts.Ten18,
				},
				tokens.Ethereum1inch: {
					Storage: map[string]string{
						tenderly.GetStorageSlotHash(os.Getenv("WALLET_ADDRESS_EMPTY"), 0): ten18Hex,
					},
				},
			},
		},
		{
			description:         "Ethereum - Swap $0.01 worth of USDC for ETH - Version 2 - Permit1",
			tenderlyDescription: "DP-1inch->ETH-Permit1",
			config: models.ClientConfig{
				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
				Web3HttpProviders: []models.Web3Provider{
					{
						ChainId: chains.Ethereum,
						Url:     web3providers.Ethereum,
					},
				},
			},
			swapParams: models.SwapTokensParams{
				AggregationControllerGetSwapParams: models.AggregationControllerGetSwapParams{
					Src:      tokens.EthereumUsdc,
					Dst:      tokens.NativeToken,
					Amount:   amounts.Ten6,
					From:     os.Getenv("WALLET_ADDRESS_EMPTY"),
					Slippage: 0.5,
				},
				ChainId:       chains.Ethereum,
				SkipWarnings:  true,
				WalletKey:     os.Getenv("WALLET_KEY_EMPTY"),
				PublicAddress: os.Getenv("WALLET_ADDRESS_EMPTY"),
				ApprovalType:  onchain.PermitAlways,
			},
			stateOverrides: map[string]tenderly.StateObject{
				os.Getenv("WALLET_ADDRESS_EMPTY"): {
					Balance: amounts.Ten18,
				},
				tokens.EthereumUsdc: {
					Storage: map[string]string{
						tenderly.GetStorageSlotHash(os.Getenv("WALLET_ADDRESS_EMPTY"), 9): ten18Hex,
					},
				},
			},
		},
	}

	//TODO set this up to have some form of configurations that enable the tests to run onchain
	tenderlyApiKey := os.Getenv("TENDERLY_API_KEY")
	if tenderlyApiKey == "" {
		fmt.Printf("No Tenderly API key present in environment, skipping e2e tests")
		return
	}
	err := cleanupForksFromPreviousTests(tenderlyApiKey)
	require.NoError(t, err, fmt.Errorf("failed to delete forks from previous test runs: %v", err))

	for _, tc := range testcases {
		t.Run(tc.description, func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			// Create the 1inch client
			c, err := NewClient(tc.config)
			require.NoError(t, err)

			ctx := context.Background()
			if tenderlyApiKey != "" {
				ctx = context.WithValue(ctx, tenderly.SwapConfigKey, tenderly.SimulationConfig{
					TenderlyApiKey: tenderlyApiKey,
					OverridesMap:   tc.stateOverrides,
				})
			}

			// Swap tokens
			err = c.Actions.swapTokens(ctx, tc.swapParams)
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
		if strings.HasPrefix(fork.Alias, "DP") {
			err := tenderly.DeleteTenderlyFork(tenderlyApiKey, fork.ID)
			if err != nil {
				return fmt.Errorf("failed to delete tenderly fork: %v", err)
			}
		}
	}
	return nil
}
