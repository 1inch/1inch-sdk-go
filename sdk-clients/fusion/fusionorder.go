package fusion

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

func NewFusionOrderParams(params FusionOrderParamsData) *FusionOrderParams {
	fusionOrderParams := &FusionOrderParams{}

	if params.Preset == "" {
		fusionOrderParams.Preset = Fast
	} else {
		fusionOrderParams.Preset = params.Preset
	}

	if params.Receiver == "" {
		fusionOrderParams.Receiver = "0x0000000000000000000000000000000000000000"
	} else {
		fusionOrderParams.Receiver = params.Receiver
	}

	fusionOrderParams.IsPermit2 = params.IsPermit2
	fusionOrderParams.Nonce = params.Nonce
	fusionOrderParams.Permit = params.Permit
	fusionOrderParams.DelayAuctionStartTimeBy = params.DelayAuctionStartTimeBy

	return fusionOrderParams
}

type FusionOrder struct {
	FusionExtension     Extension
	Inner               orderbook.OrderData
	SettlementExtension common.Address
	OrderInfo           FusionOrderV4
	AuctionDetails      AuctionDetails
	PostInteractionData SettlementPostInteractionData
	Extra               ExtraData
}

func CreateFusionOrder(settlementAddress string, orderInfo FusionOrderV4, details Details, extraParams ExtraParams) (*FusionOrder, *orderbook.MakerTraits, *orderbook.Extension, error) {
	settlementExtensionContract := settlementAddress
	auctionDetails := details.Auction
	resolverStartTime := details.ResolvingStartTime
	if details.ResolvingStartTime == nil || details.ResolvingStartTime.Cmp(big.NewInt(0)) == 0 {
		resolverStartTime = big.NewInt(time.Now().Unix())
	}
	postInteractionData := NewSettlementPostInteractionData(SettlementSuffixData{
		Whitelist:          details.Whitelist,
		IntegratorFee:      &details.Fees.IntFee,
		BankFee:            details.Fees.BankFee,
		ResolvingStartTime: resolverStartTime,
		CustomReceiver:     common.HexToAddress(orderInfo.Receiver),
	})
	extra := extraParams // TODO check defaults from SDK

	// TODO ignoring defaults for now
	//defaultExtraParams := ExtraParams{
	//	AllowPartialFills:    true,
	//	AllowMultipleFills:   true,
	//	unwrapWeth:           false,
	//	EnablePermit2:        false,
	//	OrderExpirationDelay: big.NewInt(12),
	//}

	deadline := auctionDetails.StartTime + auctionDetails.Duration + extra.OrderExpirationDelay

	makerTraitParms := orderbook.MakerTraitsParams{
		Expiry:             int64(deadline),
		AllowPartialFills:  extra.AllowPartialFills,
		AllowMultipleFills: extra.AllowMultipleFills,
		HasPostInteraction: true,
		UnwrapWeth:         extra.unwrapWeth,
		UsePermit2:         extra.EnablePermit2,
	}
	if extra.Nonce == nil || extra.Nonce.Cmp(big.NewInt(0)) == 0 {
		makerTraitParms.Nonce = extra.Nonce.Int64()
	}

	makerTraits := orderbook.NewMakerTraits(makerTraitParms)

	if makerTraits.IsBitInvalidatorMode() {
		if extraParams.Nonce == nil || extraParams.Nonce == big.NewInt(0) {
			return nil, nil, nil, errors.New("nonce required, when partial fill or multiple fill disallowed")
		}
	}

	var receiver common.Address
	if postInteractionData.IntegratorFee.Ratio != nil && postInteractionData.IntegratorFee.Ratio.Cmp(big.NewInt(0)) != 0 {
		receiver = common.HexToAddress(settlementExtensionContract)
	} else {
		receiver = common.HexToAddress(orderInfo.Receiver)
	}

	var permitInteraction *Interaction = nil
	if extraParams.Permit != "" {
		permitInteraction = &Interaction{
			Target: common.HexToAddress(orderInfo.MakerAsset),
			Data:   extraParams.Permit,
		}
	}
	extension := NewFusionExtension(common.HexToAddress(settlementExtensionContract), auctionDetails, postInteractionData, permitInteraction)
	builtExtension := extension.Build() //TODO Optionally handle an order-provided base salt value here to enable Injected Salts (not sure what this is yet)
	// TODO injected salts are supposed to be used whenever the provided orderinfo does not have a salt on it

	salt := builtExtension.BuildSalt()

	orderData := orderbook.OrderData{
		MakerAsset:   orderInfo.MakerAsset,
		TakerAsset:   orderInfo.TakerAsset,
		MakingAmount: orderInfo.MakingAmount,
		TakingAmount: orderInfo.TakingAmount,
		Salt:         fmt.Sprintf("%x", salt),
		Maker:        orderInfo.Maker,
		Receiver:     receiver.Hex(),
		MakerTraits:  makerTraits.Encode(),
		Extension:    fmt.Sprintf("%s", builtExtension.keccak256()),
	}

	return &FusionOrder{
		FusionExtension:     builtExtension,
		Inner:               orderData,
		SettlementExtension: common.HexToAddress(settlementExtensionContract),
		OrderInfo:           orderInfo,
		AuctionDetails:      auctionDetails,
		PostInteractionData: postInteractionData,
		Extra: ExtraData{
			UnwrapWETH:           extraParams.unwrapWeth,
			Nonce:                extraParams.Nonce,
			Permit:               extraParams.Permit,
			AllowPartialFills:    extraParams.AllowPartialFills,
			AllowMultipleFills:   extraParams.AllowMultipleFills,
			OrderExpirationDelay: extraParams.OrderExpirationDelay,
			EnablePermit2:        extraParams.EnablePermit2,
			Source:               extraParams.Source,
		},
	}, makerTraits, &orderbook.Extension{}, nil //TODO extension is wrong
}
