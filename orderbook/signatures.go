package orderbook

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

	// Adjust v according to Ethereum's usual 27-28 scheme if necessary
	if v < 27 {
		v += 27
	}

	// Encode v into s (last byte, first bit)
	if v == 28 {
		s[31] |= 0x80 // Set the first bit if v is 28
	}

	return &CompactSignature{
		R:  r,
		VS: s,
	}, nil
}
