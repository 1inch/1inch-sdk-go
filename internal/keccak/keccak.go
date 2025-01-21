package keccak

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

func Keccak256Legacy(value []byte) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(value)
	return fmt.Sprintf("0x%x", hash.Sum(nil))
}
