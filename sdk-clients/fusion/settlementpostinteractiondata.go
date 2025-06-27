package fusion

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/1inch/1inch-sdk-go/internal/addresses"
	"github.com/1inch/1inch-sdk-go/internal/bytesbuilder"
	"github.com/ethereum/go-ethereum/common"
)

type SettlementPostInteractionData struct {
	Whitelist          []WhitelistItem
	ResolvingStartTime *big.Int
	CustomReceiver     common.Address
	AuctionFees        *FeesIntegratorResolver
}

var uint16Max = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 16), big.NewInt(1))

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

const CUSTOM_RECEIVER_FLAG_BIT = 0

func CreateEncodedPostInteractionData(extension *Extension) (string, error) {
	builder := bytesbuilder.New()

	customReceiver := extension.PostInteractionData.CustomReceiver
	if customReceiver == (common.Address{}) {
		customReceiver = common.HexToAddress(addresses.ZeroAddress)
	}

	// Set bit flags
	flags := big.NewInt(0)
	if customReceiver.Hex() != addresses.ZeroAddress {
		flags.SetBit(flags, CUSTOM_RECEIVER_FLAG_BIT, 1)
	}
	builder.AddUint8(uint8(flags.Uint64()))

	// Set receivers
	integratorReceiver := common.HexToAddress(addresses.ZeroAddress)
	if extension.PostInteractionData.AuctionFees != nil && extension.PostInteractionData.AuctionFees.Integrator.Integrator != "" && extension.PostInteractionData.AuctionFees.Integrator.Integrator != addresses.ZeroAddress {
		integratorReceiver = common.HexToAddress(extension.PostInteractionData.AuctionFees.Integrator.Integrator)
	}

	protocolReceiver := common.HexToAddress(addresses.ZeroAddress)
	if extension.PostInteractionData.AuctionFees != nil && extension.PostInteractionData.AuctionFees.Integrator.Protocol != "" && extension.PostInteractionData.AuctionFees.Integrator.Protocol != addresses.ZeroAddress {
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
		customReceiver = common.HexToAddress(addresses.ZeroAddress)
	}

	// Set bit flags
	flags := big.NewInt(0)
	if customReceiver.Hex() != addresses.ZeroAddress {
		flags.SetBit(flags, CUSTOM_RECEIVER_FLAG_BIT, 1)
	}
	builder.AddUint8(uint8(flags.Uint64()))

	// Set receivers
	integratorReceiver := common.HexToAddress(addresses.ZeroAddress)
	if spid.AuctionFees != nil && spid.AuctionFees.Integrator.Integrator != "" && spid.AuctionFees.Integrator.Integrator != addresses.ZeroAddress {
		integratorReceiver = common.HexToAddress(spid.AuctionFees.Integrator.Integrator)
	}

	protocolReceiver := common.HexToAddress(addresses.ZeroAddress)
	if spid.AuctionFees != nil && spid.AuctionFees.Integrator.Protocol != "" && spid.AuctionFees.Integrator.Protocol != addresses.ZeroAddress {
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
