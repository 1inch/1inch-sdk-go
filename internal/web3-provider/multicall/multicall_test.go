package multicall

import (
	"context"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/constants"
)

const (
	nameMethod   = "name"
	symbolMethod = "symbol"

	Arbitrum = "https://arb1.arbitrum.io/rpc"
	Aurora   = "https://mainnet.aurora.dev"
	Ethereum = "https://eth-mainnet.public.blastapi.io"
	Polygon  = "https://polygon.llamarpc.com"
	Bsc      = "https://bsc-dataseed1.binance.org"

	NativeToken = "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"

	AuroraFrax = "0xDA2585430fEf327aD8ee44Af8F1f989a2A91A3d2"
	AuroraRose = "0xdcD6D4e2B3e1D1E1E6Fa8C21C8A323DcbecfF970"

	EthereumWeth  = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	EthereumUsdc  = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	EthereumFrax  = "0x853d955aCEf822Db058eb8505911ED77F175b99e"
	EthereumDai   = "0x6B175474E89094C44Da98b954EedeAC495271d0F"
	Ethereum1inch = "0x111111111117dC0aa78b770fA6A738034120C302"

	PolygonDai  = "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063"
	PolygonFrax = "0x45c32fA6DF82ead1e2EF74d17b76547EDdFaFF89"
	PolygonWeth = "0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"
	PolygonUsdc = "0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359"

	BscUsdc = "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d"
	BscDai  = "0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3"
	BscFrax = "0x90C97F71E18723b0Cf0dfa30ee176Ab653E89F40"

	ArbitrumUsdc = "0xaf88d065e77c8cc2239327c5edb3a432268e5831"
	ArbitrumDai  = "0xda10009cbd5d07dd0cecc66161fc93d7c9000da1"
	ArbitrumFrax = "0x17fc002b466eec40dae837fc4be5c67993ddbd6f"
)

func TestMulticallEthereumSuccess(t *testing.T) {
	client, err := ethclient.Dial(Ethereum)
	require.NoError(t, err)

	instance, err := NewMulticall(client, constants.EthereumChainId)
	require.NoError(t, err)

	var callData []CallData
	parsedABI, err := abi.JSON(strings.NewReader(constants.Erc20ABI)) // Make a generic version of this ABI
	require.NoError(t, err)

	nameRequest, err := parsedABI.Pack(nameMethod)
	require.NoError(t, err)
	callData = append(callData, BuildCallData(EthereumUsdc, hexutil.Encode(nameRequest), 0))
	callData = append(callData, BuildCallData(EthereumDai, hexutil.Encode(nameRequest), 0))

	resp, err := instance.Execute(context.Background(), callData)
	require.NoError(t, err)

	var tokenName string
	err = parsedABI.UnpackIntoInterface(&tokenName, nameMethod, resp[0])
	require.NoError(t, err)
	require.Equal(t, "USD Coin", tokenName)

	err = parsedABI.UnpackIntoInterface(&tokenName, nameMethod, resp[1])
	require.NoError(t, err)
	require.Equal(t, "Dai Stablecoin", tokenName)
}

func TestMulticallPolygonSuccess(t *testing.T) {
	client, err := ethclient.Dial(Polygon)
	require.NoError(t, err)

	instance, err := NewMulticall(client, constants.PolygonChainId)
	require.NoError(t, err)

	var callData []CallData

	parsedABI, err := abi.JSON(strings.NewReader(constants.Erc20ABI)) // Make a generic version of this ABI
	require.NoError(t, err)

	symbolRequest, err := parsedABI.Pack(symbolMethod)
	require.NoError(t, err)
	callData = append(callData, BuildCallData(PolygonUsdc, hexutil.Encode(symbolRequest), 0))
	callData = append(callData, BuildCallData(PolygonDai, hexutil.Encode(symbolRequest), 0))

	resp, err := instance.Execute(context.Background(), callData)
	require.NoError(t, err)

	var tokenName string
	err = parsedABI.UnpackIntoInterface(&tokenName, symbolMethod, resp[0])
	require.NoError(t, err)
	require.Equal(t, "USDC", tokenName)

	err = parsedABI.UnpackIntoInterface(&tokenName, symbolMethod, resp[1])
	require.NoError(t, err)
	require.Equal(t, "DAI", tokenName)
}
