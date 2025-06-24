package fusion

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/1inch/1inch-sdk-go/internal/bytesbuilder"
	"github.com/1inch/1inch-sdk-go/internal/hexadecimal"
	"github.com/ethereum/go-ethereum/common"
)

type SettlementPostInteractionData struct {
	Whitelist          []WhitelistItem
	IntegratorFee      *IntegratorFee
	BankFee            *big.Int
	ResolvingStartTime *big.Int
	CustomReceiver     common.Address
	AuctionFees        *FeesNew
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

//func GenerateWhitelist(data *SettlementSuffixData) ([]WhitelistItem, error) {
//	if len(data.Whitelist) == 0 {
//		return nil, errors.New("whitelist cannot be empty")
//	}
//
//	sumDelay := big.NewInt(0)
//	whitelist := make([]WhitelistItem, len(data.Whitelist))
//
//	// Transform timestamps to cumulative delays
//	sort.Slice(data.Whitelist, func(i, j int) bool {
//		return data.Whitelist[i].AllowFrom.Cmp(data.Whitelist[j].AllowFrom) < 0
//	})
//
//	for i, d := range data.Whitelist {
//		allowFrom := d.AllowFrom
//		if d.AllowFrom.Cmp(data.ResolvingStartTime) < 0 {
//			allowFrom = data.ResolvingStartTime
//		}
//
//		zero := big.NewInt(0)
//		delay := new(big.Int).Sub(allowFrom, data.ResolvingStartTime)
//		delay.Sub(delay, sumDelay)
//		// If the resulting value of delay is zero, set it to a fresh big.Int of value zero (for comparisons in tests)
//		if delay.Cmp(zero) == 0 {
//			delay = zero
//		}
//		whitelist[i] = WhitelistItem{
//			AddressHalf: strings.ToLower(d.Address.Hex())[len(d.Address.Hex())-20:],
//			Delay:       delay,
//		}
//
//		sumDelay.Add(sumDelay, whitelist[i].Delay)
//
//		if whitelist[i].Delay.Cmp(uint16Max) >= 0 {
//			return nil, fmt.Errorf("delay too big - %d must be less than %d", whitelist[i].Delay, uint16Max)
//		}
//	}
//
//	return whitelist, nil
//}

func GenerateWhitelist(whitelistStrings []string, resolvingStartTime *big.Int) ([]WhitelistItem, error) {
	if len(whitelistStrings) == 0 {
		return nil, errors.New("whitelist cannot be empty")
	}

	//whitelistAddresses := make([]AuctionWhitelistItem, 0)
	//for _, address := range whitelistStrings {
	//	whitelistAddresses = append(whitelistAddresses, AuctionWhitelistItem{
	//		Address:   geth_common.HexToAddress(address),
	//		AllowFrom: big.NewInt(0), // TODO generating the correct list here requires checking for an exclusive resolver. This needs to be checked for later. The generated object does not see exclusive resolver correctly
	//	})
	//}

	sumDelay := big.NewInt(0)
	whitelist := make([]WhitelistItem, len(whitelistStrings))

	// Transform timestamps to cumulative delays

	// TODO this sorting step is currently skipped since we do not calculate AllowFrom
	//sort.Slice(data.Whitelist, func(i, j int) bool {
	//	return data.Whitelist[i].AllowFrom.Cmp(data.Whitelist[j].AllowFrom) < 0
	//})
	for i, d := range whitelistStrings {
		allowFrom := big.NewInt(0).Set(resolvingStartTime)
		//allowFrom := d.AllowFrom
		//if d.AllowFrom.Cmp(data.ResolvingStartTime) < 0 {
		//	allowFrom = data.ResolvingStartTime
		//}

		zero := big.NewInt(0)
		delay := new(big.Int).Sub(allowFrom, resolvingStartTime)
		delay.Sub(delay, sumDelay)
		// If the resulting value of delay is zero, set it to a fresh big.Int of value zero (for comparisons in tests)
		if delay.Cmp(zero) == 0 {
			delay = zero
		}
		whitelist[i] = WhitelistItem{
			AddressHalf: strings.ToLower(d)[len(d)-20:],
			Delay:       delay,
		}

		sumDelay.Add(sumDelay, whitelist[i].Delay)

		if whitelist[i].Delay.Cmp(uint16Max) >= 0 {
			return nil, fmt.Errorf("delay too big - %d must be less than %d", whitelist[i].Delay, uint16Max)
		}
	}

	return whitelist, nil
}

func Decode(data string) (SettlementPostInteractionData, error) {
	bytes, err := hex.DecodeString(hexadecimal.Trim0x(data))
	if err != nil {
		return SettlementPostInteractionData{}, errors.New("invalid hex string")
	}

	flags := big.NewInt(int64(bytes[len(bytes)-1]))
	bytesWithoutFlags := bytes[:len(bytes)-1]

	iter := NewBytesIter(bytesWithoutFlags)
	var bankFee *big.Int
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
	var whitelist []WhitelistItem

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

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"
const CUSTOM_RECEIVER_FLAG_BIT = 0

func CreateEncodedPostInteractionData(extension *Extension) (string, error) {
	builder := bytesbuilder.New()

	customReceiver := extension.PostInteractionData.CustomReceiver
	if customReceiver == (common.Address{}) {
		customReceiver = common.HexToAddress(ZERO_ADDRESS)
	}

	// Set bit flags
	flags := big.NewInt(0)
	if customReceiver.Hex() != ZERO_ADDRESS {
		flags.SetBit(flags, CUSTOM_RECEIVER_FLAG_BIT, 1)
	}
	builder.AddUint8(uint8(flags.Uint64()))

	// Set receivers
	integratorReceiver := common.HexToAddress(ZERO_ADDRESS)
	if extension.PostInteractionData.AuctionFees != nil && extension.PostInteractionData.AuctionFees.Integrator.Integrator != "" && extension.PostInteractionData.AuctionFees.Integrator.Integrator != ZERO_ADDRESS {
		integratorReceiver = common.HexToAddress(extension.PostInteractionData.AuctionFees.Integrator.Integrator)
	}

	protocolReceiver := common.HexToAddress(ZERO_ADDRESS)
	if extension.PostInteractionData.AuctionFees != nil && extension.PostInteractionData.AuctionFees.Integrator.Protocol != "" && extension.PostInteractionData.AuctionFees.Integrator.Protocol != ZERO_ADDRESS {
		protocolReceiver = common.HexToAddress(extension.PostInteractionData.AuctionFees.Integrator.Protocol)
	}

	// TODO verify 0x is not appended
	builder.AddAddress(integratorReceiver)
	builder.AddAddress(protocolReceiver)

	// Optional customReceiver
	if flags.Bit(CUSTOM_RECEIVER_FLAG_BIT) == 1 {
		builder.AddAddress(customReceiver)
	}

	params := &BuildAmountGetterDataParams{
		AuctionDetails:      extension.AuctionDetails,
		PostInteractionData: extension.PostInteractionData,
		ResolvingStartTime:  extension.ResolvingStartTime,
	}
	// Add amount getter data (forAmountGetters = false)
	amountGetterData, err := BuildAmountGetterData(params, false)
	if err != nil {
		return "", fmt.Errorf("failed to build amount getter data: %w", err)
	}
	if err := builder.AddBytes(amountGetterData); err != nil {
		return "", fmt.Errorf("failed to add amount getter data: %w", err)
	}

	builder.AddUint256(extension.Surplus.EstimatedTakerAmount)

	// Add protocol fee as uint8 percent
	protocolFeePercent := extension.Surplus.ProtocolFee.ToPercent(GetDefaultBase())
	builder.AddUint8(uint8(protocolFeePercent))

	return fmt.Sprintf("0x%s", builder.AsHex()), nil
}

func (spid SettlementPostInteractionData) EncodeNew(extension Extension) (string, error) {
	builder := bytesbuilder.New()

	customReceiver := spid.CustomReceiver
	if customReceiver == (common.Address{}) {
		customReceiver = common.HexToAddress(ZERO_ADDRESS)
	}

	// Set bit flags
	flags := big.NewInt(0)
	if customReceiver.Hex() != ZERO_ADDRESS {
		flags.SetBit(flags, CUSTOM_RECEIVER_FLAG_BIT, 1)
	}
	builder.AddUint8(uint8(flags.Uint64()))

	// Set receivers
	integratorReceiver := common.HexToAddress(ZERO_ADDRESS)
	if spid.AuctionFees != nil && spid.AuctionFees.Integrator.Integrator != "" && spid.AuctionFees.Integrator.Integrator != ZERO_ADDRESS {
		integratorReceiver = common.HexToAddress(spid.AuctionFees.Integrator.Integrator)
	}

	protocolReceiver := common.HexToAddress(ZERO_ADDRESS)
	if spid.AuctionFees != nil && spid.AuctionFees.Integrator.Protocol != "" && spid.AuctionFees.Integrator.Protocol != ZERO_ADDRESS {
		protocolReceiver = common.HexToAddress(spid.AuctionFees.Integrator.Protocol)
	}

	// TODO verify 0x is not appended
	builder.AddAddress(integratorReceiver)
	builder.AddAddress(protocolReceiver)

	// Optional customReceiver
	if flags.Bit(CUSTOM_RECEIVER_FLAG_BIT) == 1 {
		builder.AddAddress(customReceiver)
	}

	params := &BuildAmountGetterDataParams{
		AuctionDetails:      extension.AuctionDetails,
		PostInteractionData: extension.PostInteractionData,
		ResolvingStartTime:  extension.ResolvingStartTime,
	}

	// Add amount getter data (forAmountGetters = false)
	amountGetterData, err := BuildAmountGetterData(params, false)
	if err != nil {
		return "", fmt.Errorf("failed to build amount getter data: %w", err)
	}
	if err := builder.AddBytes(amountGetterData); err != nil {
		return "", fmt.Errorf("failed to add amount getter data: %w", err)
	}

	builder.AddUint256(extension.Surplus.EstimatedTakerAmount)

	// Add protocol fee as uint8 percent
	protocolFeePercent := extension.Surplus.ProtocolFee.ToPercent(GetDefaultBase())
	builder.AddUint8(uint8(protocolFeePercent))

	return fmt.Sprintf("0x%s", builder.AsHex()), nil
}

func (spid SettlementPostInteractionData) Encode() (string, error) {
	bitMask := big.NewInt(0)
	bytes := bytesbuilder.New()

	if spid.BankFee != nil && spid.BankFee.Cmp(big.NewInt(0)) != 0 {
		bitMask.SetBit(bitMask, 0, 1)
		bytes.AddUint32(spid.BankFee)
	}

	if spid.IntegratorFee != nil && spid.IntegratorFee.Ratio != nil && spid.IntegratorFee.Ratio.Cmp(big.NewInt(0)) != 0 {
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
