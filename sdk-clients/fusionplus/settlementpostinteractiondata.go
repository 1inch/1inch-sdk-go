package fusionplus

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/internal/bytesbuilder"
	"github.com/1inch/1inch-sdk-go/internal/bytesiterator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// SettlementPostInteractionData represents post interaction data for FusionPlus orders.
// The IntegratorFee and BankFee fields are optional - when nil, they are not encoded.
type SettlementPostInteractionData struct {
	Whitelist          []fusionorder.WhitelistItem
	IntegratorFee      *IntegratorFee // Optional: only used in extension encoding
	BankFee            *big.Int       // Optional: only used in extension encoding
	ResolvingStartTime *big.Int
	CustomReceiver     common.Address
}

// NewSettlementPostInteractionData creates a new SettlementPostInteractionData from suffix data.
// This version does not include fees (IntegratorFee and BankFee will be nil).
func NewSettlementPostInteractionData(data SettlementSuffixData) (*SettlementPostInteractionData, error) {
	whitelist, err := fusionorder.GenerateWhitelistFromItems(data.Whitelist, data.ResolvingStartTime)
	if err != nil {
		return nil, err
	}

	return &SettlementPostInteractionData{
		Whitelist:          whitelist,
		IntegratorFee:      nil,
		BankFee:            nil,
		ResolvingStartTime: data.ResolvingStartTime,
		CustomReceiver:     data.CustomReceiver,
	}, nil
}

// NewSettlementPostInteractionDataWithFees creates a new SettlementPostInteractionData including fee information.
// This is used for extension encoding which requires fee data.
func NewSettlementPostInteractionDataWithFees(data SettlementSuffixData) (*SettlementPostInteractionData, error) {
	whitelist, err := fusionorder.GenerateWhitelistFromItems(data.Whitelist, data.ResolvingStartTime)
	if err != nil {
		return nil, err
	}

	return &SettlementPostInteractionData{
		Whitelist:          whitelist,
		IntegratorFee:      data.IntegratorFee,
		BankFee:            data.BankFee,
		ResolvingStartTime: data.ResolvingStartTime,
		CustomReceiver:     data.CustomReceiver,
	}, nil
}

// Encode encodes the post interaction data to hex string.
// If IntegratorFee and BankFee are present, they are included in the encoding.
func (spid SettlementPostInteractionData) Encode() (string, error) {
	bitMask := big.NewInt(0)
	bytes := bytesbuilder.New()

	// Encode fees if present
	if spid.BankFee != nil && spid.BankFee.Cmp(big.NewInt(0)) != 0 {
		bitMask.SetBit(bitMask, 0, 1)
		bytes.AddUint32(spid.BankFee)
	}

	if spid.IntegratorFee != nil && spid.IntegratorFee.Ratio.Cmp(big.NewInt(0)) != 0 {
		bitMask.SetBit(bitMask, 1, 1)
		bytes.AddUint16(spid.IntegratorFee.Ratio)
		bytes.AddAddress(spid.IntegratorFee.Receiver)

		if spid.CustomReceiver.Hex() != "0x0000000000000000000000000000000000000000" {
			bitMask.SetBit(bitMask, 2, 1)
			bytes.AddAddress(spid.CustomReceiver)
		}
	}

	bytes.AddUint32(spid.ResolvingStartTime)

	for _, wl := range spid.Whitelist {
		err := bytes.AddBytes(wl.AddressHalf)
		if err != nil {
			return "", err
		}
		bytes.AddUint16(wl.Delay)
	}

	bitMask.Or(bitMask, big.NewInt(int64(len(spid.Whitelist)<<3)))
	bytes.AddUint8(uint8(bitMask.Int64()))

	return fmt.Sprintf("0x%s", bytes.AsHex()), nil
}

// DecodeSettlementPostInteractionData decodes hex data into SettlementPostInteractionData
func DecodeSettlementPostInteractionData(data string) (*SettlementPostInteractionData, error) {
	bytes, err := hexutil.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
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
			return nil, err
		}
	}

	if flags.Bit(1) == 1 {
		ratio, err := iter.NextUint16()
		if err != nil {
			return nil, err
		}

		receiver, err := iter.NextUint160()
		if err != nil {
			return nil, err
		}

		integratorFee = &IntegratorFee{
			Ratio:    ratio,
			Receiver: common.HexToAddress(receiver.Text(16)),
		}

		if flags.Bit(2) == 1 {
			customReceiverRaw, err := iter.NextUint160()
			if err != nil {
				return nil, err
			}
			customReceiver = common.HexToAddress(customReceiverRaw.Text(16))
		}
	}

	resolvingStartTime, err := iter.NextUint32()
	if err != nil {
		return nil, err
	}

	var whitelist []fusionorder.WhitelistItem
	for !iter.IsEmpty() {
		addressHalfRaw, err := iter.NextBytes(10)
		if err != nil {
			return nil, err
		}
		addressHalf := hex.EncodeToString(addressHalfRaw)
		delay, err := iter.NextUint16()
		if err != nil {
			return nil, err
		}
		whitelist = append(whitelist, fusionorder.WhitelistItem{
			AddressHalf: addressHalf,
			Delay:       delay,
		})
	}

	return &SettlementPostInteractionData{
		IntegratorFee:      integratorFee,
		BankFee:            bankFee,
		ResolvingStartTime: resolvingStartTime,
		Whitelist:          whitelist,
		CustomReceiver:     customReceiver,
	}, nil
}

func (spid SettlementPostInteractionData) CanExecuteAt(executor common.Address, executionTime *big.Int) bool {
	return fusionorder.CanExecuteAt(spid.Whitelist, spid.ResolvingStartTime, executor, executionTime)
}

func (spid SettlementPostInteractionData) IsExclusiveResolver(wallet common.Address) bool {
	return fusionorder.IsExclusiveResolver(spid.Whitelist, wallet)
}
