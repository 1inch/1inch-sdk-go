package fusionorder

import (
	"encoding/json"
	"fmt"
	"math/big"

	random_number_generation "github.com/1inch/1inch-sdk-go/internal/random-number-generation"
	"golang.org/x/crypto/sha3"
)

// Keccak256Hash calculates the Keccak256 hash of any JSON-serializable data
func Keccak256Hash(data any) (*big.Int, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data for hashing: %w", err)
	}
	hash := sha3.New256()
	hash.Write(jsonData)
	return new(big.Int).SetBytes(hash.Sum(nil)), nil
}

// GenerateSaltWithExtension generates a salt value incorporating extension hash
// If extension is nil or empty (based on isEmpty check), returns a random base salt
// Otherwise combines the random salt with the extension hash
func GenerateSaltWithExtension(extensionHash *big.Int, isEmpty bool) (*big.Int, error) {
	// Define the maximum value (2^96 - 1)
	maxValue := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 96), big.NewInt(1))

	// Generate a random big.Int within the range [0, 2^96 - 1]
	baseSalt, err := random_number_generation.BigIntMaxFunc(maxValue)
	if err != nil {
		return nil, err
	}

	if isEmpty {
		return baseSalt, nil
	}

	uint160Max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))

	salt := new(big.Int).Lsh(baseSalt, 160)
	salt.Or(salt, new(big.Int).And(extensionHash, uint160Max))

	return salt, nil
}
