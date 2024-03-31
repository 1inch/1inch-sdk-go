package multicall

import (
	"context"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/abis"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/chains"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/tokens"
	"github.com/1inch/1inch-sdk-go/internal/helpers/consts/web3providers"
)

const (
	nameMethod   = "name"
	symbolMethod = "symbol"
)

func TestMulticallEthereumSuccess(t *testing.T) {
	var callData []CallData

	client, err := ethclient.Dial(web3providers.Ethereum)
	require.NoError(t, err)

	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
	require.NoError(t, err)

	nameRequest, err := parsedABI.Pack(nameMethod)
	require.NoError(t, err)
	callData = append(callData, BuildCallData(tokens.EthereumUsdc, hexutil.Encode(nameRequest), 0))
	callData = append(callData, BuildCallData(tokens.EthereumDai, hexutil.Encode(nameRequest), 0))

	resp, err := MultiCall(context.Background(), MulticallParams{
		Client:   client,
		ChainId:  chains.Ethereum,
		Calldata: callData,
	})
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
	var callData []CallData

	client, err := ethclient.Dial(web3providers.Polygon)
	require.NoError(t, err)

	parsedABI, err := abi.JSON(strings.NewReader(abis.Erc20)) // Make a generic version of this ABI
	require.NoError(t, err)

	symbolRequest, err := parsedABI.Pack(symbolMethod)
	require.NoError(t, err)
	callData = append(callData, BuildCallData(tokens.PolygonUsdc, hexutil.Encode(symbolRequest), 0))
	callData = append(callData, BuildCallData(tokens.PolygonDai, hexutil.Encode(symbolRequest), 0))

	resp, err := MultiCall(context.Background(), MulticallParams{
		Client:   client,
		ChainId:  chains.Polygon,
		Calldata: callData,
	})
	require.NoError(t, err)

	var tokenName string
	err = parsedABI.UnpackIntoInterface(&tokenName, symbolMethod, resp[0])
	require.NoError(t, err)
	require.Equal(t, "USDC", tokenName)

	err = parsedABI.UnpackIntoInterface(&tokenName, symbolMethod, resp[1])
	require.NoError(t, err)
	require.Equal(t, "DAI", tokenName)
}
