package orderbook

import (
	"context"
	"fmt"
	"log"
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/constants"
)

func (c *Client) GetSeriesNonce(ctx context.Context, publicAddress gethCommon.Address) (*big.Int, error) {

	seriesNonceManager, err := constants.GetSeriesNonceManagerFromChainId(int(c.Wallet.ChainId()))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get series nonce manager address: %v", err))
	}

	function := "nonce"

	seriesNonceData, err := c.SeriesNonceManager.Pack(function, big.NewInt(0), publicAddress)
	if err != nil {
		return nil, err
	}

	result, err := c.Wallet.Call(ctx, gethCommon.HexToAddress(seriesNonceManager), seriesNonceData)
	if err != nil {
		return nil, err
	}

	var nonce *big.Int
	err = c.SeriesNonceManager.UnpackIntoInterface(&nonce, function, result)
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

const makerTraitsPermitOnly = "7440945280133576583328096164017418065923851860621198004784596428783616"

func (c *Client) GetFillOrderCalldata(getOrderResponse *GetOrderByHashResponseExtended, takerTraits *TakerTraits) ([]byte, error) {

	var function string
	if getOrderResponse.Data.Extension == "0x" {
		function = "fillOrder"
	} else {
		if takerTraits == nil {
			return nil, fmt.Errorf("this order has extension data, but no taker traits were provided")
		}

		function = "fillOrderArgs"
	}

	compressedSignature, err := CompressSignature(getOrderResponse.Signature[2:])
	if err != nil {
		return nil, err
	}

	rCompressed, err := bytesToBytes32(compressedSignature.R)
	if err != nil {
		return nil, err
	}

	vsCompressed, err := bytesToBytes32(compressedSignature.VS)
	if err != nil {
		return nil, err
	}

	var fillOrderData []byte

	switch function {
	case "fillOrder":
		fillOrderData, err = c.AggregationRouterV6.Pack(function, getOrderResponse.LimitOrderDataNormalized, rCompressed, vsCompressed, getOrderResponse.LimitOrderDataNormalized.TakingAmount, big.NewInt(0))
		if err != nil {
			return nil, err
		}
	case "fillOrderArgs":
		takerTraitsEncoded := takerTraits.Encode()
		fillOrderData, err = c.AggregationRouterV6.Pack(function, getOrderResponse.LimitOrderDataNormalized, rCompressed, vsCompressed, getOrderResponse.LimitOrderDataNormalized.TakingAmount, takerTraitsEncoded.TraitFlags, takerTraitsEncoded.Args)
		if err != nil {
			return nil, err
		}
	}

	return fillOrderData, nil
}

// bytesToBytes32 converts a byte slice to a [32]byte, padding with zeros if necessary,
// and truncating if it's too long.
func bytesToBytes32(b []byte) (*[32]byte, error) {
	var arr [32]byte
	if len(b) > 32 {
		// If b is longer than 32 bytes, error out to avoid losing data
		return nil, fmt.Errorf("input is longer than 32 bytes")
	} else {
		// If b is shorter than 32 bytes, copy it as is and leave the rest zeroed
		copy(arr[:], b)
	}
	return &arr, nil
}
