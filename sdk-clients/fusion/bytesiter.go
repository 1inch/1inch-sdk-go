package fusion

import "math/big"

type BytesIter struct {
	data []byte
	pos  int
}

func NewBytesIter(data []byte) *BytesIter {
	return &BytesIter{data: data}
}

func (iter *BytesIter) NextByte() byte {
	val := iter.data[iter.pos]
	iter.pos++
	return val
}

var zero = big.NewInt(0)

func (iter *BytesIter) NextUint16() *big.Int {
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+2])
	iter.pos += 2

	// If the resulting value of delay is zero, set it to a fresh big.Int of value zero (for comparisons in tests)
	if val.Cmp(zero) == 0 {
		val = zero
	}
	return val
}

func (iter *BytesIter) NextUint32() *big.Int {
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+4])
	iter.pos += 4
	return val
}

func (iter *BytesIter) NextUint160() *big.Int {
	val := new(big.Int).SetBytes(iter.data[iter.pos : iter.pos+20])
	iter.pos += 20
	return val
}

func (iter *BytesIter) NextBytes(n int) []byte {
	val := iter.data[iter.pos : iter.pos+n]
	iter.pos += n
	return val
}

func (iter *BytesIter) Rest() *big.Int {
	if iter.pos >= len(iter.data) {
		return nil
	}
	return new(big.Int).SetBytes(iter.data[iter.pos:])
}

func (iter *BytesIter) IsEmpty() bool {
	return iter.pos >= len(iter.data)
}
