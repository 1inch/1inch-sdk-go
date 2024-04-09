package validate

import (
	"fmt"
	"math/big"
)

func BigIntFromString(s string) (*big.Int, error) {
	bigInt, ok := new(big.Int).SetString(s, 10) // base 10 for decimal
	if !ok {
		return nil, fmt.Errorf("failed to convert string (%v) to big.Int", s)
	}
	return bigInt, nil
}

// HasDuplicates checks if the provided slice contains any duplicate elements.
// It accepts a slice of any comparable type and returns true if there are duplicates, otherwise it returns false.
func HasDuplicates[T comparable](slice []T) bool {
	seen := make(map[T]bool)
	for _, v := range slice {
		if seen[v] {
			return true
		}
		seen[v] = true
	}
	return false
}

// IsSubset checks if all elements of sliceA are also present in sliceB.
// It returns true if sliceA is a subset of sliceB, otherwise it returns false.
func IsSubset[T comparable](sliceA, sliceB []T) bool {
	setB := make(map[T]bool)
	for _, v := range sliceB {
		setB[v] = true
	}

	for _, v := range sliceA {
		if !setB[v] {
			return false
		}
	}
	return true
}

// Contains checks if the slice contains the given value.
func Contains[T comparable](value T, sliceB []T) bool {
	for _, v := range sliceB {
		if v == value {
			return true
		}
	}
	return false
}
