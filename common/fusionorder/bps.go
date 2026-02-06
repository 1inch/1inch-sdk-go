package fusionorder

import (
	"fmt"
	"math/big"
)

// Bps represents basis points (1/100th of a percent, so 10000 bps = 100%)
type Bps struct {
	value *big.Int
}

const defaultBase = 1

// GetDefaultBase returns the default base value for Bps calculations
func GetDefaultBase() *big.Int {
	return big.NewInt(defaultBase)
}

// BpsZero is a zero basis points value
var BpsZero = MustNewBps(big.NewInt(0))

// NewBps creates a new Bps value, returns error if value is not in [0, 10000]
func NewBps(val *big.Int) (*Bps, error) {
	if val.Cmp(big.NewInt(0)) < 0 || val.Cmp(big.NewInt(10000)) > 0 {
		return nil, fmt.Errorf("bps value %s out of range [0, 10000]", val.String())
	}
	return &Bps{value: new(big.Int).Set(val)}, nil
}

// MustNewBps creates a new Bps value, panics if value is not in [0, 10000]
// Use this only for known-valid constant values at package initialization
func MustNewBps(val *big.Int) *Bps {
	bps, err := NewBps(val)
	if err != nil {
		panic(err)
	}
	return bps
}

// FromPercent creates a Bps from a percentage value
// Example: FromPercent(1, GetDefaultBase()) creates 100 bps (1%)
func FromPercent(val float64, base *big.Int) (*Bps, error) {
	mult := new(big.Float).SetFloat64(100 * val)
	return fromFloatWithBase(mult, base)
}

// MustFromPercent creates a Bps from a percentage value, panicking on error.
// Use this only for known-valid constant values at package initialization or in tests.
func MustFromPercent(val float64, base *big.Int) *Bps {
	bps, err := FromPercent(val, base)
	if err != nil {
		panic(err)
	}
	return bps
}

// FromFraction creates a Bps from a fraction
// Example: FromFraction(0.01, GetDefaultBase()) creates 100 bps (1%)
func FromFraction(val float64, base *big.Int) (*Bps, error) {
	mult := new(big.Float).SetFloat64(10000 * val)
	return fromFloatWithBase(mult, base)
}

// MustFromFraction creates a Bps from a fraction, panicking on error.
// Use this only for known-valid constant values at package initialization or in tests.
func MustFromFraction(val float64, base *big.Int) *Bps {
	bps, err := FromFraction(val, base)
	if err != nil {
		panic(err)
	}
	return bps
}

func fromFloatWithBase(f *big.Float, base *big.Int) (*Bps, error) {
	baseFloat := new(big.Float).SetInt(base)
	res := new(big.Float).Quo(f, baseFloat)

	bpsInt, _ := res.Int(nil) // round down
	return NewBps(bpsInt)
}

// Equal returns true if two Bps values are equal
func (b *Bps) Equal(other *Bps) bool {
	return b.value.Cmp(other.value) == 0
}

// IsZero returns true if the Bps value is zero
func (b *Bps) IsZero() bool {
	return b.value.Sign() == 0
}

// ToPercent converts Bps to a percentage
func (b *Bps) ToPercent(base *big.Int) float64 {
	num := new(big.Int).Mul(b.value, base)
	f := new(big.Float).SetInt(num)
	den := big.NewFloat(100)
	result, _ := new(big.Float).Quo(f, den).Float64()
	return result
}

// ToFraction converts Bps to a fraction with the given base
func (b *Bps) ToFraction(base *big.Int) *big.Int {
	num := new(big.Int).Mul(b.value, base) // numerator = bps.value * base
	den := big.NewInt(10000)               // denominator = 10000
	result := new(big.Int).Div(num, den)   // integer division
	return result
}

// String returns the string representation of the Bps value
func (b *Bps) String() string {
	return b.value.String()
}

// Value returns the underlying big.Int value
func (b *Bps) Value() *big.Int {
	return new(big.Int).Set(b.value)
}
