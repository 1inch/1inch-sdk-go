package orderbook

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/orderbook/models"
)

func normalizeGetOrderByHashResponse(resp *models.GetOrderByHashResponse) (*models.GetOrderByHashResponseExtended, error) {
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

	return &models.GetOrderByHashResponseExtended{
		GetOrderByHashResponse: *resp,
		LimitOrderDataNormalized: models.NormalizedLimitOrderData{
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
	// Normalize the string: remove the "0x" prefix if present
	normStr := strings.TrimPrefix(hexStr, "0x")

	// Create a big.Int
	n := new(big.Int)

	// Set the big.Int based on the string (16 indicates base 16 parsing)
	_, ok := n.SetString(normStr, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hexadecimal string: %s", hexStr)
	}
	return n, nil
}

// BytesToBytes32 converts a byte slice to a [32]byte, padding with zeros if necessary,
// and truncating if it's too long.
func BytesToBytes32(b []byte) [32]byte {
	var arr [32]byte
	if len(b) > 32 {
		// If b is longer than 32 bytes, truncate it to fit into the [32]byte array
		copy(arr[:], b[:32]) // TODO this is an error case as data would be lost
	} else {
		// If b is shorter than 32 bytes, copy it as is and leave the rest zeroed
		copy(arr[:], b)
	}
	return arr
}
