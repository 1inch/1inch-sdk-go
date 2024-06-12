package fusion

import (
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Interaction struct {
	Target common.Address
	Data   string
}

func NewInteraction(target common.Address, data string) Interaction {
	if !isHexBytes(data) {
		panic("Interaction data must be valid hex bytes")
	}
	return Interaction{
		Target: target,
		Data:   data,
	}
}

func (i *Interaction) Encode() string {
	return i.Target.String() + trim0x(i.Data)
}

func DecodeInteraction(bytes string) *Interaction {
	iter := NewBytesIter(bytes)
	return &Interaction{
		Target: common.Address(iter.NextUint160()),
		Data:   iter.Rest(),
	}
}

type BytesIter struct {
	data []byte
	pos  int
}

func NewBytesIter(hexStr string) *BytesIter {
	data, err := hex.DecodeString(strings.TrimPrefix(hexStr, "0x"))
	if err != nil {
		panic("Invalid hex string")
	}
	return &BytesIter{data: data}
}

func (iter *BytesIter) NextUint160() []byte {
	if iter.pos+20 > len(iter.data) {
		panic("Not enough bytes for uint160")
	}
	val := iter.data[iter.pos : iter.pos+20]
	iter.pos += 20
	return val
}

func (iter *BytesIter) Rest() string {
	if iter.pos >= len(iter.data) {
		return ""
	}
	return hex.EncodeToString(iter.data[iter.pos:])
}

func isHexBytes(s string) bool {
	_, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	return err == nil
}
