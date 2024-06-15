package fusion

import "math/big"

// Utility structures and functions

type BytesIterNew struct {
	data []byte
	pos  int
}

func NewBytesIterNew(data []byte) *BytesIterNew {
	return &BytesIterNew{data: data}
}

func (iter *BytesIterNew) NextByte() byte {
	val := iter.data[iter.pos]
	iter.pos++
	return val
}

func (iter *BytesIterNew) NextUint16() *big.Int {
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+2])
	iter.pos += 2
	return val
}

func (iter *BytesIterNew) NextUint32() *big.Int {
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+4])
	iter.pos += 4
	return val
}

func (iter *BytesIterNew) NextUint160() *big.Int {
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+20])
	iter.pos += 20
	return val
}

func (iter *BytesIterNew) NextBytes(n int) []byte {
	val := iter.data[iter.pos : iter.pos+n]
	iter.pos += n
	return val
}

func (iter *BytesIterNew) Rest() *big.Int {
	if iter.pos >= len(iter.data) {
		return nil
	}
	return new(big.Int).SetBytes(iter.data[iter.pos:])
}

func (iter *BytesIterNew) IsEmpty() bool {
	return iter.pos >= len(iter.data)
}
