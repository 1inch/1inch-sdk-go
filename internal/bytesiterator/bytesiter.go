package bytesiterator

import (
	"errors"
	"math/big"
)

// BytesIter facilitates sequential reading of bytes from a byte slice.
type BytesIter struct {
	data []byte
	pos  int
}

// New initializes a new BytesIter with the provided data.
func New(data []byte) *BytesIter {
	return &BytesIter{data: data, pos: 0}
}

// NextByte reads the next single byte.
func (iter *BytesIter) NextByte() (byte, error) {
	if iter.pos >= len(iter.data) {
		return 0, errors.New("no more bytes to read")
	}
	val := iter.data[iter.pos]
	iter.pos++
	return val, nil
}

// NextUint16 reads the next 2 bytes and returns them as a *big.Int.
func (iter *BytesIter) NextUint16() (*big.Int, error) {
	if iter.pos+2 > len(iter.data) {
		return nil, errors.New("insufficient bytes for uint16")
	}
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+2])
	iter.pos += 2
	return val, nil
}

// NextUint24 reads the next 3 bytes and returns them as a uint32.
func (iter *BytesIter) NextUint24() (uint32, error) {
	if iter.pos+3 > len(iter.data) {
		return 0, errors.New("insufficient bytes for uint24")
	}
	val := uint32(iter.data[iter.pos])<<16 | uint32(iter.data[iter.pos+1])<<8 | uint32(iter.data[iter.pos+2])
	iter.pos += 3
	return val, nil
}

// NextUint32 reads the next 4 bytes and returns them as a *big.Int.
func (iter *BytesIter) NextUint32() (*big.Int, error) {
	if iter.pos+4 > len(iter.data) {
		return nil, errors.New("insufficient bytes for uint32")
	}
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+4])
	iter.pos += 4
	return val, nil
}

// NextUint160 reads the next 20 bytes and returns them as a *big.Int.
func (iter *BytesIter) NextUint160() (*big.Int, error) {
	if iter.pos+20 > len(iter.data) {
		return nil, errors.New("insufficient bytes for uint160")
	}
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+20])
	iter.pos += 20
	return val, nil
}

// NextUint256 reads the next 32 bytes and returns them as a *big.Int.
func (iter *BytesIter) NextUint256() (*big.Int, error) {
	if iter.pos+32 > len(iter.data) {
		return nil, errors.New("insufficient bytes for uint256")
	}
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+32])
	iter.pos += 32
	return val, nil
}

// NextBytes reads the next n bytes and returns them as a byte slice.
func (iter *BytesIter) NextBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, errors.New("negative byte count")
	}
	if iter.pos+n > len(iter.data) {
		return nil, errors.New("insufficient bytes for next bytes")
	}
	val := iter.data[iter.pos : iter.pos+n]
	iter.pos += n
	return val, nil
}

// NextString reads the next n bytes and returns them as a string.
func (iter *BytesIter) NextString(n int) (string, error) {
	bytes, err := iter.NextBytes(n)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Rest returns the remaining bytes as a byte slice.
func (iter *BytesIter) Rest() ([]byte, error) {
	if iter.pos >= len(iter.data) {
		return nil, nil
	}
	val := iter.data[iter.pos:]
	iter.pos = len(iter.data)
	return val, nil
}

// BytesLeft returns the number of bytes that are left unread.
func (iter *BytesIter) BytesLeft() int {
	return len(iter.data) - iter.pos
}

// IsEmpty checks if there are no more bytes to read.
func (iter *BytesIter) IsEmpty() bool {
	return iter.pos >= len(iter.data)
}
