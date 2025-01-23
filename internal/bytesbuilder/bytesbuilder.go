package bytesbuilder

import (
	"encoding/hex"
	"math/big"

	"github.com/1inch/1inch-sdk-go/internal/hexidecimal"
	"github.com/ethereum/go-ethereum/common"
)

type BytesBuilder struct {
	data []byte
}

func New() *BytesBuilder {
	return &BytesBuilder{data: []byte{}}
}

func (b *BytesBuilder) AddUint24(val *big.Int) {
	bytes := val.Bytes()
	switch {
	case len(bytes) < 3:
		// Pad on the left with zeros to make it 3 bytes
		padded := make([]byte, 3-len(bytes))
		bytes = append(padded, bytes...)
	case len(bytes) > 3:
		// Truncate any bytes above the 3 least significant
		bytes = bytes[len(bytes)-3:]
	}
	b.data = append(b.data, bytes...)
}

func (b *BytesBuilder) AddUint32(val *big.Int) {
	bytes := val.Bytes()
	if len(bytes) < 4 {
		padded := make([]byte, 4-len(bytes))
		bytes = append(padded, bytes...)
	}
	b.data = append(b.data, bytes...)
}

func (b *BytesBuilder) AddUint16(val *big.Int) {
	bytes := val.Bytes()
	if len(bytes) < 2 {
		padded := make([]byte, 2-len(bytes))
		bytes = append(padded, bytes...)
	}
	b.data = append(b.data, bytes...)
}

func (b *BytesBuilder) AddUint8(val uint8) {
	b.data = append(b.data, byte(val))
}

func (b *BytesBuilder) AddAddress(address common.Address) {
	b.data = append(b.data, address.Bytes()...)
}

func (b *BytesBuilder) AddBytes(data string) error {
	bytes, err := hex.DecodeString(hexidecimal.Trim0x(data))
	if err != nil {
		return err
	}
	b.data = append(b.data, bytes...)
	return nil
}

func (b *BytesBuilder) AsHex() string {
	return hex.EncodeToString(b.data)
}
