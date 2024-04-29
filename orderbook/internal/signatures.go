package internal

import (
	"encoding/hex"
	"fmt"
)

// CompactSignature represents a compacted form of an Ethereum signature.
type CompactSignature struct {
	R  []byte
	VS []byte
}

// CompressSignature converts a standard 65-byte Ethereum signature into the EIP-2098 compact format.
func CompressSignature(signature string) (*CompactSignature, error) {

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return nil, err
	}

	if len(signatureBytes) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes long")
	}

	// Extract r, s, and v components from the signature
	r := signatureBytes[:32]
	s := make([]byte, 32)
	copy(s, signatureBytes[32:64])
	v := signatureBytes[64]

	// Encode v into s (first bit)
	if v == 28 {
		s[0] |= 0x80 // Set the first bit if v is 28
	}

	return &CompactSignature{
		R:  r,
		VS: s,
	}, nil
}
