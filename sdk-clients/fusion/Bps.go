package fusion

import (
	"fmt"
	"math/big"
)

type Bps struct {
	value *big.Int
}

const defaultBase = 1

func GetDefaultBase() *big.Int {
	return big.NewInt(defaultBase)
}

var BpsZero = NewBps(big.NewInt(0))

func NewBps(val *big.Int) *Bps {
	if val.Cmp(big.NewInt(0)) < 0 || val.Cmp(big.NewInt(10000)) > 0 {
		panic(fmt.Sprintf("invalid bps %s", val.String()))
	}
	return &Bps{value: new(big.Int).Set(val)}
}

func FromPercent(val float64, base *big.Int) *Bps {
	mult := new(big.Float).SetFloat64(100 * val)
	return fromFloatWithBase(mult, base)
}

func FromFraction(val float64, base *big.Int) *Bps {
	mult := new(big.Float).SetFloat64(10000 * val)
	return fromFloatWithBase(mult, base)
}

func fromFloatWithBase(f *big.Float, base *big.Int) *Bps {
	baseFloat := new(big.Float).SetInt(base)
	res := new(big.Float).Quo(f, baseFloat)

	bpsInt, _ := res.Int(nil) // round down
	return NewBps(bpsInt)
}

func (b *Bps) Equal(other *Bps) bool {
	return b.value.Cmp(other.value) == 0
}

func (b *Bps) IsZero() bool {
	return b.value.Sign() == 0
}

func (b *Bps) ToPercent(base *big.Int) float64 {
	num := new(big.Int).Mul(b.value, base)
	f := new(big.Float).SetInt(num)
	den := big.NewFloat(100)
	result, _ := new(big.Float).Quo(f, den).Float64()
	return result
}

func (b *Bps) ToFraction(base *big.Int) *big.Int {
	num := new(big.Int).Mul(b.value, base) // numerator = bps.value * base
	den := big.NewInt(10000)               // denominator = 10000
	result := new(big.Int).Div(num, den)   // integer division
	return result
}

func (b *Bps) String() string {
	return b.value.String()
}
