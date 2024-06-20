package fusion

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func stringToHexBytes(hexStr string) ([]byte, error) {
	// Strip the "0x" prefix if it exists
	cleanedStr := strings.TrimPrefix(hexStr, "0x")

	// Ensure the string has an even length by padding with a zero if it's odd
	if len(cleanedStr)%2 != 0 {
		cleanedStr = "0" + cleanedStr
	}

	// Decode the string into bytes
	bytes, err := hex.DecodeString(cleanedStr)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func BuildSalt(extension string) string {
	if extension == "0x" {
		return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
	}

	byteConverted, err := stringToHexBytes(extension)
	if err != nil {
		panic(err)
	}

	keccakHash := crypto.Keccak256Hash(byteConverted)
	salt := new(big.Int).SetBytes(keccakHash.Bytes())
	// We need to keccak256 the extension and then bitwise & it with uint_160_max
	var uint160Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
	salt.And(salt, uint160Max)
	return fmt.Sprintf("0x%x", salt)
}
