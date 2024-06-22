package fusion

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/1inch/1inch-sdk-go/sdk-clients/orderbook"
)

var timeNow func() int64 = GetCurrentTime

func GetCurrentTime() int64 {
	return time.Now().Unix()
}

func CreateSettlementPostInteractionData(details Details, orderInfo FusionOrderV4) (*SettlementPostInteractionData, error) {
	resolverStartTime := details.ResolvingStartTime
	if details.ResolvingStartTime == nil || details.ResolvingStartTime.Cmp(big.NewInt(0)) == 0 {
		resolverStartTime = big.NewInt(timeNow())
	}
	fmt.Printf("Resolver start time: %v\n", resolverStartTime)
	return NewSettlementPostInteractionData(SettlementSuffixData{
		Whitelist:          details.Whitelist,
		IntegratorFee:      &details.Fees.IntFee,
		BankFee:            details.Fees.BankFee,
		ResolvingStartTime: resolverStartTime,
		CustomReceiver:     common.HexToAddress(orderInfo.Receiver),
	})
}

type CreateExtensionParams struct {
	settlementAddress   string
	postInteractionData *SettlementPostInteractionData
	orderInfo           FusionOrderV4
	details             Details
	extraParams         ExtraParams
}

func CreateExtension(params CreateExtensionParams) (*Extension, error) {

	var permitInteraction *Interaction
	if params.extraParams.Permit != "" {
		permitInteraction = &Interaction{
			Target: common.HexToAddress(params.orderInfo.MakerAsset),
			Data:   params.extraParams.Permit,
		}
	}

	extensionBuilder := ExtensionBuilder{}
	settlementAddressContract := common.HexToAddress(params.settlementAddress)
	auctionDetails := params.details.Auction.Encode()
	extensionBuilder.WithMakingAmountData(settlementAddressContract, auctionDetails)
	extensionBuilder.WithTakingAmountData(settlementAddressContract, auctionDetails)
	postInteraction := NewInteraction(settlementAddressContract, params.postInteractionData.Encode())
	extensionBuilder.WithPostInteraction(postInteraction)
	if permitInteraction != nil {
		extensionBuilder.WithMakerPermit(permitInteraction.Target, permitInteraction.Data)
	}
	return extensionBuilder.Build(), nil
}

func CreateMakerTraits(details Details, extraParams ExtraParams) (*orderbook.MakerTraits, error) {
	deadline := details.Auction.StartTime + details.Auction.Duration + extraParams.OrderExpirationDelay
	makerTraitParms := orderbook.MakerTraitsParams{
		Expiry:             int64(deadline),
		AllowPartialFills:  extraParams.AllowPartialFills,
		AllowMultipleFills: extraParams.AllowMultipleFills,
		HasPostInteraction: true,
		UnwrapWeth:         extraParams.unwrapWeth,
		UsePermit2:         extraParams.EnablePermit2,
		HasExtension:       true,
		Nonce:              extraParams.Nonce.Int64(),
	}
	makerTraits := orderbook.NewMakerTraits(makerTraitParms)
	if makerTraits.IsBitInvalidatorMode() {
		if extraParams.Nonce == nil || extraParams.Nonce.Cmp(big.NewInt(0)) == 0 {
			return nil, errors.New("nonce required, when partial fill or multiple fill disallowed")
		}
	}
	return makerTraits, nil
}

type CreateOrderDataParams struct {
	settlementAddress   string
	postInteractionData *SettlementPostInteractionData
	extension           *Extension
	orderInfo           FusionOrderV4
	details             Details
	extraParams         ExtraParams
	makerTraits         *orderbook.MakerTraits
}

func CreateOrder(params CreateOrderDataParams) (*Order, error) {
	var receiver common.Address
	if params.postInteractionData.IntegratorFee.Ratio != nil && params.postInteractionData.IntegratorFee.Ratio.Cmp(big.NewInt(0)) != 0 {
		receiver = common.HexToAddress(params.settlementAddress)
	} else {
		receiver = common.HexToAddress(params.orderInfo.Receiver)
	}

	salt, err := params.extension.GenerateSalt() // TODO this is not the right salt
	if err != nil {
		return nil, fmt.Errorf("error generating salt: %v", err)
	}

	return &Order{
		FusionExtension: params.extension,
		Inner: orderbook.OrderData{
			MakerAsset:   params.orderInfo.MakerAsset,
			TakerAsset:   params.orderInfo.TakerAsset,
			MakingAmount: params.orderInfo.MakingAmount,
			TakingAmount: params.orderInfo.TakingAmount,
			Salt:         fmt.Sprintf("%x", salt),
			Maker:        params.orderInfo.Maker,
			Receiver:     receiver.Hex(),
			MakerTraits:  params.makerTraits.Encode(),
			Extension:    fmt.Sprintf("%x", params.extension.keccak256()),
		},
		SettlementExtension: common.HexToAddress(params.settlementAddress),
		OrderInfo:           params.orderInfo,
		AuctionDetails:      params.details.Auction,
		PostInteractionData: params.postInteractionData,
		Extra: ExtraData{
			UnwrapWETH:           params.extraParams.unwrapWeth,
			Nonce:                params.extraParams.Nonce,
			Permit:               params.extraParams.Permit,
			AllowPartialFills:    params.extraParams.AllowPartialFills,
			AllowMultipleFills:   params.extraParams.AllowMultipleFills,
			OrderExpirationDelay: params.extraParams.OrderExpirationDelay,
			EnablePermit2:        params.extraParams.EnablePermit2,
			Source:               params.extraParams.Source,
		},
	}, nil
}
