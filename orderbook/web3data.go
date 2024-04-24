package orderbook

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/constants"
	"github.com/1inch/1inch-sdk-go/orderbook/models"
)

func (c *Client) GetSeriesNonce(ctx context.Context, publicAddress gethCommon.Address) (*big.Int, error) {

	seriesNonceManager, err := constants.GetSeriesNonceManagerFromChainId(int(c.Wallet.ChainId()))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get series nonce manager address: %v", err))
	}

	function := "nonce"

	seriesNonceManagerABI, err := abi.JSON(strings.NewReader(constants.SeriesNonceManagerABI)) // Make a generic version of this ABI
	if err != nil {
		return nil, err
	}

	seriesNonceData, err := seriesNonceManagerABI.Pack(function, big.NewInt(0), publicAddress)
	if err != nil {
		return nil, err
	}

	result, err := c.Wallet.Call(ctx, gethCommon.HexToAddress(seriesNonceManager), seriesNonceData)
	if err != nil {
		return nil, err
	}

	var nonce *big.Int
	err = seriesNonceManagerABI.UnpackIntoInterface(&nonce, function, result)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func (c *Client) GetFillOrderCalldata(getOrderResponse *models.GetOrderByHashResponseExtended) ([]byte, error) {
	function := "fillOrder"

	aggregationRouterV6, err := abi.JSON(strings.NewReader(constants.AggregationRouterV6ABI))
	if err != nil {
		return nil, err
	}

	compressedSignature, err := CompressSignature(getOrderResponse.Signature[2:])
	if err != nil {
		return nil, err
	}

	fillOrderData, err := aggregationRouterV6.Pack(function, getOrderResponse.LimitOrderDataNormalized, BytesToBytes32(compressedSignature.R), BytesToBytes32(compressedSignature.VS), getOrderResponse.LimitOrderDataNormalized.TakingAmount, big.NewInt(0))
	if err != nil {
		return nil, err
	}

	return fillOrderData, nil
}
