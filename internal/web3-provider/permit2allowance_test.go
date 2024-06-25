package web3_provider

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// Define the maximum values
var (
	MaxAllowanceTransferAmount = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
	MaxAllowanceExpiration     = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
	MaxOrderedNonce            = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 48), big.NewInt(1))
	MaxSigDeadline             = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
)

func TestMaxValuesPass(t *testing.T) {
	permitDetails := PermitDetails{
		Token:      common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Amount:     MaxAllowanceTransferAmount,
		Expiration: MaxAllowanceExpiration,
		Nonce:      MaxOrderedNonce,
	}

	permit := PermitSingle{
		Details:     permitDetails,
		Spender:     common.HexToAddress("0x0000000000000000000000000000000000000000"),
		SigDeadline: MaxSigDeadline,
	}

	permit2Address := common.HexToAddress("0x0000000000000000000000000000000000000000")
	chainId := big.NewInt(1) // For example, 1 for Ethereum mainnet

	_, err := hashPermitSingle(permit, permit2Address, chainId)
	assert.NoError(t, err)
}
