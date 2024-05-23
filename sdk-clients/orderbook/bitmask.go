package orderbook

import (
	"fmt"
	"math/big"
)

type BitMask struct {
	Offset *big.Int
	Mask   *big.Int
}

// NewBitMask creates a new BitMask with the given start and end bit positions.
func NewBitMask(startBit, endBit *big.Int) *BitMask {
	if startBit.Cmp(endBit) >= 0 {
		panic("BitMask: startBit must be less than endBit")
	}

	bitCount := new(big.Int).Sub(endBit, startBit)                                                    // endBit - startBit
	mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(bitCount.Uint64())), big.NewInt(1)) // (1 << bitCount) - 1

	return &BitMask{
		Offset: startBit,
		Mask:   mask,
	}
}

func (b *BitMask) SetBits(value, bits *big.Int) *big.Int {
	// Create the shifted mask
	shiftedMask := new(big.Int).Set(b.Mask)
	shiftedMask.Lsh(shiftedMask, uint(b.Offset.Uint64()))
	// Clear the bits at the mask location
	value.And(value, new(big.Int).Not(shiftedMask))
	// Shift the bits to the correct location
	shiftedBits := new(big.Int).Lsh(bits, uint(b.Offset.Uint64()))
	value.Or(value, shiftedBits)
	return value
}

// ToString returns the string representation of the mask shifted by the offset.
func (b *BitMask) ToString() string {
	shiftedMask := new(big.Int).Lsh(b.Mask, uint(b.Offset.Uint64()))
	return fmt.Sprintf("0x%x", shiftedMask)
}

// ToBigInt returns the mask value as a big.Int, shifted by the offset.
func (b *BitMask) ToBigInt() *big.Int {
	return new(big.Int).Lsh(b.Mask, uint(b.Offset.Uint64()))
}
