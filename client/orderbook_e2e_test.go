package client

// TODO commenting this test file out while permits are disabled on orderbook
//import (
//	"context"
//	"fmt"
//	"os"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//
//	"github.com/1inch/1inch-sdk-go/client/models"
//	"github.com/1inch/1inch-sdk-go/helpers"
//	"github.com/1inch/1inch-sdk-go/helpers/consts/addresses"
//	"github.com/1inch/1inch-sdk-go/helpers/consts/amounts"
//	"github.com/1inch/1inch-sdk-go/helpers/consts/chains"
//	"github.com/1inch/1inch-sdk-go/helpers/consts/tokens"
//	"github.com/1inch/1inch-sdk-go/helpers/consts/web3providers"
//	"github.com/1inch/1inch-sdk-go/internal/onchain"
//	"github.com/1inch/1inch-sdk-go/internal/tenderly"
//)
//
//func TestCreateOrderE2E(t *testing.T) {
//
//	testcases := []struct {
//		description       string
//		config            models.ClientConfig
//		createOrderParams models.CreateOrderParams
//		expectedOutput    string
//	}{
//		{
//			description: "Arbitrum - Create limit order offering 1 FRAX for 1 DAI",
//			config: models.ClientConfig{
//				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
//				Web3HttpProviders: []models.Web3Provider{
//					{
//						ChainId: chains.Arbitrum,
//						Url:     web3providers.Arbitrum,
//					},
//				},
//			},
//			createOrderParams: models.CreateOrderParams{
//				ChainId:      chains.Arbitrum,
//				PrivateKey:   os.Getenv("WALLET_KEY_EMPTY"),
//				Maker:        os.Getenv("WALLET_ADDRESS_EMPTY"),
//				MakerAsset:   tokens.ArbitrumFrax,
//				TakerAsset:   tokens.ArbitrumDai,
//				MakingAmount: amounts.Ten18,
//				TakingAmount: amounts.Ten18,
//				Taker:        addresses.Zero,
//			},
//		},
//		{
//			description: "Polygon - Create limit order offering 1 FRAX for 1 DAI",
//			config: models.ClientConfig{
//				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
//				Web3HttpProviders: []models.Web3Provider{
//					{
//						ChainId: chains.Polygon,
//						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY_POLYGON"),
//					},
//				},
//			},
//			createOrderParams: models.CreateOrderParams{
//				ChainId:      chains.Polygon,
//				PrivateKey:   os.Getenv("WALLET_KEY_EMPTY"),
//				Maker:        os.Getenv("WALLET_ADDRESS_EMPTY"),
//				MakerAsset:   tokens.PolygonFrax,
//				TakerAsset:   tokens.PolygonDai,
//				MakingAmount: amounts.Ten18,
//				TakingAmount: amounts.Ten18,
//				Taker:        addresses.Zero,
//			},
//		},
//		{
//			description: "Ethereum - Create limit order offering 1 1INCH for 1 DAI",
//			config: models.ClientConfig{
//				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
//				Web3HttpProviders: []models.Web3Provider{
//					{
//						ChainId: chains.Ethereum,
//						Url:     os.Getenv("WEB_3_HTTP_PROVIDER_URL_WITH_KEY"),
//					},
//				},
//			},
//			createOrderParams: models.CreateOrderParams{
//				ChainId:      chains.Ethereum,
//				PrivateKey:   os.Getenv("WALLET_KEY_EMPTY"),
//				Maker:        os.Getenv("WALLET_ADDRESS_EMPTY"),
//				MakerAsset:   tokens.Ethereum1inch,
//				TakerAsset:   tokens.EthereumDai,
//				MakingAmount: amounts.Ten18,
//				TakingAmount: amounts.Ten18,
//				Taker:        addresses.Zero,
//				ApprovalType: onchain.PermitAlways,
//			},
//		},
//		{
//			description: "BSC - Create limit order offering 1 USDC for 1 DAI",
//			config: models.ClientConfig{
//				DevPortalApiKey: os.Getenv("DEV_PORTAL_TOKEN"),
//				Web3HttpProviders: []models.Web3Provider{
//					{
//						ChainId: chains.Bsc,
//						Url:     web3providers.Bsc,
//					},
//				},
//			},
//			createOrderParams: models.CreateOrderParams{
//				ApprovalType: onchain.PermitAlways,
//				ChainId:      chains.Bsc,
//				PrivateKey:   os.Getenv("WALLET_KEY_EMPTY"),
//				Maker:        os.Getenv("WALLET_ADDRESS_EMPTY"),
//				MakerAsset:   tokens.BscFrax,
//				TakerAsset:   tokens.BscDai,
//				MakingAmount: amounts.Ten18,
//				TakingAmount: amounts.Ten18,
//				Taker:        addresses.Zero,
//			},
//		},
//	}
//
//	//TODO set this up to have some form of configurations that enable the tests to run onchain and should also cleanup any previous test runs
//	tenderlyApiKey := os.Getenv("TENDERLY_API_KEY")
//	if tenderlyApiKey == "" {
//		fmt.Printf("No Tenderly API key present in environment, skipping e2e tests")
//		return
//	}
//
//	for _, tc := range testcases {
//		t.Run(tc.description, func(t *testing.T) {
//
//			t.Cleanup(func() {
//				helpers.Sleep()
//			})
//
//			c, err := NewClient(tc.config)
//			require.NoError(t, err)
//
//			ctx := context.Background()
//			if tenderlyApiKey != "" {
//				ctx = context.WithValue(ctx, tenderly.SwapConfigKey, tenderly.SimulationConfig{
//					TenderlyApiKey: tenderlyApiKey,
//				})
//			}
//			_, _, err = c.OrderbookApi.CreateOrder(ctx, tc.createOrderParams)
//			require.NoError(t, err)
//		})
//	}
//}
