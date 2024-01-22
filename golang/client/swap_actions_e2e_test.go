//go:build e2e
// +build e2e

package client

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/1inch/1inch-sdk/golang/client/onchain"
	"github.com/1inch/1inch-sdk/golang/client/swap"
	"github.com/1inch/1inch-sdk/golang/helpers"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/chains"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/contracts"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/tokens"
	"github.com/1inch/1inch-sdk/golang/helpers/consts/web3providers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestSwapTokensE2E(t *testing.T) {

	testcases := []struct {
		description         string
		config              Config
		swapParams          swap.AggregationControllerGetSwapParams
		removeApprovalAfter bool
		approvalType        swap.ApprovalType
		expectedOutput      string
	}{
		{
			description: "Polygon - Swap 0.01 DAI for USDC - Approval - Does not support traditional permit interface",
			config: Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				Web3HttpProvider: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
				ChainId:          chains.Polygon,
			},
			swapParams: swap.AggregationControllerGetSwapParams{
				Src:      tokens.PolygonDai,
				Dst:      tokens.PolygonUsdc,
				Amount:   "10000000000000000",
				From:     os.Getenv("WALLET_ADDRESS"),
				Slippage: 0.5,
			},
			removeApprovalAfter: true,
			approvalType:        swap.PermitIfPossible,
		},
		{
			description: "Polygon - Swap 0.01 FRAX for USDC - Approval - Forced",
			config: Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				Web3HttpProvider: os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
				ChainId:          chains.Polygon,
			},
			swapParams: swap.AggregationControllerGetSwapParams{
				Src:      tokens.PolygonFrax,
				Dst:      tokens.PolygonUsdc,
				Amount:   "10000000000000000",
				From:     os.Getenv("WALLET_ADDRESS"),
				Slippage: 0.5,
			},
			removeApprovalAfter: true,
			approvalType:        swap.ApprovalAlways,
		},
		{
			description: "Arbitrum - Swap 0.01 USDC for DAI - Approve - Arbitrum unsuported right now",
			config: Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				Web3HttpProvider: web3providers.Arbitrum,
				ChainId:          chains.Arbitrum,
			},
			swapParams: swap.AggregationControllerGetSwapParams{
				Src:      tokens.ArbitrumUsdc,
				Dst:      tokens.ArbitrumDai,
				Amount:   "10000",
				From:     os.Getenv("WALLET_ADDRESS"),
				Slippage: 0.5,
			},
			approvalType: swap.PermitIfPossible,
		},
		{
			description: "Arbitrum - Swap $0.01 worth of ETH for USDC - Approve - Native token for ERC20",
			config: Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				Web3HttpProvider: web3providers.Arbitrum,
				ChainId:          chains.Arbitrum,
			},
			swapParams: swap.AggregationControllerGetSwapParams{
				Src:      tokens.NativeToken,
				Dst:      tokens.ArbitrumUsdc,
				Amount:   "3974301376798",
				From:     os.Getenv("WALLET_ADDRESS"),
				Slippage: 0.5,
			},

			approvalType: swap.PermitIfPossible,
		},
		{
			description: "Ethereum - Swap $0.01 worth of 1inch for ETH - Force Permit1",
			config: Config{
				DevPortalApiKey:  os.Getenv("DEV_PORTAL_TOKEN"),
				WalletKey:        os.Getenv("WALLET_KEY"),
				Web3HttpProvider: web3providers.Ethereum,
				ChainId:          chains.Ethereum,
			},
			swapParams: swap.AggregationControllerGetSwapParams{
				Src:      tokens.Ethereum1inch,
				Dst:      tokens.NativeToken,
				Amount:   "20000000000000000",
				From:     os.Getenv("WALLET_ADDRESS"),
				Slippage: 0.5,
			},
			approvalType: swap.PermitAlways,
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%v", tc.description), func(t *testing.T) {

			t.Cleanup(func() {
				helpers.Sleep()
			})

			// Create the 1inch client
			c, err := NewClient(tc.config)
			require.NoError(t, err)

			// Swap tokens
			err = c.Actions.SwapTokens(tc.swapParams, true, tc.approvalType)
			if err != nil {
				log.Fatalf("Failed to swap tokens: %v", err)
			}
			require.NoError(t, err)
			if tc.removeApprovalAfter {

				allowance, err := onchain.ReadContractAllowance(c.EthClient, common.HexToAddress(tc.swapParams.Src), common.HexToAddress(tc.swapParams.From), common.HexToAddress(contracts.AggregationRouterV5))
				require.NoError(t, err)

				erc20Config := onchain.Erc20RevokeConfig{
					ChainId:                 tc.config.ChainId,
					Key:                     tc.config.WalletKey,
					Erc20Address:            common.HexToAddress(tc.swapParams.Src),
					PublicAddress:           common.HexToAddress(tc.swapParams.From),
					SpenderAddress:          common.HexToAddress(contracts.AggregationRouterV5),
					AllowanceDecreaseAmount: allowance,
				}
				err = onchain.RevokeApprovalForRouter(c.EthClient, c.NonceCache, erc20Config)
				require.NoError(t, err)
			}
		})
	}
}
