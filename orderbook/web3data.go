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
