package web3_provider

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestMaxValuesPass(t *testing.T) {

	permit := AllowancePermitSingle{
		Details: AllowancePermitDetails{
			Token:      "0x0000000000000000000000000000000000000000",
			Amount:     "1461501637330902918203684832716283019655932542975",
			Expiration: "281474976710655",
			Nonce:      "281474976710655",
		},
		Spender:     "0x0000000000000000000000000000000000000000",
		SigDeadline: "115792089237316195423570985008687907853269984665640564039457584007913129639935",
	}

	permit2Address := common.HexToAddress("0x0000000000000000000000000000000000000000")
	chainId := big.NewInt(1)

	_, err := GetTypedDataAllowancePermitSingle(permit, permit2Address, int(chainId.Int64()))
	assert.NoError(t, err)

	hash, err := AllowancePermitSingleTypedDataHash(permit, permit2Address, int(chainId.Int64()))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	assert.Equal(t, "0x7c7685afe45d5d39b6279f05214f9bb9aa275f541f950d0a97d0c18aa43158c8", hash)
}
