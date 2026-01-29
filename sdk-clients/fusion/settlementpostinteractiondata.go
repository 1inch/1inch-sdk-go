package fusion

import (
	"fmt"
	"math/big"

	"github.com/1inch/1inch-sdk-go/common/fusionorder"
	"github.com/1inch/1inch-sdk-go/internal/addresses"
	"github.com/1inch/1inch-sdk-go/internal/bytesbuilder"
	"github.com/ethereum/go-ethereum/common"
)

type SettlementPostInteractionData struct {
	Whitelist          []WhitelistItem
	ResolvingStartTime *big.Int
	CustomReceiver     common.Address
	AuctionFees        *FeesIntegratorAndResolver
}

// GenerateWhitelist converts a list of address strings into WhitelistItems with delays.
// This is an alias for fusionorder.GenerateWhitelist.
var GenerateWhitelist = fusionorder.GenerateWhitelist

const customReceiverBitFlag = 0

func CreateEncodedPostInteractionData(extension *Extension) (string, error) {
	builder := bytesbuilder.New()

	customReceiver := extension.PostInteractionData.CustomReceiver
	if customReceiver == (common.Address{}) {
		customReceiver = common.HexToAddress(addresses.ZeroAddress)
	}

	flags := big.NewInt(0)
	if customReceiver.Hex() != addresses.ZeroAddress {
		flags.SetBit(flags, customReceiverBitFlag, 1)
	}
	builder.AddUint8(uint8(flags.Uint64()))

	integratorReceiver := common.HexToAddress(addresses.ZeroAddress)
	if extension.PostInteractionData.AuctionFees != nil && extension.PostInteractionData.AuctionFees.Integrator.Integrator != "" && extension.PostInteractionData.AuctionFees.Integrator.Integrator != addresses.ZeroAddress {
		integratorReceiver = common.HexToAddress(extension.PostInteractionData.AuctionFees.Integrator.Integrator)
	}

	protocolReceiver := common.HexToAddress(addresses.ZeroAddress)
	if extension.PostInteractionData.AuctionFees != nil && extension.PostInteractionData.AuctionFees.Integrator.Protocol != "" && extension.PostInteractionData.AuctionFees.Integrator.Protocol != addresses.ZeroAddress {
		protocolReceiver = common.HexToAddress(extension.PostInteractionData.AuctionFees.Integrator.Protocol)
	}

	builder.AddAddress(integratorReceiver)
	builder.AddAddress(protocolReceiver)

	if flags.Bit(customReceiverBitFlag) == 1 {
		builder.AddAddress(customReceiver)
	}

	params := &BuildAmountGetterDataParams{
		AuctionDetails:      extension.AuctionDetails,
		PostInteractionData: extension.PostInteractionData,
		ResolvingStartTime:  extension.ResolvingStartTime,
	}

	amountGetterData, err := BuildAmountGetterData(params, false)
	if err != nil {
		return "", fmt.Errorf("failed to build amount getter data: %w", err)
	}
	if err := builder.AddBytes(amountGetterData); err != nil {
		return "", fmt.Errorf("failed to add amount getter data: %w", err)
	}

	builder.AddUint256(extension.Surplus.EstimatedTakerAmount)

	protocolFeePercent := extension.Surplus.ProtocolFee.ToPercent(fusionorder.GetDefaultBase())
	builder.AddUint8(uint8(protocolFeePercent))

	return fmt.Sprintf("0x%s", builder.AsHex()), nil
}

func (spid SettlementPostInteractionData) CanExecuteAt(executor common.Address, executionTime *big.Int) bool {
	return fusionorder.CanExecuteAt(spid.Whitelist, spid.ResolvingStartTime, executor, executionTime)
}

func (spid SettlementPostInteractionData) IsExclusiveResolver(wallet common.Address) bool {
	return fusionorder.IsExclusiveResolver(spid.Whitelist, wallet)
}
