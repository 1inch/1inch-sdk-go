package random_number_generation

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var BigIntMaxFunc func(*big.Int) (*big.Int, error) = BigIntMax

// BigIntMax generates a random big.Int from 0 to max
func BigIntMax(max *big.Int) (*big.Int, error) {
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, fmt.Errorf("error generating random number: %v", err)
	}
	return n, nil
}
