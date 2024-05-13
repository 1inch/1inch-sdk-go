package orderbook

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitMask(t *testing.T) {
	tests := []struct {
		name               string
		startBit           int64
		endBit             int64
		valueToUpdate      string
		inputBits          string
		expectedOutputBits string
	}{
		{
			name:               "Simple single bit mask",
			startBit:           0,
			endBit:             1,
			valueToUpdate:      "0",
			inputBits:          "1",
			expectedOutputBits: "1",
		},
		{
			name:               "Set middle bits",
			startBit:           4,
			endBit:             8,
			valueToUpdate:      "110000000000",
			inputBits:          "1111",
			expectedOutputBits: "110011110000",
		},
		{
			name:               "Clear bits",
			startBit:           4,
			endBit:             8,
			valueToUpdate:      "11111111",
			inputBits:          "0",
			expectedOutputBits: "00001111",
		},
		{
			name:               "Set bits in an existing value",
			startBit:           4,
			endBit:             6,
			valueToUpdate:      "11110000",
			inputBits:          "11",
			expectedOutputBits: "11110000",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			startBit := big.NewInt(tc.startBit)
			endBit := big.NewInt(tc.endBit)
			valueToUpdate := bitStringToBigInt(tc.valueToUpdate)
			inputBits := bitStringToBigInt(tc.inputBits)
			expectedOutputBits := bitStringToBigInt(tc.expectedOutputBits)

			bitmask := NewBitMask(startBit, endBit)
			result := bitmask.SetBits(valueToUpdate, inputBits)
			assert.Equal(t, expectedOutputBits, result)
		})
	}
}

func TestBitMaskToString(t *testing.T) {
	tests := []struct {
		name           string
		startBit       int64
		endBit         int64
		expectedString string
	}{
		{
			name:           "Simple mask",
			startBit:       4,
			endBit:         8,
			expectedString: "0xf0",
		},
		{
			name:           "Single bit mask",
			startBit:       0,
			endBit:         1,
			expectedString: "0x1",
		},
		{
			name:           "Full byte mask",
			startBit:       0,
			endBit:         8,
			expectedString: "0xff",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bitmask := NewBitMask(big.NewInt(tc.startBit), big.NewInt(tc.endBit))
			assert.Equal(t, tc.expectedString, bitmask.ToString())
		})
	}
}

func TestBitMaskToBigInt(t *testing.T) {
	tests := []struct {
		name           string
		startBit       int64
		endBit         int64
		expectedBigInt *big.Int
	}{
		{
			name:           "Simple mask",
			startBit:       4,
			endBit:         8,
			expectedBigInt: bitStringToBigInt("11110000"),
		},
		{
			name:           "Single bit mask",
			startBit:       0,
			endBit:         1,
			expectedBigInt: bitStringToBigInt("00000001"),
		},
		{
			name:           "Full byte mask",
			startBit:       0,
			endBit:         8,
			expectedBigInt: bitStringToBigInt("11111111"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bitmask := NewBitMask(big.NewInt(tc.startBit), big.NewInt(tc.endBit))
			assert.Equal(t, tc.expectedBigInt, bitmask.ToBigInt())
		})
	}
}

func bitStringToBigInt(bitStr string) *big.Int {
	i := new(big.Int)
	i.SetString(bitStr, 2)
	return i
}
