package orderbook

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

type GetSaltParams struct {
	Extension string
	Source    *string // Optional string for tracking code
	UseRandom bool    // If true, uses random bits for the middle section; otherwise uses timestamp
}

// GenerateSaltNew generates a salt value with specific bit patterns
func GenerateSaltNew(params *GetSaltParams) (*big.Int, error) {
	salt := big.NewInt(0)

	// Generate upper 32 bits (bits 224-255) - tracking code mask
	trackingSource := "sdk"
	if params.Source != nil {
		trackingSource = *params.Source
	}
	trackingHash := crypto.Keccak256Hash([]byte(trackingSource))
	trackingCodeMask := new(big.Int).Lsh(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 32), big.NewInt(1)), 224) // (2^32 - 1) << 224
	trackingBits := new(big.Int).SetBytes(trackingHash.Bytes())
	trackingBits.And(trackingBits, trackingCodeMask)
	salt.Or(salt, trackingBits)

	// Generate middle 64 bits (bits 160-223)
	var middleBits *big.Int
	if params.UseRandom {
		// Generate random 64 bits
		randomBytes := make([]byte, 8)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random bytes: %w", err)
		}
		middleBits = new(big.Int).SetBytes(randomBytes)
	} else {
		middleBits = big.NewInt(time.Now().Unix())
	}
	// Mask to 64 bits and shift to position 160-223
	mask64 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 64), big.NewInt(1))
	middleBits.And(middleBits, mask64)
	middleBits.Lsh(middleBits, 160)
	salt.Or(salt, middleBits)

	// Handle extension for lower 160 bits
	if params.Extension == "0x" || params.Extension == "" {
		// If there is no extension, salt can be anything for the lower bits
		// (middle bits already set above, lower 160 bits remain 0 or can be left as-is)
	} else {
		// Lower 160 bits must be from keccak256 hash of the extension
		extensionBytes, err := stringToHexBytes(params.Extension)
		if err != nil {
			return nil, fmt.Errorf("failed to convert extension to bytes: %w", err)
		}
		extensionHash := crypto.Keccak256Hash(extensionBytes)
		mask160 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 160), big.NewInt(1))
		extensionBits := new(big.Int).SetBytes(extensionHash.Bytes())
		extensionBits.And(extensionBits, mask160)
		salt.Or(salt, extensionBits)
	}

	return salt, nil
}
