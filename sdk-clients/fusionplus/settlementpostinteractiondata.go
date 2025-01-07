package fusionplus

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/1inch/1inch-sdk-go/internal/bytesbuilder"
	"github.com/1inch/1inch-sdk-go/internal/bytesiterator"
	"github.com/ethereum/go-ethereum/common"
)

type SettlementPostInteractionData struct {
	Whitelist          []WhitelistItem
	IntegratorFee      *IntegratorFee
	BankFee            *big.Int
	ResolvingStartTime *big.Int
	CustomReceiver     common.Address
}

var uint16Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 16), big.NewInt(1))

func NewSettlementPostInteractionData(data SettlementSuffixData) (*SettlementPostInteractionData, error) {
	if len(data.Whitelist) == 0 {
		return nil, errors.New("whitelist cannot be empty")
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

		zero := big.NewInt(0)
		delay := new(big.Int).Sub(allowFrom, data.ResolvingStartTime)
		delay.Sub(delay, sumDelay)
		// If the resulting value of delay is zero, set it to a fresh big.Int of value zero (for comparisons in tests)
		if delay.Cmp(zero) == 0 {
			delay = zero
		}
		whitelist[i] = WhitelistItem{
			AddressHalf: strings.ToLower(d.Address.Hex())[len(d.Address.Hex())-20:],
			Delay:       delay,
		}

		sumDelay.Add(sumDelay, whitelist[i].Delay)

		if whitelist[i].Delay.Cmp(uint16Max) >= 0 {
			return nil, fmt.Errorf("delay too big - %d must be less than %d", whitelist[i].Delay, uint16Max)
		}
	}

	return &SettlementPostInteractionData{
		Whitelist:          whitelist,
		IntegratorFee:      data.IntegratorFee,
		BankFee:            data.BankFee,
		ResolvingStartTime: data.ResolvingStartTime,
		CustomReceiver:     data.CustomReceiver,
	}, nil
}

func Decode(data string) (SettlementPostInteractionData, error) {
	bytes, err := hex.DecodeString(strings.TrimPrefix(data, "0x"))
	if err != nil {
		return SettlementPostInteractionData{}, errors.New("invalid hex string")
	}

	flags := big.NewInt(int64(bytes[len(bytes)-1]))
	bytesWithoutFlags := bytes[:len(bytes)-1]

	iter := bytesiterator.New(bytesWithoutFlags)
	var bankFee *big.Int
	var integratorFee *IntegratorFee
	var customReceiver common.Address

	if flags.Bit(0) == 1 {
		bankFee, err = iter.NextUint32()
		if err != nil {
			return SettlementPostInteractionData{}, err
		}
	}

	if flags.Bit(1) == 1 {

		ratio, err := iter.NextUint16()
		if err != nil {
			return SettlementPostInteractionData{}, err
		}

		receiver, err := iter.NextUint160()
		if err != nil {
			return SettlementPostInteractionData{}, err
		}

		integratorFee = &IntegratorFee{
			Ratio:    ratio,
			Receiver: common.HexToAddress(receiver.Text(16)),
		}

		if flags.Bit(2) == 1 {

			customReceiverRaw, err := iter.NextUint160()
			if err != nil {
				return SettlementPostInteractionData{}, err
			}

			customReceiver = common.HexToAddress(customReceiverRaw.Text(16))
		}
	}

	resolvingStartTime, err := iter.NextUint32()
	if err != nil {
		return SettlementPostInteractionData{}, err
	}
	var whitelist []WhitelistItem

	for !iter.IsEmpty() {
		addressHalfRaw, err := iter.NextBytes(10)
		if err != nil {
			return SettlementPostInteractionData{}, err
		}
		addressHalf := hex.EncodeToString(addressHalfRaw)
		delay, err := iter.NextUint16()
		if err != nil {
			return SettlementPostInteractionData{}, err
		}
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
	bytes := bytesbuilder.New()

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

	bitMask.Or(bitMask, big.NewInt(int64(len(spid.Whitelist)<<3)))
	bytes.AddUint8(uint8(bitMask.Int64()))

	output := fmt.Sprintf("0x%s", bytes.AsHex())

	return output
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
