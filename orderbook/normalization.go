package orderbook

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func NormalizeGetOrderByHashResponse(resp *GetOrderByHashResponse) (*GetOrderByHashResponseExtended, error) {
	saltBigInt, ok := new(big.Int).SetString(resp.Data.Salt, 10)
	if !ok {
		return nil, fmt.Errorf("invalid salt value")
	}
	makingAmountBigInt, ok := new(big.Int).SetString(resp.Data.MakingAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid making amount value")
	}
	takingAmountBigInt, ok := new(big.Int).SetString(resp.Data.TakingAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid taking amount value")
	}

	makerAssetBigInt := AddressStringToBigInt(resp.Data.MakerAsset)
	takerAssetBigInt := AddressStringToBigInt(resp.Data.TakerAsset)
	makerBigInt := AddressStringToBigInt(resp.Data.Maker)
	receiverBigInt := AddressStringToBigInt(resp.Data.Receiver)

	makerTraits, err := HexStringToBigInt(resp.Data.MakerTraits)
	if err != nil {
		return nil, fmt.Errorf("invalid maker traits value")
	}

	return &GetOrderByHashResponseExtended{
		GetOrderByHashResponse: *resp,
		LimitOrderDataNormalized: NormalizedLimitOrderData{
			Salt:         saltBigInt,
			MakerAsset:   makerAssetBigInt,
			TakerAsset:   takerAssetBigInt,
			Maker:        makerBigInt,
			Receiver:     receiverBigInt,
			MakingAmount: makingAmountBigInt,
			TakingAmount: takingAmountBigInt,
			MakerTraits:  makerTraits,
		},
	}, nil
}

func AddressStringToBigInt(addressString string) *big.Int {
	address := common.HexToAddress(addressString)
	addressBytes := new(big.Int).SetBytes(address.Bytes())
	return addressBytes
}

func HexStringToBigInt(hexStr string) (*big.Int, error) {
	normStr := strings.TrimPrefix(hexStr, "0x")
	n := new(big.Int)
	_, ok := n.SetString(normStr, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hexadecimal string: %s", hexStr)
	}
	return n, nil
}
