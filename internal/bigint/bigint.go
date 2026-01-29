package bigint

import (
	"fmt"
	"math/big"
)

var Base1E5 = big.NewInt(100000)
var Base1E2 = big.NewInt(100)

func FromString(s string) (*big.Int, error) {
	bigInt, ok := new(big.Int).SetString(s, 10) // base 10 for decimal
	if !ok {
		return nil, fmt.Errorf("invalid big.Int string: %s", s)
	}
	return bigInt, nil
}
