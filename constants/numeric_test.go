package constants

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint16Max(t *testing.T) {
	// uint16 max = 2^16 - 1 = 65535
	expected := big.NewInt(65535)
	assert.Equal(t, 0, Uint16Max.Cmp(expected))
}

func TestUint24Max(t *testing.T) {
	// uint24 max = 2^24 - 1 = 16777215
	assert.Equal(t, 16777215, Uint24Max)
}

func TestUint32Max(t *testing.T) {
	// uint32 max = 2^32 - 1 = 4294967295
	assert.Equal(t, 4294967295, Uint32Max)
}

func TestUint40Max(t *testing.T) {
	// uint40 max = 2^40 - 1 = 1099511627775
	expected := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 40), big.NewInt(1))
	assert.Equal(t, 0, Uint40Max.Cmp(expected))
}

func TestUint256Max(t *testing.T) {
	// uint256 max = 2^256 - 1
	expected := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	assert.Equal(t, 0, Uint256Max.Cmp(expected))

	// Should be 78 digits
	assert.Equal(t, 78, len(Uint256Max.String()))
}
