package fusion

import (
	"encoding/hex"
	"errors"
	"math/big"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type SettlementPostInteractionData struct {
	Whitelist          []WhitelistItem
	IntegratorFee      *IntegratorFee
	BankFee            *big.Int
	ResolvingStartTime *big.Int
	CustomReceiver     common.Address
}

var UINT_16_MAX = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 16), big.NewInt(1))

func NewSettlementPostInteractionData(data SettlementSuffixData) SettlementPostInteractionData {
	if len(data.Whitelist) == 0 {
		panic("Whitelist cannot be empty")
	}

	sumDelay := big.NewInt(0)
	whitelist := make([]WhitelistItem, len(data.Whitelist))

	// Transform timestamps to cumulative delays
	sort.Slice(data.Whitelist, func(i, j int) bool {
		return data.Whitelist[i].AllowFrom.Cmp(data.Whitelist[j].AllowFrom) < 0
	})

	for i, d := range data.Whitelist {
		allowFrom := d.AllowFrom
		if d.AllowFrom.Cmp(data.ResolvingStartTime) < 0 {
			allowFrom = data.ResolvingStartTime
		}

		delay := new(big.Int).Sub(allowFrom, data.ResolvingStartTime)
		delay.Sub(delay, sumDelay)
		whitelist[i] = WhitelistItem{
			AddressHalf: d.Address.Hex()[len(d.Address)-20:],
			Delay:       delay,
		}

		sumDelay.Add(sumDelay, whitelist[i].Delay)

		if whitelist[i].Delay.Cmp(UINT_16_MAX) >= 0 {
			panic("Too big diff between timestamps")
		}
	}

	return SettlementPostInteractionData{
		Whitelist:          whitelist,
		IntegratorFee:      data.IntegratorFee,
		BankFee:            data.BankFee,
		ResolvingStartTime: data.ResolvingStartTime,
		CustomReceiver:     data.CustomReceiver,
	}
}

func DecodeSettlementPostInteractionData(data string) (SettlementPostInteractionData, error) {
	bytes, err := hex.DecodeString(strings.TrimPrefix(data, "0x"))
	if err != nil {
		return SettlementPostInteractionData{}, errors.New("invalid hex string")
	}

	iter := NewBytesIterNew(bytes)
	flags := big.NewInt(int64(iter.NextByte()))
	bankFee := big.NewInt(0)
	var integratorFee *IntegratorFee
	var customReceiver common.Address

	if flags.Bit(0) == 1 {
		bankFee = iter.NextUint32()
	}

	if flags.Bit(1) == 1 {
		integratorFee = &IntegratorFee{
			Ratio:    iter.NextUint16(),
			Receiver: common.HexToAddress(iter.NextUint160().Text(16)),
		}

		if flags.Bit(2) == 1 {
			customReceiver = common.HexToAddress(iter.NextUint160().Text(16))
		}
	}

	resolvingStartTime := iter.NextUint32()
	whitelist := []WhitelistItem{}

	for !iter.IsEmpty() {
		addressHalf := hex.EncodeToString(iter.NextBytes(10))
		delay := iter.NextUint16()
		whitelist = append(whitelist, WhitelistItem{
			AddressHalf: addressHalf,
			Delay:       delay,
		})
	}

	return SettlementPostInteractionData{
		IntegratorFee:      integratorFee,
		BankFee:            bankFee,
		ResolvingStartTime: resolvingStartTime,
		Whitelist:          whitelist,
		CustomReceiver:     customReceiver,
	}, nil
}

func (spid SettlementPostInteractionData) Encode() string {
	bitMask := big.NewInt(0)
	bytes := NewBytesBuilder()

	if spid.BankFee != nil && spid.BankFee.Cmp(big.NewInt(0)) != 0 {
		bitMask.SetBit(bitMask, 0, 1)
		bytes.AddUint32(spid.BankFee)
	}

	if spid.IntegratorFee != nil && spid.IntegratorFee.Ratio.Cmp(big.NewInt(0)) != 0 {
		bitMask.SetBit(bitMask, 1, 1)
		bytes.AddUint16(spid.IntegratorFee.Ratio)
		bytes.AddAddress(spid.IntegratorFee.Receiver)

		// TODO this check is probably not good enough
		if spid.CustomReceiver.Hex() != "0x0000000000000000000000000000000000000000" {
			bitMask.SetBit(bitMask, 2, 1)
			bytes.AddAddress(spid.CustomReceiver)
		}
	}

	bytes.AddUint32(spid.ResolvingStartTime)

	for _, wl := range spid.Whitelist {
		bytes.AddBytes(wl.AddressHalf)
		bytes.AddUint16(wl.Delay)
	}

	bitMask.SetInt64(int64(len(spid.Whitelist)) << 3)
	bytes.AddUint8(uint8(bitMask.Int64()))

	return bytes.AsHex()
}

func (spid SettlementPostInteractionData) CanExecuteAt(executor common.Address, executionTime *big.Int) bool {
	addressHalf := executor.Hex()[len(executor.Hex())-20:]

	allowedFrom := spid.ResolvingStartTime

	for _, whitelist := range spid.Whitelist {
		allowedFrom.Add(allowedFrom, whitelist.Delay)

		if addressHalf == whitelist.AddressHalf {
			return executionTime.Cmp(allowedFrom) >= 0
		} else if executionTime.Cmp(allowedFrom) < 0 {
			return false
		}
	}

	return false
}

func (spid SettlementPostInteractionData) IsExclusiveResolver(wallet common.Address) bool {
	addressHalf := wallet.Hex()[len(wallet.Hex())-20:]

	if len(spid.Whitelist) == 1 {
		return addressHalf == spid.Whitelist[0].AddressHalf
	}

	if spid.Whitelist[0].Delay.Cmp(spid.Whitelist[1].Delay) == 0 {
		return false
	}

	return addressHalf == spid.Whitelist[0].AddressHalf
}

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

func (iter *BytesIterNew) IsEmpty() bool {
	return iter.pos >= len(iter.data)
}

type BytesBuilder struct {
	data []byte
}

func NewBytesBuilder() *BytesBuilder {
	return &BytesBuilder{data: []byte{}}
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
	b.data = append(b.data, []byte(address.Hex())...)
}

func (b *BytesBuilder) AddBytes(data string) {
	bytes, err := hex.DecodeString(strings.TrimPrefix(data, "0x"))
	if err != nil {
		panic("invalid hex string")
	}
	b.data = append(b.data, bytes...)
}

func (b *BytesBuilder) AsHex() string {
	return hex.EncodeToString(b.data)
}
