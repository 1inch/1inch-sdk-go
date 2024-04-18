package web3_provider

import (
	"context"
	"math/big"

	"github.com/1inch/1inch-sdk-go/internal/web3-provider/multicall"

	gethCommon "github.com/ethereum/go-ethereum/common"
)

func (w Wallet) GetSeriesNonce(ctx context.Context, seriesNonceManager gethCommon.Address, publicAddress gethCommon.Address) (*big.Int, error) {
	function := "nonce"

	seriesNonceData, err := w.seriesNonceManagerABI.Pack(function, big.NewInt(0), publicAddress)
	if err != nil {
		return nil, err
	}

	callDataArray := []multicall.CallData{
		multicall.BuildCallData(seriesNonceManager, seriesNonceData, 0),
	}

	mResult, err := w.multicall.Execute(ctx, callDataArray)
	if err != nil {
		return nil, err
	}

	var nonce *big.Int
	err = w.seriesNonceManagerABI.UnpackIntoInterface(&nonce, function, mResult[0])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}
